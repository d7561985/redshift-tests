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
	"github.com/d7561985/redshift-test/model"
	"github.com/d7561985/redshift-test/pkg/decoder"
	"github.com/d7561985/redshift-test/pkg/decoder/csvutil"
	"github.com/d7561985/redshift-test/pkg/decoder/gz"
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

func (s *Store) PlayerInsert(ctx context.Context, p []*model.Player) (model.Copy, error) {
	return s.Bulk(ctx, "players", p)
}

func (s *Store) CasinoBetInsert(ctx context.Context, p []*model.CBet) (model.Copy, error) {
	return s.Bulk(ctx, "cb", p)
}

func (s *Store) Bulk(cxt context.Context, table string, arr interface{}) (model.Copy, error) {
	ctx, cancel := context.WithTimeout(cxt, timeout)
	defer cancel()

	path := fmt.Sprintf("%s/%d.csv.gz", table, time.Now().Unix())

	x := bytes.NewBuffer(nil)
	header, dc := decoder.Decorate(csvutil.Marshal, gz.Marshal)
	if err := dc(&arr, x); err != nil {
		return model.Copy{}, errors.WithStack(err)
	}

	// Uploads the object to S3. The Context will interrupt the request if the
	// timeout expires.
	_, err := s.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(x.Bytes()),
		//ContentEncoding: aws.String("gzip"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			return model.Copy{}, errors.WithStack(fmt.Errorf("upload canceled due to timeout: %w", err))
		}

		return model.Copy{}, errors.WithStack(fmt.Errorf("failed to upload object: %w", err))
	}

	return model.Copy{
		Path:   path,
		Table:  table,
		Fields: *header,
	}, nil
}
