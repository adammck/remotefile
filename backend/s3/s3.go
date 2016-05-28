package s3

import (
	"bytes"
	"io"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type S3 struct {
	api    s3iface.S3API
	Bucket *string
	Key    *string
}

func New(S3API s3iface.S3API, bucket string, key string) *S3 {
	return &S3{
		api:    S3API,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
}

func (s *S3) Get() (bool, io.Reader, error) {
	res, err := s.api.GetObject(&s3.GetObjectInput{
		Bucket: s.Bucket,
		Key:    s.Key,
	})

	if isNoSuchKey(err) {
		return false, &bytes.Buffer{}, nil
	}

	if err != nil {
		return false, &bytes.Buffer{}, err
	}

	return true, res.Body, nil
}

func (s *S3) Put(r io.ReadSeeker) error {
	_, err := s.api.PutObject(&s3.PutObjectInput{
		Bucket:               s.Bucket,
		Key:                  s.Key,
		ServerSideEncryption: aws.String("AES256"),
		Body:                 r,
	})

	return err
}

func (s *S3) Delete() error {
	_, err := s.api.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.Bucket,
		Key:    s.Key,
	})

	return err
}

func (s *S3) Filename() string {
	return path.Base(aws.StringValue(s.Key))
}

// isNoSuchKey returns true if the given error is NoSuchKey.
// See: http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#ErrorCodeList
func isNoSuchKey(err error) bool {
	if err == nil {
		return false
	}

	awsErr, ok := err.(awserr.Error)
	if !ok {
		return false
	}

	return awsErr.Code() == "NoSuchKey"
}
