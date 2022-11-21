package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptySlice(t *testing.T) {
	input := streamOfVals(Slice(1, 2), Slice[int](), Slice(3), Slice[int]())
	withoutEmpty := streamRemove(input, IsEmptySlice[int]).ToSlice()
	assert.Equal(t, Slice(Slice(1, 2), Slice(3)), withoutEmpty)
}

func TestIsEmptyStr(t *testing.T) {
	input := streamOfVals("a", "b", "")
	withoutEmpty := streamRemove(input, IsEmptyStr).ToSlice()
	assert.Equal(t, Slice("a", "b"), withoutEmpty)
}

func TestIsEmptyMap(t *testing.T) {
	input := streamOfVals(map[int]string{}, map[int]string{1: "a"})
	withoutEmpty := streamRemove(input, IsEmptyMap[int, string]).ToSlice()
	assert.Equal(t, Slice(map[int]string{1: "a"}), withoutEmpty)
}

func TestMisc(t *testing.T) {
	type testCase[T any] struct {
		fn func(T) bool

		expected bool
	}

	assert.True(t,
		And(
			Equal(22),
			Not(Equal(3)),
			Or(
				Lesser(20),
				Greater(10),
			),
			GreaterOrEqual(22),
		)(22))

	assert.True(t,
		And(
			Equal("hello"),
		)("hello"))

	assert.False(t,
		And(
			Not(Equal("hello")),
		)("hello"))

	assert.True(t,
		Or(
			Equal("hello"),
			IsEmptyStr,
		)("hello"))

	assert.True(t, LesserOrEqual(3)(1))
	assert.True(t, LesserOrEqual(3)(3))
	assert.False(t, LesserOrEqual(3)(5))

}

func streamRemove[T any](srcSt Stream[T], removeOp func(T) bool) Stream[T] {
	return StreamTransform(srcSt, func(val T, nextOp func(T) bool) bool {
		if !removeOp(val) {
			return nextOp(val)
		}
		return true
	})
}

func streamOfVals[T any](vals ...T) Stream[T] {
	return Stream[T]{
		srcIter: &SliceIter[T]{
			Vals: vals,
		},
	}
}
