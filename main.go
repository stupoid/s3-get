package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// newS3Client returns a s3.Client with configurable endpoint
// If the endpoint provided is an empty string, it will default to AWS
func newS3Client(endpoint string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	return client, nil
}

// ParseS3URI extracts the bucket and object key from a S3 URI
// example S3 URI: s3://bucket-name/object-key
func ParseS3URI(s3URI string) (bucket, key string, err error) {
	u, err := url.Parse(s3URI)
	if err != nil {
		return "", "", err
	}
	if u.Scheme != "s3" {
		return "", "", fmt.Errorf("invalid scheme %v", u.Scheme)
	}
	bucket = u.Host
	key = strings.TrimLeft(u.Path, "/")
	return bucket, key, nil
}

// Download retrieves an object from S3 and writes it to a file.
func Download(d *manager.Downloader, timeout time.Duration, bucket, key, fileName string) (int64, error) {
	ctx := context.TODO()
	var cancel context.CancelFunc

	if timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	file, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	bytesDownloaded, err := d.Download(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		os.Remove(fileName) // cleanup empty file if fail to get Object
		return 0, err
	}
	return bytesDownloaded, nil
}

func main() {
	var s3URI, fileName, endpoint string
	var timeout time.Duration
	var quiet bool
	var partSize int64
	var concurrency int

	flag.StringVar(&s3URI, "src", "", "Source S3 URI.")
	flag.StringVar(&fileName, "dst", "", "Destination output path.")
	flag.StringVar(&endpoint, "endpoint", "", `AWS endpoint to use (for connecting to S3 compatible interfaces).
Alternatively this can be set via environment variable 'S3_GET_ENDPOINT'.

If this is not set, the default AWS service endpoint will be used.

The value set by this flag will take precedence over the environment variable.`,
	)
	flag.DurationVar(&timeout, "timeout", 0, `Download timeout.
If this is set to zero, no timeout will be set.

A duration string is a possibly signed sequence of decimal numbers,
each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	`)
	flag.BoolVar(&quiet, "quiet", true, "Silence error output.")
	flag.Int64Var(&partSize, "part-size", 0, `The size (in bytes) to request from S3 for each part.
The minimum allowed part size is 5MB, and  if this value is set to zero,
the DefaultDownloadPartSize value will be used.

PartSize is ignored if the Range input parameter is provided.`)
	flag.IntVar(&concurrency, "concurrency", 0, `The number of goroutines to spin up in parallel when sending parts.
If this is set to zero, the DefaultDownloadConcurrency value will be used.

Concurrency of 1 will download the parts sequentially.

Concurrency is ignored if the Range input parameter is provided.`)
	flag.Parse()

	logger := log.New(os.Stdout, "", 0)

	if quiet {
		logger.SetOutput(io.Discard)
	}

	if s3URI == "" {
		logger.Fatal("source S3 URI not set")
	}

	if fileName == "" {
		logger.Fatal("destination output path not set")
	}

	bucket, key, err := ParseS3URI(s3URI)
	if err != nil {
		logger.Fatal(err)
	}

	if _, err = os.Stat(fileName); !os.IsNotExist(err) {
		logger.Fatalf("file already exists: %s\n", fileName)
	}

	if endpoint == "" {
		endpoint = os.Getenv("S3_GET_ENDPOINT")
	}

	client, err := newS3Client(endpoint)
	if err != nil {
		logger.Fatal(err)
	}

	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.PartSize = partSize
		d.Concurrency = concurrency
	})
	bytesDownloaded, err := Download(downloader, timeout, bucket, key, fileName)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("downloaded %dB: %v to %v\n", bytesDownloaded, s3URI, fileName)
}
