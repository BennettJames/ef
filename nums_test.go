package ef

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNum(t *testing.T) {

	t.Run("Order", func(t *testing.T) {
		assert.Equal(t, PairOf(5, 22), PairOf(Order(5, 22)))
		assert.Equal(t, PairOf(5, 22), PairOf(Order(22, 5)))
	})

	t.Run("MinNumber", func(t *testing.T) {
		checkEqualNum(t, MinNumber[int](), math.MinInt)
		checkEqualNum(t, MinNumber[int8](), math.MinInt8)
		checkEqualNum(t, MinNumber[int16](), math.MinInt16)
		checkEqualNum(t, MinNumber[int32](), math.MinInt32)
		checkEqualNum(t, MinNumber[int64](), math.MinInt64)
		checkEqualNum(t, MinNumber[uint](), 0)
		checkEqualNum(t, MinNumber[uint8](), 0)
		checkEqualNum(t, MinNumber[uint16](), 0)
		checkEqualNum(t, MinNumber[uint32](), 0)
		checkEqualNum(t, MinNumber[uint64](), 0)
		checkEqualNum(t, MinNumber[uintptr](), 0)
		checkEqualNum(t, MinNumber[float32](), float32(math.Inf(-1)))
		checkEqualNum(t, MinNumber[float64](), math.Inf(-1))
	})

	t.Run("MaxNumber", func(t *testing.T) {
		checkEqualNum(t, MaxNumber[int](), math.MaxInt)
		checkEqualNum(t, MaxNumber[int8](), math.MaxInt8)
		checkEqualNum(t, MaxNumber[int16](), math.MaxInt16)
		checkEqualNum(t, MaxNumber[int32](), math.MaxInt32)
		checkEqualNum(t, MaxNumber[int64](), math.MaxInt64)
		checkEqualNum(t, MaxNumber[uint](), math.MaxUint)
		checkEqualNum(t, MaxNumber[uint8](), math.MaxUint8)
		checkEqualNum(t, MaxNumber[uint16](), math.MaxUint16)
		checkEqualNum(t, MaxNumber[uint32](), math.MaxUint32)
		checkEqualNum(t, MaxNumber[uint64](), math.MaxUint64)
		checkEqualNum(t, MaxNumber[uintptr](), math.MaxUint)
		checkEqualNum(t, MaxNumber[float32](), float32(math.Inf(0)))
		checkEqualNum(t, MaxNumber[float64](), math.Inf(0))
	})
}

func checkEqualNum[N Number](t *testing.T, expected, actual N) {
	assert.Equal(t, expected, actual)
}
