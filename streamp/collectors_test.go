package streamp

import (
	"testing"

	"github.com/BennettJames/ef"
	"github.com/BennettJames/ef/stream"
	"github.com/stretchr/testify/assert"
)

func TestPStreamReduceInit(t *testing.T) {
	input := stream.OfVals(
		ef.PairOf("a", 1),
		ef.PairOf("a", 2),
		ef.PairOf("b", 3),
	)
	expected := map[string][]int{
		"a": ef.Slice(1, 2),
		"b": ef.Slice(3),
	}
	assert.Equal(t,
		expected,
		ReduceInit(
			input, map[string][]int{},
			addToMultiMap[string, int]))
}

func TestPStreamReduce(t *testing.T) {
	t.Run("Sum", func(t *testing.T) {
		input := stream.OfVals(
			ef.PairOf("a", 1),
			ef.PairOf("a", 2),
			ef.PairOf("b", 3),
		)
		assert.Equal(t,
			6,
			Reduce(input, func(total int, key string, val int) int {
				total += val
				return total
			}))
	})
}

func TestPStreamFind(t *testing.T) {

	t.Run("HasValue", func(t *testing.T) {
		st := stream.OfVals(
			ef.PairOf("a", 1),
			ef.PairOf("b", 2),
			ef.PairOf("c", 3),
			ef.PairOf("d", 4),
			ef.PairOf("e", 5),
		)
		iterCount := 0
		foundVal := Find(st, func(k string, v int) bool {
			iterCount++
			return v%3 == 0
		})
		assert.Equal(t, 3, iterCount)
		assert.Equal(t, ef.NewOptValue(ef.PairOf("c", 3)), foundVal)
	})

	t.Run("MissingValue", func(t *testing.T) {
		st := stream.OfVals(
			ef.PairOf("a", 1),
			ef.PairOf("b", 2),
			ef.PairOf("c", 3),
			ef.PairOf("d", 4),
			ef.PairOf("e", 5),
		)
		iterCount := 0
		foundVal := Find(st, func(k string, v int) bool {
			iterCount++
			return v%7 == 0
		})
		assert.Equal(t, 5, iterCount)
		assert.Equal(t, ef.Opt[ef.Pair[string, int]]{}, foundVal)
	})
}

func TestPStreamAnyMatch(t *testing.T) {
	st := stream.OfVals(
		ef.PairOf("a", 1),
		ef.PairOf("b", 2),
		ef.PairOf("c", 3),
		ef.PairOf("d", 4),
		ef.PairOf("e", 5),
	)
	iterCount := 0
	foundVal := AnyMatch(st, func(k string, v int) bool {
		iterCount++
		return v%3 == 0
	})
	assert.Equal(t, 3, iterCount)
	assert.True(t, foundVal)
}

func TestPStreamAllMatch(t *testing.T) {

	st := stream.OfVals(
		ef.PairOf("a", 2),
		ef.PairOf("b", 4),
		ef.PairOf("c", 6),
		ef.PairOf("d", 8),
		ef.PairOf("e", 10),
	)
	iterCount := 0
	foundVal := AllMatch(st, func(k string, v int) bool {
		iterCount++
		return v%2 == 0
	})
	assert.Equal(t, 5, iterCount)
	assert.True(t, foundVal)
}

func addToMultiMap[K comparable, V any](m map[K][]V, key K, val V) map[K][]V {
	// this is interesting. I think this would be a good addition, but not
	// sure now is the time?
	if existing, ok := m[key]; ok {
		m[key] = append(existing, val)
	} else {
		m[key] = ef.Slice(val)
	}
	return m
}
