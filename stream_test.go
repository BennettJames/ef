package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	t.Run("Each", func(t *testing.T) {
		input := Slice("a", "b", "c")
		st := streamOfSlice(input)
		readValues := Slice[string]()
		st.Each(func(v string) {
			readValues = append(readValues, v)
		})
		assert.Equal(t, input, readValues)

	})

	t.Run("ToSlice", func(t *testing.T) {
		input := Slice("a", "b", "c")
		assert.Equal(t, input, streamOfSlice(input).ToSlice())
	})
}

func streamOfSlice[T any](values []T) Stream[T] {
	return Stream[T]{
		srcIter: &SliceIter[T]{
			Vals: values,
		},
	}
}
