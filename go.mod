module github.com/couchbase/gocb/v2

replace github.com/couchbase/gocbcoreps => ../gocbcoreps

require (
	github.com/couchbase/gocbcore/v10 v10.2.3
	github.com/couchbase/gocbcoreps v0.0.1-2
	github.com/couchbase/goprotostellar v0.0.1-3
	github.com/couchbaselabs/gocaves/client v0.0.0-20230404095311-05e3ba4f0259
	github.com/couchbaselabs/gocbconnstr/v2 v2.0.0-20230419073701-6f1441e69de9
	github.com/google/uuid v1.3.0
	github.com/stretchr/testify v1.8.2
	google.golang.org/protobuf v1.30.0
)

go 1.13

replace github.com/couchbaselabs/gocbconnstr/v2 => ../gocbconnstr
