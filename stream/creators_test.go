package stream

import (
	"testing"

	"github.com/BennettJames/ef"
	"github.com/stretchr/testify/assert"
)

func TestStreamOfSlice(t *testing.T) {
	t.Run("WithValues", func(t *testing.T) {
		assert.Equal(t,
			ef.Slice(1, 2, 3),
			OfVals(1, 2, 3).ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t,
			[]int{},
			OfVals[int]().ToSlice())
	})
}

func TestStreamOfIndexedSlice(t *testing.T) {
	vals := ef.Slice("a", "b", "c")
	st := OfIndexedSlice(vals)
	asList := st.ToSlice()
	assert.Equal(t, ef.Slice(
		ef.PairOf(0, "a"),
		ef.PairOf(1, "b"),
		ef.PairOf(2, "c"),
	), asList)
}

func TestStreamEmpty(t *testing.T) {
	st := Empty[string]()
	assert.Equal(t, ef.Slice[string](), st.ToSlice())
}

func TestStreamOfMap(t *testing.T) {
	// note [bs]: this has a nondeterministic order on account of it being a map.
	// That's fine, but note that I'd be happier w/ this with a good order method.
	st := OfMap(map[int]string{
		1: "a",
	})
	assert.Equal(t, ef.Slice(
		ef.PairOf(1, "a"),
	), st.ToSlice())
}

func TestStreamConcat(t *testing.T) {

	var _ = StreamMap(Of[int](ef.Slice(1, 2, 3)), func(v int) string {
		return "a"
	})

	t.Run("Basic", func(t *testing.T) {
		assert.Equal(t,
			ef.Slice(1, 2, 3, 4, 5),
			Concat(
				Of[int](ef.Slice(1, 2, 3)),
				Of[int](ef.NewOptValue(4)),
				Of[int](ef.Ptr(5)),
			).ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, ef.Slice[int](), Concat[int]().ToSlice())
	})
}

func TestStreamOf(t *testing.T) {
	t.Run("Slice", func(t *testing.T) {
		st := Of[int](ef.Slice(1, 2, 3))
		assert.Equal(t, ef.Slice(1, 2, 3), st.ToSlice())
	})

	t.Run("Pointer", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			st := Of[int](ef.Ptr(1))
			assert.Equal(t, ef.Slice(1), st.ToSlice())
		})

		t.Run("Nil", func(t *testing.T) {
			var value *int
			st := Of[int](value)
			assert.Equal(t, ef.Slice[int](), st.ToSlice())
		})
	})

	t.Run("Optional", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			st := Of[int](ef.NewOptValue(1))
			assert.Equal(t, ef.Slice(1), st.ToSlice())
		})

		t.Run("Empty", func(t *testing.T) {
			st := Of[int](ef.Opt[int]{})
			assert.Equal(t, ef.Slice[int](), st.ToSlice())
		})
	})

	t.Run("Stream", func(t *testing.T) {
		st1 := Of[int](ef.Slice(1, 2, 3))
		st2 := Of[int](st1)
		assert.Equal(t, ef.Slice(1, 2, 3), st2.ToSlice())
	})

}
