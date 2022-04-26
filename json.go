package ef

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type Readable interface {
	[]byte | string | *string | []rune | func(p []byte) (n int, err error)
}

type explicitReader struct {
	fn func(p []byte) (n int, err error)
}

func (sr *explicitReader) Read(p []byte) (n int, err error) {
	return sr.fn(p)
}

func AutoReader[R Readable](r R) io.Reader {
	switch narrowed := (interface{})(r).(type) {
	case []byte:
		return bytes.NewReader(narrowed)
	case string:
		return strings.NewReader(narrowed)
	case *string:
		if narrowed != nil {
			return strings.NewReader(*narrowed)
		} else {
			return strings.NewReader("")
		}
	case []rune:
		// ques [bs]: is this particularly inefficient? I kinda suspect so.
		return strings.NewReader(string(narrowed))
	case func(p []byte) (n int, err error):
		return &explicitReader{
			fn: narrowed,
		}
	default:
		panic("unreachable")
	}
}

// so interestingly, the mostly-original design doc suggests the following as a
// legal constraint. I'm guessing they decided along the way to restrain
// generics further. A little limiting but oh well.
//
// This may explain the issue -
// https://github.com/golang/go/issues/45346#issuecomment-862505803

func JSONReader[S Readable, V any](s S) func() V {
	return nil
}

func ReadJSON[Out any](s io.Reader) (out *Out, err error) {
	out = new(Out)
	err = json.NewDecoder(s).Decode(out)
	return out, err
}
