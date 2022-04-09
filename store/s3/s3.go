package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/d7561985/redshift-test/store/postgres"
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

func (s *Store) Bulk(cxt context.Context, journals []*postgres.Journal) error {
	ctx, cancel := context.WithTimeout(cxt, timeout)
	defer cancel()

	x := bytes.NewBuffer(nil)
	if err := json.NewEncoder(x).Encode(&journals); err != nil {
		return errors.WithStack(err)
	}

	// Uploads the object to S3. The Context will interrupt the request if the
	// timeout expires.
	p, err := s.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(time.Now().String() + ".json"),
		Body:   bytes.NewReader(x.Bytes()),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}

		os.Exit(1)
	}

	fmt.Printf("successfully uploaded file to %s/%s\n", p.GoString(), s.bucket)

	return nil
}
