package gz

import (
	"compress/gzip"
	"io"

	"github.com/pkg/errors"
)

var errNotBytes = errors.New("required only []byte in argument")

func Marshal(in interface{}, out io.Writer) (err error) {
	body, ok := in.([]byte)
	if !ok {
		return errors.WithStack(errNotBytes)
	}

	gw := gzip.NewWriter(out)
	if _, err = gw.Write(body); err != nil {
		return errors.WithStack(err)
	}

	if err = gw.Close(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
