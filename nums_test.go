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

	// todo [bs]: should audit ranges for sane overflow behavior.

	t.Run("Range", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{-1, 0, 1},
				Range(-1, 2).ToList())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{},
				Range(1, 1).ToList())
		})

		t.Run("OutOfBounds", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{},
				Range(5, 0).ToList())
		})
	})

	t.Run("RangeIncl", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{-1, 0, 1, 2},
				RangeIncl(-1, 2).ToList())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{1},
				RangeIncl(1, 1).ToList())
		})

		t.Run("OutOfBounds", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{},
				RangeIncl(5, 0).ToList())
		})
	})

	t.Run("RangeRev", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{1, 0, -1, -2},
				RangeRev(-2, 2).ToList())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{},
				RangeRev(1, 1).ToList())
		})

		t.Run("OutOfBounds", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{},
				RangeRev(5, 0).ToList())
		})
	})

	t.Run("RangeRevIncl", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{2, 1, 0, -1, -2},
				RangeRevIncl(-2, 2).ToList())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{1},
				RangeRevIncl(1, 1).ToList())
		})

		t.Run("OutOfBounds", func(t *testing.T) {
			assert.Equal(
				t,
				[]int{},
				RangeRevIncl(5, 0).ToList())
		})
	})
}

func checkEqualNum[N Number](t *testing.T, expected, actual N) {
	assert.Equal(t, expected, actual)
}

var updater int

func BenchmarkRangeEach(b *testing.B) {

	// compares the performance of using a range for a for loop
	// over ints vs a conventional for loop.

	// this might be a good time to think some about whether I can use an
	// alternate structure for streams. Having an opt-fn in an iterator isn't
	// always a bad idea, but I suspect it is radically less efficient than being
	// able to just fire up a range and let it do it's thing.
	//
	// let's also see if opt itself is adding any overhead - would doing a (val,
	// ok) be better? let's find out.
	//
	// There may be a few other cases where I'd want to

	// let's try another variant of this with

	max := 1024 * 8

	b.Run("forLoop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < max; j++ {
				updater = j
			}
		}
	})

	// b.Run("eachStreamerRange", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		eachStreamerRange(0, max).each(func(i int) {
	// 			updater = i
	// 		})
	// 	}
	// })
	b.Run("eachStreamerRangePtr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			eachStreamerRangePtr(0, max).each(func(i int) {
				updater = i
			})
		}
	})

	b.Run("streamSliceInterfacePtr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var r eachStreamer = &eachStreamerRangeStruct{
				start: 0,
				end:   max,
			}
			r.each(func(i int) {
				updater = i
			})
		}
	})

	b.Run("streamSliceGenericInterfacePtr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var r eachStreamerGeneric[int] = &eachStreamerRangeGeneric[int]{
				start: 0,
				end:   max,
			}
			r.each(func(i int) {
				updater = i
			})
		}
	})

	b.Run("rangeFuncDirectStruct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := eachStreamerRangeStruct{
				start: 0,
				end:   max,
			}
			r.each(func(i int) {
				updater = i
			})
		}
	})

	b.Run("rangeFuncPtrStruct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := &eachStreamerRangeStruct{
				start: 0,
				end:   max,
			}
			r.each(func(i int) {
				updater = i
			})
		}
	})
}

type eachStreamer interface {
	each(func(int))
}

type eachStreamerGeneric[T Number] interface {
	each(func(T))
}

func eachStreamerRange(start, end int) eachStreamer {
	return eachStreamerRangeStruct{
		start: start,
		end:   end,
	}
}

func eachStreamerRangePtr(start, end int) eachStreamer {
	return &eachStreamerRangeStruct{
		start: start,
		end:   end,
	}
}

type eachStreamerRangeStruct struct {
	start, end int
}

func (s eachStreamerRangeStruct) each(fn func(int)) {
	for i := s.start; i < s.end; i++ {
		fn(i)
	}
}

type eachStreamerRangeGeneric[T Number] struct {
	start, end T
}

func (s eachStreamerRangeGeneric[T]) each(fn func(T)) {
	for i := s.start; i < s.end; i++ {
		fn(i)
	}
}
