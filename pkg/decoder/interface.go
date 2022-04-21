package decoder

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type Decoder func(interface{}, io.Writer) error

func Decorate(c ...Decoder) (*string, Decoder) {
	header := new(string)

	return header, func(in interface{}, out io.Writer) error {
		w := &bytes.Buffer{}

		for i, mw := range c {
			w = &bytes.Buffer{}
			w.Reset()
			if err := mw(in, w); err != nil {
				return errors.WithStack(err)
			}

			in = w.Bytes()

			if i == 0 {
				buf := w.Bytes()

				if len(buf) > 0 {
					line, _, err := bufio.NewReader(bytes.NewBuffer(buf)).ReadLine()
					if err != nil {
						return errors.WithStack(err)
					}

					*header = string(line)
				}
			}
		}

		_, err := out.Write(w.Bytes())

		return errors.WithStack(err)
	}
}
