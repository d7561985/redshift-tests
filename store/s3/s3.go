package s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/d7561985/redshift-test/store/postgres"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

const timeout = time.Minute

type Store struct {
	bucket string
	*s3.S3
}

func New(bucket string) *Store {
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.
	svc := s3.New(sess, &aws.Config{
		Region: aws.String("eu-central-1")},
	)

	return &Store{bucket: bucket, S3: svc}
}

func (s *Store) Bulk(cxt context.Context, journals []*postgres.Journal) (string, error) {
	ctx, cancel := context.WithTimeout(cxt, timeout)
	defer cancel()

	path := fmt.Sprintf("journal/%d.csv", time.Now().Unix())
	x := bytes.NewBuffer(nil)

	if err := gocsv.Marshal(&journals, x); err != nil {
		return "", errors.WithStack(err)
	}

	// Uploads the object to S3. The Context will interrupt the request if the
	// timeout expires.
	_, err := s.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(x.Bytes()),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			return "", errors.WithStack(fmt.Errorf("upload canceled due to timeout: %w", err))
		}

		return "", errors.WithStack(fmt.Errorf("failed to upload object: %w", err))
	}

	return path, nil
}
