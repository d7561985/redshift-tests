package csv

import (
	"io"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
)

func Marshal(in interface{}, out io.Writer) (err error) {
	if err := gocsv.Marshal(in, out); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
