package stream

import (
	"fmt"
	"testing"

	"github.com/BennettJames/ef"
	"github.com/stretchr/testify/assert"
)

func TestStreamMap(t *testing.T) {
	input := OfVals(1, 2, 3)
	assert.Equal(t,
		ef.Slice("2", "4", "6"),
		StreamMap(input, func(v int) string {
			return fmt.Sprintf("%d", v*2)
		}).ToSlice())
}

func TestStreamPeek(t *testing.T) {
	count := 0
	input := OfVals(1, 2, 3)
	StreamPeek(input, func(v int) {
		count += v
	}).ToSlice()
	assert.Equal(t, 6, count)
}

func TestStreamKeep(t *testing.T) {
	input := OfVals(0, 1, 2, 3, 4)
	isEven := func(v int) bool {
		return v%2 == 0
	}
	filtered := StreamKeep(input, isEven).ToSlice()
	assert.Equal(t, ef.Slice(0, 2, 4), filtered)
}

func TestStreamRemove(t *testing.T) {
	input := OfVals(0, 1, 2, 3, 4)
	isEven := func(v int) bool {
		return v%2 == 0
	}
	filtered := StreamRemove(input, isEven).ToSlice()
	assert.Equal(t, ef.Slice(1, 3), filtered)
}

func TestEach(t *testing.T) {

	// todo [bs]: let's see if there are any interesting pstream options
	// here.

	t.Run("Slice", func(t *testing.T) {
		in := ef.Slice(1, 2, 3)
		readVals := ef.Slice[int]()
		Each(in, func(v int) {
			readVals = append(readVals, v)
		})
		assert.Equal(t, in, readVals)
	})

	t.Run("Stream", func(t *testing.T) {
		in := ef.Slice(1, 2, 3)
		readVals := ef.Slice[int]()
		Each(OfSlice(in), func(v int) {
			readVals = append(readVals, v)
		})
		assert.Equal(t, in, readVals)
	})
}
