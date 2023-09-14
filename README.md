# s3-get

Simple utility to download S3 Object with standard AWS Credentials resolution.

## Installation:

```sh
$ go install github.com/stupoid/s3-get
```

## Usage:

```sh
$ ./s3-get -h
Usage of ./s3-get:
  -concurrency int
    	The number of goroutines to spin up in parallel when sending parts.
    	If this is set to zero, the DefaultDownloadConcurrency value will be used.

    	Concurrency of 1 will download the parts sequentially.

    	Concurrency is ignored if the Range input parameter is provided.
  -dst string
    	Destination output path.
  -endpoint string
    	AWS endpoint to use (for connecting to S3 compatible interfaces).
    	Alternatively this can be set via environment variable 'S3_GET_ENDPOINT'.

    	If this is not set, the default AWS service endpoint will be used.

    	The value set by this flag will take precedence over the environment variable.
  -part-size int
    	The size (in bytes) to request from S3 for each part.
    	The minimum allowed part size is 5MB, and  if this value is set to zero,
    	the DefaultDownloadPartSize value will be used.

    	PartSize is ignored if the Range input parameter is provided.
  -quiet
    	Silence error output. (default true)
  -src string
    	Source S3 URI.
  -timeout duration
    	Download timeout.
    	If this is set to zero, no timeout will be set.

    	A duration string is a possibly signed sequence of decimal numbers,
    	each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
    	Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
```

### Example usage

Typical usage with AWS creds already loaded

```sh
$ ./s3-get -src=s3://my-bucket/credentials/private.crt -dst=./cert.pem
```

Usage with local mock S3 server

```sh
# start up a local mock s3 server (e.g. localstack)
$ docker --rm -d run -p 4566:4566 localstack/localstack -e SERVICES s3
```

```
$ ./s3-get -src=s3://my-bucket/credentials/private.crt -dst=./cert.pem -endpoint=http://127.0.0.1:4566
```

Alternate way to set custom aws endpoint

```sh
$ S3_GET_ENDPOINT=http://127.0.0.1:4566 ./s3-get -src=s3://my-bucket/credentials/private.crt -dst=./cert.pem
```
