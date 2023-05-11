package gocb

import (
	"errors"
	"fmt"
	"sync"

	"github.com/couchbase/gocbcore/v10"
	gocbcoreps "github.com/couchbase/gocbcoreps"
)

var ErrNotImplemented = errors.New("not implemented")

type psConnectionMgr struct {
	host   string
	lock   sync.Mutex
	config *gocbcoreps.DialOptions
	agent  *gocbcoreps.RoutingClient

	timeouts TimeoutsConfig
	tracer   RequestTracer
	meter    *meterWrapper
}

func (c *psConnectionMgr) connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	client, err := gocbcoreps.Dial(c.host, c.config)
	if err != nil {
		return err
	}

	c.agent = client

	return nil
}

func (c *psConnectionMgr) openBucket(bucketName string) error {
	c.agent.OpenBucket(bucketName)
	return nil
}

func (c *psConnectionMgr) buildConfig(cluster *Cluster) error {
	c.host = fmt.Sprintf("%s:%d", cluster.connSpec().Addresses[0].Host, cluster.connSpec().Addresses[0].Port)

	creds, err := cluster.authenticator().Credentials(AuthCredsRequest{})
	if err != nil {
		return err
	}

	c.config = &gocbcoreps.DialOptions{
		Username: creds[0].Username,
		Password: creds[0].Password,
	}

	return nil
}
func (c *psConnectionMgr) getKvProvider(bucketName string) (kvProvider, error) {
	kv := c.agent.KvV1()
	return &kvProviderPs{client: kv}, nil
}
func (c *psConnectionMgr) getKvCapabilitiesProvider(bucketName string) (kvCapabilityVerifier, error) {
	return &gocbcore.AgentInternal{}, ErrNotImplemented
}
func (c *psConnectionMgr) getViewProvider(bucketName string) (viewProvider, error) {
	return &viewProviderWrapper{}, ErrNotImplemented
}
func (c *psConnectionMgr) getQueryProvider() (queryProvider, error) {
	provider := c.agent.QueryV1()
	return &queryProviderPs{
		provider: provider,
		timeouts: c.timeouts,
		tracer:   c.tracer,
		meter:    c.meter,
	}, nil
}
func (c *psConnectionMgr) getAnalyticsProvider() (analyticsProvider, error) {
	return &analyticsProviderWrapper{}, ErrNotImplemented
}
func (c *psConnectionMgr) getSearchProvider() (searchProvider, error) {
	return &searchProviderWrapper{}, ErrNotImplemented
}
func (c *psConnectionMgr) getHTTPProvider(bucketName string) (httpProvider, error) {
	return &httpProviderWrapper{}, ErrNotImplemented
}
func (c *psConnectionMgr) getDiagnosticsProvider(bucketName string) (diagnosticsProvider, error) {
	return &diagnosticsProviderWrapper{}, ErrNotImplemented
}
func (c *psConnectionMgr) getWaitUntilReadyProvider(bucketName string) (waitUntilReadyProvider, error) {
	return &waitUntilReadyProviderWrapper{}, ErrNotImplemented
}
func (c *psConnectionMgr) connection(bucketName string) (*gocbcore.Agent, error) {
	return &gocbcore.Agent{}, ErrNotImplemented
}
func (c *psConnectionMgr) close() error {
	return ErrNotImplemented
}
