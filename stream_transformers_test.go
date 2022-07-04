package ef

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
