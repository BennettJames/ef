package ef

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
