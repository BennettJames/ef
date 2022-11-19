package json

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// note [bs]: this is an experiment to try to use automatic type mapping to make
// it easier to read things into json structs. It's a little useful, but I don't
// think it quite clears the bar - in particular, the inability to union w/
// interfaces dampens the ergonomics.

// Readable is a union of a few different types that are "close" to being
// readable.
//
// Unfortunately, due to current limitations in generics the io.Reader interface
// itself cannot be included in this set (discussion -
// https://github.com/golang/go/issues/45346#issuecomment-862505803)
type Readable interface {
	~[]byte | ~string | ~*string | ~[]rune
}

// ToReader maps a "Readable" type - that is, a type that can be treated as a
// stream of characters - to an appropriate io.Reader.
func ToReader[R Readable](r R) io.Reader {
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
	default:
		panic("unreachable")
	}
}

// Parse tries to read the input as a JSON object of the provided type.
func Parse[T any](in io.Reader) (out T, err error) {
	err = json.NewDecoder(in).Decode(&out)
	return out, err
}
