package decoder

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type Decoder func(interface{}, io.Writer) error

func Decorate(c ...Decoder) Decoder {
	return func(in interface{}, out io.Writer) error {
		w := &bytes.Buffer{}

		for _, cb := range c {
			w = &bytes.Buffer{}
			w.Reset()
			if err := cb(in, w); err != nil {
				return errors.WithStack(err)
			}

			in = w.Bytes()
		}

		_, err := out.Write(w.Bytes())

		return errors.WithStack(err)
	}
}
