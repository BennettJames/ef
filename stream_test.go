package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamOf(t *testing.T) {
	t.Run("Slice", func(t *testing.T) {
		st := StreamOf[int]([]int{1, 2, 3})
		assert.Equal(t, []int{1, 2, 3}, st.ToList())
	})

	t.Run("Pointer", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			st := StreamOf[int](Ptr(1))
			assert.Equal(t, []int{1}, st.ToList())
		})

		t.Run("Nil", func(t *testing.T) {
			var value *int
			st := StreamOf[int](value)
			assert.Equal(t, []int{}, st.ToList())
		})
	})

	t.Run("Optional", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			st := StreamOf[int](OptOf(1))
			assert.Equal(t, []int{1}, st.ToList())
		})

		t.Run("Empty", func(t *testing.T) {
			st := StreamOf[int](OptEmpty[int]())
			assert.Equal(t, []int{}, st.ToList())
		})
	})

	t.Run("Stream", func(t *testing.T) {
		st1 := StreamOf[int]([]int{1, 2, 3})
		st2 := StreamOf[int](st1)
		assert.Equal(t, []int{1, 2, 3}, st2.ToList())
	})

	// note [bs]: now seems like a good time to dive a littler deeper into
	// certain nil behavior w/ functions.
	t.Run("Func", func(t *testing.T) {
		counter := 1
		fn := func() Opt[int] {
			if counter > 3 {
				return OptEmpty[int]()
			}
			val := counter
			counter++
			return OptOf(val)
		}
		st := StreamOf[int](fn)
		assert.Equal(t, []int{1, 2, 3}, st.ToList())
	})
}

func TestStreamOfIndexedSlice(t *testing.T) {
	vals := []string{"a", "b", "c"}
	st := StreamOfIndexedSlice(vals)
	asList := st.ToList()
	assert.Equal(t, []Pair[int, string]{
		PairOf(0, "a"),
		PairOf(1, "b"),
		PairOf(2, "c"),
	}, asList)
}

func TestStreamEmpty(t *testing.T) {
	st := StreamEmpty[string]()
	assert.Equal(t, []string{}, st.ToList())
}

func TestNewPStream(t *testing.T) {
	// note [bs]: this has a nondeterministic order on account of it being a map.
	// That's fine, but note that I'd be happier w/ this with a good order method.
	st := StreamOfMap(map[int]string{
		1: "a",
	})
	assert.Equal(t, []Pair[int, string]{
		PairOf(1, "a"),
	}, st.ToList())
}

func TestStream(t *testing.T) {

}

func TestStreamJoinString(t *testing.T) {
	t.Run("SimpleStrings", func(t *testing.T) {
		st := StreamOfSlice([]string{"a", "b", "c"})
		assert.Equal(t, "a-b-c", StreamJoinString(st, "-"))
	})

	t.Run("NoSep", func(t *testing.T) {
		st := StreamOfSlice([]string{"a", "b", "c"})
		assert.Equal(t, "abc", StreamJoinString(st, ""))
	})

	t.Run("Empty", func(t *testing.T) {
		st := StreamOfSlice([]string{})
		assert.Equal(t, "", StreamJoinString(st, "-"))
	})

	t.Run("Single", func(t *testing.T) {
		st := StreamOfSlice([]string{"a"})
		assert.Equal(t, "a", StreamJoinString(st, "-"))
	})
}
