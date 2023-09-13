package main

import (
	"testing"
)

func TestParseS3URI(t *testing.T) {
	bucket, key, err := ParseS3URI("s3://a/b/c/d")
	if err != nil {
		t.Fatal(err)
	}

	if bucket != "a" || key != "b/c/d" {
		t.Fatal("bucket or key not correctly parsed")
	}
}

func TestParseS3URIInvalidScheme(t *testing.T) {
	_, _, err := ParseS3URI("http://test-bucket/a/b/c")
	if err.Error() != "invalid scheme http" {
		t.Fatal("scheme check failed")
	}
}
