package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestStreamFind(t *testing.T) {

	t.Run("Empty", func(t *testing.T) {
		vals := Slice[int]()
		iterCount := 0
		foundVal := StreamFind(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 0, iterCount)
		assert.Equal(t, OptEmpty[int](), foundVal)
	})

	t.Run("HasValue", func(t *testing.T) {
		vals := Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := StreamFind(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 3, iterCount)
		assert.Equal(t, OptOf(3), foundVal)
	})

	t.Run("MissingValue", func(t *testing.T) {
		vals := Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := StreamFind(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%7 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.Equal(t, OptEmpty[int](), foundVal)
	})
}

func TestPStreamFind(t *testing.T) {

	t.Run("HasValue", func(t *testing.T) {
		st := StreamOfVals(
			PairOf("a", 1),
			PairOf("b", 2),
			PairOf("c", 3),
			PairOf("d", 4),
			PairOf("e", 5),
		)
		iterCount := 0
		foundVal := PStreamFind(st, func(k string, v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 3, iterCount)
		assert.Equal(t, OptOf(PairOf("c", 3)), foundVal)
	})

	t.Run("MissingValue", func(t *testing.T) {
		st := StreamOfVals(
			PairOf("a", 1),
			PairOf("b", 2),
			PairOf("c", 3),
			PairOf("d", 4),
			PairOf("e", 5),
		)
		iterCount := 0
		foundVal := PStreamFind(st, func(k string, v int) bool {
			iterCount++
			return v%7 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.Equal(t, OptEmpty[Pair[string, int]](), foundVal)
	})
}

func TestStreamAnyMatch(t *testing.T) {

	t.Run("Empty", func(t *testing.T) {
		vals := Slice[int]()
		iterCount := 0
		foundVal := StreamAnyMatch(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 0, iterCount)
		assert.False(t, foundVal)
	})

	t.Run("HasValue", func(t *testing.T) {
		vals := Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := StreamAnyMatch(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 3, iterCount)
		assert.True(t, foundVal)
	})

	t.Run("MissingValue", func(t *testing.T) {
		vals := Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := StreamAnyMatch(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%7 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.False(t, foundVal)
	})

	t.Run("ForPairs", func(t *testing.T) {
		st := StreamOfVals(
			PairOf("a", 1),
			PairOf("b", 2),
			PairOf("c", 3),
			PairOf("d", 4),
			PairOf("e", 5),
		)
		iterCount := 0
		foundVal := PStreamAnyMatch(st, func(k string, v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 3, iterCount)
		assert.True(t, foundVal)
	})
}

func TestStreamAllMatch(t *testing.T) {

	t.Run("Empty", func(t *testing.T) {
		vals := Slice[int]()
		iterCount := 0
		foundVal := StreamAllMatch(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%2 == 0
		})
		assert.Equal(t, 0, iterCount)
		assert.True(t, foundVal)
	})

	t.Run("HasValue", func(t *testing.T) {
		vals := Slice(2, 4, 6, 8, 10)
		iterCount := 0
		foundVal := StreamAllMatch(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%2 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.True(t, foundVal)
	})

	t.Run("BadMatch", func(t *testing.T) {
		vals := Slice(2, 4, 6, 7, 10)
		iterCount := 0
		foundVal := StreamAllMatch(StreamOfSlice(vals), func(v int) bool {
			iterCount++
			return v%2 == 0
		})
		assert.Equal(t, 4, iterCount)
		assert.False(t, foundVal)
	})

	t.Run("ForPairs", func(t *testing.T) {
		st := StreamOfVals(
			PairOf("a", 2),
			PairOf("b", 4),
			PairOf("c", 6),
			PairOf("d", 8),
			PairOf("e", 10),
		)
		iterCount := 0
		foundVal := PStreamAllMatch(st, func(k string, v int) bool {
			iterCount++
			return v%2 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.True(t, foundVal)
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
