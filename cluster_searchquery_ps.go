package gocb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	cbsearch "github.com/couchbase/gocb/v2/search"
	"github.com/couchbase/goprotostellar/genproto/search_v1"
)

type searchProviderPs struct {
	agent search_v1.SearchServiceClient
}

var _ searchProvider = &searchProviderPs{}

// executes a search query against PS, taking care of the translation.
func (search *searchProviderPs) SearchQuery(indexName string, query cbsearch.Query, opts *SearchOptions) (*SearchResult, error) {
	psQuery, err := cbsearch.Internal{}.MapQueryToPs(query)
	if err != nil {
		return nil, err
	}

	psSort, err := cbsearch.Internal{}.MapSortToPs(opts.Sort)
	if err != nil {
		return nil, err
	}

	facets, err := cbsearch.Internal{}.MapFacetsToPs(opts.Facets)
	if err != nil {
		return nil, err
	}

	request := search_v1.SearchQueryRequest{
		IndexName:          indexName,
		Query:              psQuery,
		Limit:              opts.Limit,
		Skip:               opts.Skip,
		IncludeExplanation: opts.Explain,
		HighlightFields:    opts.Highlight.Fields,
		Fields:             opts.Fields,
		Sort:               psSort,
		DisableScoring:     opts.DisableScoring,
		Collections:        opts.Collections,
		IncludeLocations:   opts.IncludeLocations,
		Facets:             facets,
	}

	switch opts.ScanConsistency {
	case SearchScanConsistencyNotBounded:
		request.ScanConsistency = search_v1.SearchQueryRequest_SCAN_CONSISTENCY_NOT_BOUNDED
	default:
		return nil, fmt.Errorf("invalid scan consistency option specified")
	}

	if opts.Highlight != nil {
		switch opts.Highlight.Style {
		case DefaultHighlightStyle:
			request.HighlightStyle = search_v1.SearchQueryRequest_HIGHLIGHT_STYLE_DEFAULT
		case AnsiHightlightStyle:
			request.HighlightStyle = search_v1.SearchQueryRequest_HIGHLIGHT_STYLE_ANSI
		case HTMLHighlightStyle:
			request.HighlightStyle = search_v1.SearchQueryRequest_HIGHLIGHT_STYLE_HTML
		default:
			return nil, fmt.Errorf("invalid highlight option specified")
		}

	}

	client, err := search.agent.SearchQuery(opts.Context, &request)
	if err != nil {
		return nil, err
	}

	return newSearchResult(&psSearchRowReader{client: client}), nil
}

// wrapper around the PS result to make it compatible with
// the searchRowReader interface.
type psSearchRowReader struct {
	client        search_v1.SearchService_SearchQueryClient
	nextRowsIndex int
	nextRows      []*search_v1.SearchQueryResponse_SearchQueryRow
	err           error
	meta          *search_v1.SearchQueryResponse_MetaData
}

// returns the next search row, either from local or fetches it from the client.
func (reader *psSearchRowReader) NextRow() []byte {
	// we have results so lets use them.
	if reader.nextRowsIndex < len(reader.nextRows) {
		row := reader.nextRows[reader.nextRowsIndex]
		reader.nextRowsIndex++

		convertedRow, err := psSearchRowToJsonBytes(row)
		if err != nil {
			reader.finishWithError(err)
			return nil
		}

		return convertedRow
	}

	// check if there are anymore available results.
	res, err := reader.client.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			reader.finishWithoutError()
			return nil
		}
		reader.finishWithError(err)
		return nil
	}

	reader.nextRows = res.GetHits()
	reader.nextRowsIndex = 1
	reader.meta = res.MetaData

	if len(res.Hits) > 0 {
		convertedRow, err := psSearchRowToJsonBytes(res.Hits[0])
		if errors.Is(err, io.EOF) {
			reader.finishWithoutError()
			return nil
		}
		return convertedRow
	}

	return nil
}

