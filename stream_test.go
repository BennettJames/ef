package ef

import (
	"fmt"
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

func TestStreamToMap(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		m := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		assert.Equal(t, m, StreamToMap(StreamOfMap(m)))
	})

	t.Run("Collision", func(t *testing.T) {
		assert.Panics(t, func() {
			StreamToMap(StreamOfVals(
				PairOf("a", 1),
				PairOf("a", 2),
			))
		})
	})
}

func TestStreamToMapMerge(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		st := StreamOfVals(
			PairOf("a", 1),
			PairOf("a", 2),
		)
		assert.Equal(t,
			map[string]int{
				"a": 2,
			},
			StreamToMapMerge(st, func(key string, v1, v2 int) int {
				return Max(v1, v2)
			}))
	})

	t.Run("Collision", func(t *testing.T) {
		assert.Panics(t, func() {
			StreamToMap(StreamOfVals(
				PairOf("a", 1),
				PairOf("a", 2),
			))
		})
	})
}

func TestStreamMap(t *testing.T) {
	input := StreamOfVals(1, 2, 3)
	assert.Equal(t,
		Slice("2", "4", "6"),
		StreamMap(input, func(v int) string {
			return fmt.Sprintf("%d", v*2)
		}).ToSlice())
}

func TestPStreamMap(t *testing.T) {
	input := StreamOfVals(
		PairOf(1, "a"),
		PairOf(2, "b"),
		PairOf(3, "c"),
	)
	assert.Equal(t,
		Slice(
			PairOf("a", 1),
			PairOf("b", 4),
			PairOf("c", 9),
		),
		PStreamMap(input, func(num int, str string) (string, int) {
			return str, num * num
		}).ToSlice())
}

func TestPStreamMapKey(t *testing.T) {
	input := StreamOfVals(
		PairOf(1, "a"),
		PairOf(2, "b"),
		PairOf(3, "c"),
	)
	assert.Equal(t,
		Slice(
			PairOf("1", "a"),
			PairOf("4", "b"),
			PairOf("9", "c"),
		),
		PStreamMapKey(input, func(num int, str string) string {
			return fmt.Sprintf("%d", num*num)
		}).ToSlice())
}

func TestPStreamMapValue(t *testing.T) {
	input := StreamOfVals(
		PairOf(1, "a"),
		PairOf(2, "b"),
		PairOf(3, "c"),
	)
	assert.Equal(t,
		Slice(
			PairOf(1, "a!"),
			PairOf(2, "b!"),
			PairOf(3, "c!"),
		),
		PStreamMapValue(input, func(num int, str string) string {
			return str + "!"
		}).ToSlice())
}

func TestStreamPeek(t *testing.T) {
	count := 0
	input := StreamOfVals(1, 2, 3)
	StreamPeek(input, func(v int) {
		count += v
	}).ToSlice()
	assert.Equal(t, 6, count)
}

func TestPStreamPeek(t *testing.T) {
	count := 0
	input := StreamOfVals(
		PairOf(1, 10),
		PairOf(2, 20),
		PairOf(3, 30),
	)
	PStreamPeek(input, func(v1, v2 int) {
		count += v1 + v2
	}).ToSlice()
	assert.Equal(t, 66, count)
}

func TestStreamKeep(t *testing.T) {
	input := StreamOfVals(0, 1, 2, 3, 4)
	isEven := func(v int) bool {
		return v%2 == 0
	}
	filtered := StreamKeep(input, isEven).ToSlice()
	assert.Equal(t, Slice(0, 2, 4), filtered)
}

func TestPStreamKeep(t *testing.T) {
	input := StreamOfVals(PairOf(1, 1), PairOf(2, 3), PairOf(4, 4))
	match := func(v1, v2 int) bool {
		return v1 == v2
	}
	filtered := PStreamKeep(input, match).ToSlice()
	assert.Equal(t, Slice(PairOf(1, 1), PairOf(4, 4)), filtered)
}

func TestStreamRemove(t *testing.T) {
	input := StreamOfVals(0, 1, 2, 3, 4)
	isEven := func(v int) bool {
		return v%2 == 0
	}
	filtered := StreamRemove(input, isEven).ToSlice()
	assert.Equal(t, Slice(1, 3), filtered)
}

