package gocb

import (
	"github.com/couchbase/gocbcore/v10"
)

type connectionManager interface {
	connect() error
	openBucket(bucketName string) error
	buildConfig(cluster *Cluster) error
	getKvProvider(bucketName string) (kvProvider, error)
	getKvCapabilitiesProvider(bucketName string) (kvCapabilityVerifier, error)
	getViewProvider(bucketName string) (viewProvider, error)
	getQueryProvider() (queryProvider, error)
	getAnalyticsProvider() (analyticsProvider, error)
	getSearchProvider() (searchProvider, error)
	getHTTPProvider(bucketName string) (httpProvider, error)
	getDiagnosticsProvider(bucketName string) (diagnosticsProvider, error)
	getWaitUntilReadyProvider(bucketName string) (waitUntilReadyProvider, error)
	connection(bucketName string) (*gocbcore.Agent, error)
	close() error
}

// TODO: change how managers are selected.
func newConnectionMgr(protocol string) connectionManager {
	switch protocol {
	case "protostellar":
		return &psConnectionMgr{}
	default:
		client := &stdConnectionMgr{}
		return client
	}
}
