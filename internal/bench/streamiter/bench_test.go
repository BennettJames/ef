package streamiter

import (
	"math"
	"testing"
)

var streamEachVal int

func BenchmarkStreamEachExpensive(b *testing.B) {

	size := 1024
	vals := make([]int, size)
	for i := 0; i < len(vals); i++ {
		vals[i] = i
	}

	b.Run("forLoop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range vals {
				streamEachVal = int(math.Pow(float64(v), 1.5))
				// streamEachVal = v * v / 3
				// streamEachVal = v
			}
		}
	})

	b.Run("genericIter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			genericIter(vals, func(v int) {
				streamEachVal = int(math.Pow(float64(v), 1.5))
				// streamEachVal = v * v / 3
				// streamEachVal = v
			})
		}
	})

	b.Run("intIter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			intIter(vals, func(v int) {
				streamEachVal = int(math.Pow(float64(v), 1.5))
				// streamEachVal = v * v / 3
				// streamEachVal = v
			})
		}
	})

	b.Run("stream4Iter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			st := Stream4OfSlice(vals)
			st.Each(func(v int) bool {
				streamEachVal = int(math.Pow(float64(v), 1.5))
				// streamEachVal = v * v / 3
				// streamEachVal = v
				return true
			})
		}
	})

	b.Run("stream5Iter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			st := Stream5OfSlice(vals)
			st.Each(func(v int) bool {
				streamEachVal = int(math.Pow(float64(v), 1.5))
				// streamEachVal = v
				return true
			})
		}
	})

	b.Run("streamIter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			st := StreamOfSlice(vals)
			st.Each(func(v int) {
				streamEachVal = int(math.Pow(float64(v), 1.5))
				// streamEachVal = v
			})
		}
	})
}
