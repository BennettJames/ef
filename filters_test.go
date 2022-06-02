package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptySlice(t *testing.T) {
	input := StreamOfVals(Slice(1, 2), Slice[int](), Slice(3), Slice[int]())
	withoutEmpty := StreamRemove(input, IsEmptySlice[int]).ToSlice()
	assert.Equal(t, Slice(Slice(1, 2), Slice(3)), withoutEmpty)
}

func TestIsEmptyStr(t *testing.T) {
	input := StreamOfVals("a", "b", "")
	withoutEmpty := StreamRemove(input, IsEmptyStr).ToSlice()
	assert.Equal(t, Slice("a", "b"), withoutEmpty)
}

func TestIsEmptyMap(t *testing.T) {
	input := StreamOfVals(map[int]string{}, map[int]string{1: "a"})
	withoutEmpty := StreamRemove(input, IsEmptyMap[int, string]).ToSlice()
	assert.Equal(t, Slice(map[int]string{1: "a"}), withoutEmpty)
}