func (reader *psSearchRowReader) Close() error {
	if reader.err != nil {
		return reader.err
	}
	// if the client is nil then we must be closed already.
	if reader.client == nil {
		return nil
	}
	err := reader.client.CloseSend()
	reader.client = nil
	return err
}

func (reader *psSearchRowReader) MetaData() ([]byte, error) {
	if reader.err != nil {
		return nil, reader.Err()
	}
	if reader.client != nil {
		return nil, errors.New("the result must be fully read before accessing the meta-data")
	}
	if reader.meta == nil {
		return nil, errors.New("an error occurred during querying which has made the meta-data unavailable")
	}

	meta := jsonSearchResponse{
		TotalHits: reader.meta.Metrics.TotalRows,
		MaxScore:  reader.meta.Metrics.MaxScore,
		Took:      uint64(reader.meta.Metrics.ExecutionTime.GetNanos()), // this is in nano seconds
		Status: jsonSearchResponseStatus{
			Errors: reader.meta.Errors,
		},
		//Facets: //TODO: add support for facets.
	}

	return json.Marshal(meta)
}

func (reader *psSearchRowReader) Err() error {
	return reader.err
}

func (reader *psSearchRowReader) finishWithoutError() {
	// Close the stream now that we are done with it
	err := reader.client.CloseSend()
	if err != nil {
		logWarnf("query stream close failed after meta-data: %s", err)
	}

	reader.client = nil
}

func (reader *psSearchRowReader) finishWithError(err error) {
	// Lets record the error that happened
	reader.err = err

	// Lets close the underlying stream
	closeErr := reader.client.CloseSend()
	if closeErr != nil {
		// We log this at debug level, but its almost always going to be an
		// error since thats the most likely reason we are in finishWithError
		logDebugf("query stream close failed after error: %s", closeErr)
	}

	// Our client is invalidated as soon as an error occurs
	reader.client = nil
}

// Helper functions to convert from PS world into something gocb can process
func psSearchRowLocationToJsonSearchRowLocations(locations []*search_v1.SearchQueryResponse_Location) jsonSearchRowLocations {
	jsonForm := make(jsonSearchRowLocations)

	for _, location := range locations {
		field := location.GetField()
		term := location.GetTerm()
		jsonForm[field][term] = append(jsonForm[field][term], jsonRowLocation{
			Field:          field,
			Term:           term,
			Position:       location.GetPosition(),
			Start:          location.GetStart(),
			End:            location.GetEnd(),
			ArrayPositions: location.GetArrayPositions(),
		})
	}

	return jsonForm
}

func psSearchRowFragmentToMap(fragmentMap map[string]*search_v1.SearchQueryResponse_Fragment) map[string][]string {
	var result = make(map[string][]string)
	for key, fragment := range fragmentMap {
		result[key] = fragment.GetContent()
	}

	return result
}

// helper util to convert PS's SearchQueryRow to jsonSearchRow.
func psSearchRowToJsonSearchRow(row *search_v1.SearchQueryResponse_SearchQueryRow) (jsonSearchRow, error) {
	fieldRaw, err := json.Marshal(row.Fields)
	if err != nil {
		return jsonSearchRow{}, err
	}

	return jsonSearchRow{
		ID:          row.Id,
		Index:       row.Index,
		Score:       row.Score,
		Explanation: row.Explanation,
		Locations:   psSearchRowLocationToJsonSearchRowLocations(row.Locations),
		Fragments:   psSearchRowFragmentToMap(row.Fragments),
		Fields:      fieldRaw,
	}, nil

}

// converts from ps search results to jsonRowMessage as bytes for compatibility with existing gocb code.
func psSearchRowToJsonBytes(row *search_v1.SearchQueryResponse_SearchQueryRow) ([]byte, error) {
	convertedRow, err := psSearchRowToJsonSearchRow(row)
	if err != nil {
		return nil, err
	}

	rowBytes, err := json.Marshal(convertedRow)
	if err != nil {
		return nil, err
	}
	return rowBytes, nil

}
