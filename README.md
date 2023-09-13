# s3-get

Simple utility to download S3 Object with standard AWS Credentials resolution.

## Installation:

```sh
$ go get github.com/stupoid/s3-get
```

## Usage:

```sh
$ ./s3-get -h
Usage of ./s3-get:
  -dst string
        Destination output path.
  -endpoint string
        Optional custom endpoint to use (for connecting to S3 compatible interfaces).
        Alternatively this can be set via environment variable 'AWS_ENDPOINT_URL'.
        But this flag will take precedence.
  -src string
        Source S3 URI.
  -timeout duration
        Download timeout. (default 1m0s)
```

### Example usage

Typical usage with AWS creds already loaded

```sh
$ ./s3-get -src=s3://my-bucket/credentials/private.crt -dst=./cert.pem
```

Usage with local mock S3 server

```sh
$ ./s3-get -src=s3://my-bucket/credentials/private.crt -dst=./cert.pem -endpoint=http://127.0.0.1:4566
```

Alternate way to set custom aws endpoint

```sh
$ export AWS_ENDPOINT_URL=http://127.0.0.1:4566
$ ./s3-get -src=s3://my-bucket/credentials/private.crt -dst=./cert.pem
```
