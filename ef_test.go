package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtr(t *testing.T) {
	x := "hello"
	assert.Equal(t, &x, Ptr("hello"))
}

func TestDeRef(t *testing.T) {
	t.Run("NonNil", func(t *testing.T) {
		val := 22
		assert.Equal(t, val, DeRef(&val))
	})

	t.Run("Nil", func(t *testing.T) {
		var val *int
		assert.Equal(t, 0, DeRef(val))
	})
}

func TestAsType(t *testing.T) {
	t.Run("RightType", func(t *testing.T) {
		var x any = 22
		xAsInt := AsType[int](x)
		assert.Equal(t, 22, xAsInt)
	})

	t.Run("WrongType", func(t *testing.T) {
		var x any = 22
		assert.Panics(t, func() {
			AsType[string](x)
		})
	})
}

func TestSlice(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		expected := []Pair[string, int]{
			PairOf("a", 1),
			PairOf("b", 2),
		}
		assert.Equal(t,
			expected,
			Slice(
				PairOf("a", 1),
				PairOf("b", 2),
			))
	})

	t.Run("Empty", func(t *testing.T) {
		expected := []Pair[string, int]{}
		assert.Equal(t,
			expected,
			Slice[Pair[string, int]]())
	})
}
