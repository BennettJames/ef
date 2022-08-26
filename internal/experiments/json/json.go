package ef

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Readable interface {
	~[]byte | ~string | ~*string | ~[]rune | ~func(p []byte) (n int, err error)
}

type explicitReader struct {
	fn func(p []byte) (n int, err error)
}

func (sr *explicitReader) Read(p []byte) (n int, err error) {
	return sr.fn(p)
}

func AutoReader[R Readable](r R) io.Reader {
	switch narrowed := any(r).(type) {
	case nil:
		return strings.NewReader("")
	case []byte:
		return bytes.NewReader(narrowed)
	case string:
		return strings.NewReader(narrowed)
	case *string:
		if narrowed != nil {
			return strings.NewReader(*narrowed)
		} else {
			// ques [bs]: is this redundant w/ the prior nil case?
			return strings.NewReader("")
		}
	case []rune:
		// ques [bs]: is this particularly inefficient? I kinda suspect so.
		return strings.NewReader(string(narrowed))
	case func(p []byte) (n int, err error):
		if narrowed != nil {
			fmt.Println("@@@ got reader fn; nonnil")
			return &explicitReader{
				fn: narrowed,
			}
		} else {
			fmt.Println("@@@ got reader fn; nil")
			return strings.NewReader("")
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

func ReadJSON[Out any](s io.Reader) (out *Out, err error) {
	out = new(Out)
	err = json.NewDecoder(s).Decode(out)
	return out, err
}

func ReadJSON2[Out any, In Readable](r In) (out *Out, err error) {
	out = new(Out)
	err = json.NewDecoder(AutoReader(r)).Decode(out)
	return out, err
}
