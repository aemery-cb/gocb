package gocb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	gocbcore "github.com/couchbase/gocbcore/v10"
)

type queryProviderCoreProvider interface {
	N1QLQuery(ctx context.Context, opts gocbcore.N1QLQueryOptions) (qOut queryRowReader, errOut error)
	PreparedN1QLQuery(ctx context.Context, opts gocbcore.N1QLQueryOptions) (qOut queryRowReader, errOut error)
}

type queryProviderCore struct {
	provider queryProviderCoreProvider

	retryStrategyWrapper *retryStrategyWrapper
	transcoder           Transcoder
	timeouts             TimeoutsConfig
	tracer               RequestTracer
	meter                *meterWrapper
}

func (qpc *queryProviderCore) Query(statement string, s *Scope, opts *QueryOptions) (*QueryResult, error) {
	start := time.Now()
	defer qpc.meter.ValueRecord(meterValueServiceQuery, "query", start)

	span := createSpan(qpc.tracer, opts.ParentSpan, "query", "query")
	span.SetAttribute("db.statement", statement)
	if s != nil {
		span.SetAttribute("db.name", s.BucketName())
		span.SetAttribute("db.couchbase.scope", s.Name())
	}
	defer span.End()

	retryStrategy := qpc.retryStrategyWrapper
	if opts.RetryStrategy != nil {
		retryStrategy = newRetryStrategyWrapper(opts.RetryStrategy)
	}

	queryOpts, err := opts.toMap()
	if err != nil {
		return nil, QueryError{
			InnerError:      wrapError(err, "failed to generate query options"),
			Statement:       statement,
			ClientContextID: opts.ClientContextID,
		}
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = qpc.timeouts.QueryTimeout
	}
	deadline := time.Now().Add(timeout)

	queryOpts["statement"] = statement
	if s != nil {
		queryOpts["query_context"] = fmt.Sprintf("%s.%s", s.BucketName(), s.Name())
	}

	eSpan := qpc.tracer.RequestSpan(span.Context(), "request_encoding")
	eSpan.SetAttribute("db.system", "couchbase")
	reqBytes, err := json.Marshal(queryOpts)
	eSpan.End()
	if err != nil {
		return nil, QueryError{
			InnerError:      wrapError(err, "failed to marshall query body"),
			Statement:       maybeGetQueryOption(queryOpts, "statement"),
			ClientContextID: maybeGetQueryOption(queryOpts, "client_context_id"),
		}
	}

	var res queryRowReader
	var qErr error
	if opts.Adhoc {
		res, qErr = qpc.provider.N1QLQuery(opts.Context, gocbcore.N1QLQueryOptions{
			Payload:       reqBytes,
			RetryStrategy: retryStrategy,
			Deadline:      deadline,
			TraceContext:  span.Context(),
			User:          opts.Internal.User,
			Endpoint:      opts.Internal.Endpoint,
		})
	} else {
		res, qErr = qpc.provider.PreparedN1QLQuery(opts.Context, gocbcore.N1QLQueryOptions{
			Payload:       reqBytes,
			RetryStrategy: retryStrategy,
			Deadline:      deadline,
			TraceContext:  span.Context(),
			User:          opts.Internal.User,
			Endpoint:      opts.Internal.Endpoint,
		})
	}
	if qErr != nil {
		return nil, maybeEnhanceQueryError(qErr)
	}

	return newQueryResult(res), nil
}

func maybeGetQueryOption(options map[string]interface{}, name string) string {
	if value, ok := options[name].(string); ok {
		return value
	}
	return ""
}
