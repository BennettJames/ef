package stream

import (
	"testing"

	"github.com/BennettJames/ef"
	"github.com/stretchr/testify/assert"
)

func TestStreamToMap(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		m := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}
		assert.Equal(t, m, ToMap(OfMap(m)))
	})

	t.Run("Collision", func(t *testing.T) {
		assert.Panics(t, func() {
			ToMap(OfVals(
				ef.PairOf("a", 1),
				ef.PairOf("a", 2),
			))
		})
	})
}

func TestStreamToMapMerge(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		st := OfVals(
			ef.PairOf("a", 1),
			ef.PairOf("a", 2),
		)
		assert.Equal(t,
			map[string]int{
				"a": 2,
			},
			ToMapMerge(st, func(key string, v1, v2 int) int {
				return ef.Max(v1, v2)
			}))
	})

	t.Run("Collision", func(t *testing.T) {
		assert.Panics(t, func() {
			ToMap(OfVals(
				ef.PairOf("a", 1),
				ef.PairOf("a", 2),
			))
		})
	})
}

func TestStreamReduce(t *testing.T) {
	t.Run("Sum", func(t *testing.T) {
		st := OfVals(1, 2, 3)
		assert.Equal(t, 6, Reduce(st, ef.Add[int]))
	})

	t.Run("Max", func(t *testing.T) {
		st := OfVals(1, 2, 3)
		assert.Equal(t, 3, Reduce(st, ef.Max[int]))
	})
}

func TestStreamReduceInit(t *testing.T) {
	t.Run("Mult", func(t *testing.T) {
		st := OfVals(1, 2, 3, 4)
		assert.Equal(t, 24, ReduceInit(st, 1, ef.Mult[int]))
	})
}

func TestStreamJoinString(t *testing.T) {
	t.Run("SimpleStrings", func(t *testing.T) {
		st := OfVals("a", "b", "c")
		assert.Equal(t, "a-b-c", JoinString(st, "-"))
	})

	t.Run("NoSep", func(t *testing.T) {
		st := OfVals("a", "b", "c")
		assert.Equal(t, "abc", JoinString(st, ""))
	})

	t.Run("Empty", func(t *testing.T) {
		st := OfVals[string]()
		assert.Equal(t, "", JoinString(st, "-"))
	})

	t.Run("Single", func(t *testing.T) {
		st := OfVals("a")
		assert.Equal(t, "a", JoinString(st, "-"))
	})
}

func TestStreamFind(t *testing.T) {

	t.Run("Empty", func(t *testing.T) {
		vals := ef.Slice[int]()
		iterCount := 0
		foundVal := Find(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 0, iterCount)
		assert.Equal(t, ef.Opt[int]{}, foundVal)
	})

	t.Run("HasValue", func(t *testing.T) {
		vals := ef.Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := Find(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 3, iterCount)
		assert.Equal(t, ef.NewOptValue(3), foundVal)
	})

	t.Run("MissingValue", func(t *testing.T) {
		vals := ef.Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := Find(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%7 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.Equal(t, ef.Opt[int]{}, foundVal)
	})
}

func TestStreamAnyMatch(t *testing.T) {

	t.Run("Empty", func(t *testing.T) {
		vals := ef.Slice[int]()
		iterCount := 0
		foundVal := Match(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 0, iterCount)
		assert.False(t, foundVal)
	})

	t.Run("HasValue", func(t *testing.T) {
		vals := ef.Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := Match(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 3, iterCount)
		assert.True(t, foundVal)
	})

	t.Run("MissingValue", func(t *testing.T) {
		vals := ef.Slice(1, 2, 3, 4, 5)
		iterCount := 0
		foundVal := Match(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%7 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.False(t, foundVal)
	})
}

func TestStreamAllMatch(t *testing.T) {

	t.Run("Empty", func(t *testing.T) {
		vals := ef.Slice[int]()
		iterCount := 0
		foundVal := AllMatch(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%2 == 0
		})
		assert.Equal(t, 0, iterCount)
		assert.True(t, foundVal)
	})

	t.Run("HasValue", func(t *testing.T) {
		vals := ef.Slice(2, 4, 6, 8, 10)
		iterCount := 0
		foundVal := AllMatch(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%2 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.True(t, foundVal)
	})

	t.Run("BadMatch", func(t *testing.T) {
		vals := ef.Slice(2, 4, 6, 7, 10)
		iterCount := 0
		foundVal := AllMatch(OfSlice(vals), func(v int) bool {
			iterCount++
			return v%2 == 0
		})
		assert.Equal(t, 4, iterCount)
		assert.False(t, foundVal)
	})
}

func TestStreamStats(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		st := Empty[int]()
		assert.Equal(t,
			ef.SummaryStats[int]{
				Average: 0,
				Size:    0,
				Total:   0,
				Min:     ef.MaxNumber[int](),
				Max:     ef.MinNumber[int](),
			},
			Stats(st))
	})

	t.Run("Int", func(t *testing.T) {
		st := OfVals(1, 2, 3, 4, 5)
		assert.Equal(t,
			ef.SummaryStats[int]{
				Average: 3,
				Size:    5,
				Total:   15,
				Min:     1,
				Max:     5,
			},
			Stats(st))
	})

	t.Run("Float64", func(t *testing.T) {
		st := OfVals(1.0, 2.5, -10.0, 5.0)
		assert.Equal(t,
			ef.SummaryStats[float64]{
				Average: -0.375,
				Size:    4,
				Total:   -1.5,
				Min:     -10,
				Max:     5,
			},
			Stats(st))
	})
}
