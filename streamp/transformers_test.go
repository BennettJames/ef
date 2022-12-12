package streamp

import (
	"fmt"
	"testing"

	"github.com/BennettJames/ef"
	"github.com/BennettJames/ef/stream"
	"github.com/stretchr/testify/assert"
)

func TestPStreamMap(t *testing.T) {
	input := stream.OfVals(
		ef.PairOf(1, "a"),
		ef.PairOf(2, "b"),
		ef.PairOf(3, "c"),
	)
	assert.Equal(t,
		ef.Slice(
			ef.PairOf("a", 1),
			ef.PairOf("b", 4),
			ef.PairOf("c", 9),
		),
		Map(input, func(num int, str string) (string, int) {
			return str, num * num
		}).ToSlice())
}

func TestPStreamMapKey(t *testing.T) {
	input := stream.OfVals(
		ef.PairOf(1, "a"),
		ef.PairOf(2, "b"),
		ef.PairOf(3, "c"),
	)
	assert.Equal(t,
		ef.Slice(
			ef.PairOf("1", "a"),
			ef.PairOf("4", "b"),
			ef.PairOf("9", "c"),
		),
		MapKey(input, func(num int, str string) string {
			return fmt.Sprintf("%d", num*num)
		}).ToSlice())
}

func TestPStreamMapValue(t *testing.T) {
	input := stream.OfVals(
		ef.PairOf(1, "a"),
		ef.PairOf(2, "b"),
		ef.PairOf(3, "c"),
	)
	assert.Equal(t,
		ef.Slice(
			ef.PairOf(1, "a!"),
			ef.PairOf(2, "b!"),
			ef.PairOf(3, "c!"),
		),
		MapValue(input, func(num int, str string) string {
			return str + "!"
		}).ToSlice())
}

func TestPStreamPeek(t *testing.T) {
	count := 0
	input := stream.OfVals(
		ef.PairOf(1, 10),
		ef.PairOf(2, 20),
		ef.PairOf(3, 30),
	)
	Peek(input, func(v1, v2 int) {
		count += v1 + v2
	}).ToSlice()
	assert.Equal(t, 66, count)
}

func TestPStreamKeep(t *testing.T) {
	input := stream.OfVals(ef.PairOf(1, 1), ef.PairOf(2, 3), ef.PairOf(4, 4))
	match := func(v1, v2 int) bool {
		return v1 == v2
	}
	filtered := Keep(input, match).ToSlice()
	assert.Equal(t, ef.Slice(ef.PairOf(1, 1), ef.PairOf(4, 4)), filtered)
}

func TestPStreamRemove(t *testing.T) {
	input := stream.OfVals(ef.PairOf(1, 1), ef.PairOf(2, 3), ef.PairOf(4, 4))
	match := func(v1, v2 int) bool {
		return v1 == v2
	}
	filtered := Remove(input, match).ToSlice()
	assert.Equal(t, ef.Slice(ef.PairOf(2, 3)), filtered)
}
