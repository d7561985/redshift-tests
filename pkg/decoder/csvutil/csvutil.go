package csvutil

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
)

func Marshal(in interface{}, out io.Writer) (err error) {
	w := csv.NewWriter(out)
	defer w.Flush()

	encoder := csvutil.NewEncoder(w)
	encoder.Register(func(f time.Time) ([]byte, error) {
		return []byte(fmt.Sprintf("%d", f.UnixMilli())), nil
	})

	// varbite should encode as hex
	encoder.Register(func(src []byte) ([]byte, error) {
		dst := make([]byte, len(src)*4)
		idx := hex.Encode(dst, src)

		return dst[:idx], nil
	})

	encoder.Register(func(src uuid.UUID) ([]byte, error) {
		dst := make([]byte, len(src)*4)
		idx := hex.Encode(dst, src[:])

		return dst[:idx], nil
	})

	if err = encoder.Encode(in); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
