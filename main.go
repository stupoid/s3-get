package main

import (
	"context"
	"flag"
	"fmt"
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
func newS3Client(endpoint string) (client *s3.Client, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return
	}

	if endpoint != "" {
		log.Printf("s3.Client using endpoint %v\n", endpoint)
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		client = s3.NewFromConfig(cfg)
	}
	return
}

// ParseS3URI extracts the bucket and object key from a S3 URI
// example S3 URI: s3://bucket-name/object-key
func ParseS3URI(s3URI string) (bucket, key string, err error) {
	u, err := url.Parse(s3URI)
	if err != nil {
		return
	}
	if u.Scheme != "s3" {
		err = fmt.Errorf("invalid scheme %v", u.Scheme)
		return
	}
	bucket = u.Host
	key = strings.TrimLeft(u.Path, "/")
	return
}

// Download retrieves an object from S3 and writes it to a file.
func Download(c *s3.Client, timeout time.Duration, bucket, key, fileName string) error {
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	d := manager.NewDownloader(c)
	_, err = d.Download(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		os.Remove(fileName) // cleanup empty file if fail to get Object
		return err
	}
	return nil
}

func main() {
	var s3URI, fileName, AWSEndpointURL string
	var timeout time.Duration

	flag.StringVar(&s3URI, "src", "", "Source S3 URI.")
	flag.StringVar(&fileName, "dst", "", "Destination output path.")
	flag.StringVar(&AWSEndpointURL, "endpoint", "", `Optional custom endpoint to use (for connecting to S3 compatible interfaces).
Alternatively this can be set via environment variable 'AWS_ENDPOINT_URL'.
But this flag will take precedence.`,
	)
	flag.DurationVar(&timeout, "timeout", time.Minute, "Download timeout.")
	flag.Parse()

	if s3URI == "" {
		log.Fatal("Source S3 URI not set")
	}

	if fileName == "" {
		log.Fatal("Destination output path not set")
	}

	bucket, key, err := ParseS3URI(s3URI)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = os.Stat(fileName); !os.IsNotExist(err) {
		log.Fatalf("File already exists: %s\n", fileName)
	}

	if AWSEndpointURL == "" {
		AWSEndpointURL = os.Getenv("AWS_ENDPOINT_URL")
	}

	client, err := newS3Client(AWSEndpointURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := Download(client, timeout, bucket, key, fileName); err != nil {
		log.Fatal(err)
	}
}
