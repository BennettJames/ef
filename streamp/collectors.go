package streamp

import (
	"github.com/BennettJames/ef"
	"github.com/BennettJames/ef/stream"
)

// Reduce combines all the values in the pair-stream down to one of type `V`. A
// value of `V` type is initialized, and merge is called repeatedly on each
// element in the stream with the current value of `V`. Once the stream is
// finished, the final value is returned.
//
// Example:
//
//	st := stream.OfMap(map[string]int{
//	  "a": 1,
//	  "b": 2,
//	  "c": 3,
//	})
//	sum := streamp.Reduce(st, func(total int, k string, v int) int { return total + v })
func Reduce[T any, U any, V any](
	srcSt ef.Stream[ef.Pair[T, U]],
	reduceOp func(total V, first T, second U) V,
) V {
	v := new(V)
	return ReduceInit(srcSt, *v, reduceOp)
}

// ReduceInit combines all the values in the pair-stream down to one of type
// `V`. `merge“ is called repeatedly on each element in the stream with the
// current value of `V`, starting with the provided `initVal`. Once the stream
// is finished, the final value is returned.
//
// Example:
//
//	st := stream.OfMap(map[string]int{
//	  "a": 1,
//	  "b": 2,
//	  "c": 3,
//	})
//	product := streamp.ReduceInit(st, 1, func(total int, key string, val int) int {
//	  return total * val
//	})
func ReduceInit[T any, U any, V any](
	srcSt ef.Stream[ef.Pair[T, U]],
	initVal V,
	reduceOp func(total V, first T, second U) V,
) V {
	EachPair(srcSt, func(v1 T, v2 U) {
		initVal = reduceOp(initVal, v1, v2)
	})
	return initVal
}

// Find searches the pair-stream for a value that matches the provided find
// operator. If a value is found, then it is returned; otherwise an empty
// optional is.
func Find[T, U any](
	srcSt ef.Stream[ef.Pair[T, U]],
	findOp func(T, U) bool,
) ef.Opt[ef.Pair[T, U]] {
	return stream.Find(srcSt, func(p ef.Pair[T, U]) bool {
		return findOp(p.Get())
	})
}

// AnyMatch will return true if any element in the source pair stream passes the
// given match operator, and false otherwise.
func AnyMatch[T, U any](
	srcSt ef.Stream[ef.Pair[T, U]],
	matchOp func(T, U) bool,
) bool {
	return stream.Match(srcSt, func(p ef.Pair[T, U]) bool {
		return matchOp(p.Get())
	})
}

// AllMatch will return true if every element in the source pair stream passes
// the given match operator, and false otherwise.
func AllMatch[T, U any](
	srcSt ef.Stream[ef.Pair[T, U]],
	matchOp func(T, U) bool,
) bool {
	return stream.AllMatch(srcSt, func(p ef.Pair[T, U]) bool {
		return matchOp(p.Get())
	})
}
