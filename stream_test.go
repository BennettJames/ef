package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	t.Run("Each", func(t *testing.T) {
		input := Slice("a", "b", "c")
		st := StreamOfSlice(input)
		readValues := Slice[string]()
		st.Each(func(v string) {
			readValues = append(readValues, v)
		})
		assert.Equal(t, input, readValues)

	})

	t.Run("ToSlice", func(t *testing.T) {
		input := Slice("a", "b", "c")
		assert.Equal(t, input, StreamOfSlice(input).ToSlice())
	})
}

func TestEach(t *testing.T) {

	// todo [bs]: let's see if there are any interesting pstream options
	// here.

	t.Run("Slice", func(t *testing.T) {
		in := Slice(1, 2, 3)
		readVals := Slice[int]()
		Each(in, func(v int) {
			readVals = append(readVals, v)
		})
		assert.Equal(t, in, readVals)
	})

	t.Run("Stream", func(t *testing.T) {
		in := Slice(1, 2, 3)
		readVals := Slice[int]()
		Each(StreamOfSlice(in), func(v int) {
			readVals = append(readVals, v)
		})
		assert.Equal(t, in, readVals)
	})
}

func TestStreamOf(t *testing.T) {
	t.Run("Slice", func(t *testing.T) {
		st := StreamOf[int](Slice(1, 2, 3))
		assert.Equal(t, Slice(1, 2, 3), st.ToSlice())
	})

	t.Run("Pointer", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			st := StreamOf[int](Ptr(1))
			assert.Equal(t, Slice(1), st.ToSlice())
		})

		t.Run("Nil", func(t *testing.T) {
			var value *int
			st := StreamOf[int](value)
			assert.Equal(t, Slice[int](), st.ToSlice())
		})
	})

	t.Run("Optional", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			st := StreamOf[int](OptOf(1))
			assert.Equal(t, Slice(1), st.ToSlice())
		})

		t.Run("Empty", func(t *testing.T) {
			st := StreamOf[int](OptEmpty[int]())
			assert.Equal(t, Slice[int](), st.ToSlice())
		})
	})

	t.Run("Stream", func(t *testing.T) {
		st1 := StreamOf[int](Slice(1, 2, 3))
		st2 := StreamOf[int](st1)
		assert.Equal(t, Slice(1, 2, 3), st2.ToSlice())
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
		assert.Equal(t, Slice(1, 2, 3), st.ToSlice())
	})
}

func TestStreamOfSlice(t *testing.T) {
	t.Run("WithValues", func(t *testing.T) {
		assert.Equal(t,
			Slice(1, 2, 3),
			StreamOfVals(1, 2, 3).ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t,
			[]int{},
			StreamOfVals[int]().ToSlice())
	})
}

func TestStreamOfIndexedSlice(t *testing.T) {
	vals := Slice("a", "b", "c")
	st := StreamOfIndexedSlice(vals)
	asList := st.ToSlice()
	assert.Equal(t, Slice(
		PairOf(0, "a"),
		PairOf(1, "b"),
		PairOf(2, "c"),
	), asList)
}

func TestStreamEmpty(t *testing.T) {
	st := StreamEmpty[string]()
	assert.Equal(t, Slice[string](), st.ToSlice())
}

func TestStreamOfMap(t *testing.T) {
	// note [bs]: this has a nondeterministic order on account of it being a map.
	// That's fine, but note that I'd be happier w/ this with a good order method.
	st := StreamOfMap(map[int]string{
		1: "a",
	})
	assert.Equal(t, Slice(
		PairOf(1, "a"),
	), st.ToSlice())
}

func TestStreamConcat(t *testing.T) {

	var _ = StreamMap(StreamOf[int](Slice(1, 2, 3)), func(v int) string {
		return "a"
	})

	// this is interesting. I'm very tempted to adopt this as the core API system,
	// along with transitioning to a "key pair" system in lieu of "pair".
	//
	// That said - constant type checking / wrapping might have some
	var _ = StreamMap2(Slice(1, 2, 3), func(v int) string {
		return "a"
	})

	t.Run("Basic", func(t *testing.T) {
		assert.Equal(t,
			Slice(1, 2, 3, 4, 5),
			StreamConcat(
				StreamOf[int](Slice(1, 2, 3)),
				StreamOf[int](OptOf(4)),
				StreamOf[int](Ptr(5)),
			).ToSlice())
	})

	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, Slice[int](), StreamConcat[int]().ToSlice())
	})
}