func TestPStreamRemove(t *testing.T) {
	input := StreamOfVals(PairOf(1, 1), PairOf(2, 3), PairOf(4, 4))
	match := func(v1, v2 int) bool {
		return v1 == v2
	}
	filtered := PStreamRemove(input, match).ToSlice()
	assert.Equal(t, Slice(PairOf(2, 3)), filtered)
}

func TestStreamReduce(t *testing.T) {
	t.Run("Sum", func(t *testing.T) {
		st := StreamOfVals(1, 2, 3)
		assert.Equal(t, 6, StreamReduce(st, Add[int]))
	})

	t.Run("Max", func(t *testing.T) {
		st := StreamOfVals(1, 2, 3)
		assert.Equal(t, 3, StreamReduce(st, Max[int]))
	})
}

func TestStreamReduceInit(t *testing.T) {
	t.Run("Mult", func(t *testing.T) {
		st := StreamOfVals(1, 2, 3, 4)
		assert.Equal(t, 24, StreamReduceInit(st, 1, Mult[int]))
	})
}

func TestPStreamReduce(t *testing.T) {
	t.Run("Sum", func(t *testing.T) {
		input := StreamOfVals(
			PairOf("a", 1),
			PairOf("a", 2),
			PairOf("b", 3),
		)
		assert.Equal(t,
			6,
			PStreamReduce(input, func(total int, key string, val int) int {
				total += val
				return total
			}))
	})
}

func TestPStreamReduceInit(t *testing.T) {
	input := StreamOfVals(
		PairOf("a", 1),
		PairOf("a", 2),
		PairOf("b", 3),
	)
	expected := map[string][]int{
		"a": Slice(1, 2),
		"b": Slice(3),
	}
	assert.Equal(t,
		expected,
		PStreamReduceInit(
			input, map[string][]int{},
			addToMultiMap[string, int]))
}

func addToMultiMap[K comparable, V any](m map[K][]V, key K, val V) map[K][]V {
	// this is interesting. I think this would be a good addition, but not
	// sure now is the time?
	if existing, ok := m[key]; ok {
		m[key] = append(existing, val)
	} else {
		m[key] = Slice(val)
	}
	return m
}

func TestStreamJoinString(t *testing.T) {
	t.Run("SimpleStrings", func(t *testing.T) {
		st := StreamOfVals("a", "b", "c")
		assert.Equal(t, "a-b-c", StreamJoinString(st, "-"))
	})

	t.Run("NoSep", func(t *testing.T) {
		st := StreamOfVals("a", "b", "c")
		assert.Equal(t, "abc", StreamJoinString(st, ""))
	})

	t.Run("Empty", func(t *testing.T) {
		st := StreamOfVals[string]()
		assert.Equal(t, "", StreamJoinString(st, "-"))
	})

	t.Run("Single", func(t *testing.T) {
		st := StreamOfVals("a")
		assert.Equal(t, "a", StreamJoinString(st, "-"))
	})
}

func TestStreamStats(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		st := StreamEmpty[int]()
		assert.Equal(t,
			SummaryStats[int]{
				Average: 0,
				Size:    0,
				Total:   0,
				Min:     MaxNumber[int](),
				Max:     MinNumber[int](),
			},
			StreamStats(st))
	})

	t.Run("Int", func(t *testing.T) {
		st := StreamOfVals(1, 2, 3, 4, 5)
		assert.Equal(t,
			SummaryStats[int]{
				Average: 3,
				Size:    5,
				Total:   15,
				Min:     1,
				Max:     5,
			},
			StreamStats(st))
	})

	t.Run("Float64", func(t *testing.T) {
		st := StreamOfVals(1.0, 2.5, -10.0, 5.0)
		assert.Equal(t,
			SummaryStats[float64]{
				Average: -0.375,
				Size:    4,
				Total:   -1.5,
				Min:     -10,
				Max:     5,
			},
			StreamStats(st))
	})
}

func TestStreamConcat(t *testing.T) {

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
