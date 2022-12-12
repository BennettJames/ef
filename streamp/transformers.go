package streamp

import (
	"github.com/BennettJames/ef"
	"github.com/BennettJames/ef/stream"
)

// Map transforms both the values in each pair in the stream with the provided
// function.
func Map[K1, V1, K2, V2 any](
	srcSt ef.Stream[ef.Pair[K1, V1]],
	mapOp func(K1, V1) (K2, V2),
) ef.Stream[ef.Pair[K2, V2]] {
	// ques [bs]: how convinced am I that this is the right interface for this
	// behavior? particularly - could add function to map value or map key,
	// possibly in addition to or as a replacement for this.
	return stream.StreamMap(srcSt, func(p ef.Pair[K1, V1]) ef.Pair[K2, V2] {
		return ef.PairOf(mapOp(p.Get()))
	})
}

// MapValue transforms the value in each pair in the stream with the provided
// function.
func MapValue[K, V1, V2 any](
	srcSt ef.Stream[ef.Pair[K, V1]],
	mapOp func(K, V1) V2,
) ef.Stream[ef.Pair[K, V2]] {
	return stream.StreamMap(srcSt, func(p ef.Pair[K, V1]) ef.Pair[K, V2] {
		return ef.PairOf(p.First, mapOp(p.Get()))
	})
}

// MapKey transforms the key in each pair in the stream with the provided
// function.
func MapKey[K1, K2, V any](
	srcSt ef.Stream[ef.Pair[K1, V]],
	mapOp func(K1, V) K2,
) ef.Stream[ef.Pair[K2, V]] {
	return stream.StreamMap(srcSt, func(p ef.Pair[K1, V]) ef.Pair[K2, V] {
		return ef.PairOf(mapOp(p.Get()), p.Second)
	})
}

// Peek will call the function on each pair in the stream, but without any
// other side effects on the stream.
func Peek[T, U any](
	srcSt ef.Stream[ef.Pair[T, U]],
	peekOp func(v T, u U),
) ef.Stream[ef.Pair[T, U]] {
	return stream.StreamPeek(srcSt, func(p ef.Pair[T, U]) {
		peekOp(p.Get())
	})
}

// Keep returns a stream consisting of all elements of the source stream that
// match the given check.
func Keep[T, U any](
	srcSt ef.Stream[ef.Pair[T, U]],
	keepOp func(T, U) bool,
) ef.Stream[ef.Pair[T, U]] {
	return stream.StreamKeep(srcSt, func(p ef.Pair[T, U]) bool {
		return keepOp(p.Get())
	})
}

// Remove returns a stream consisting of all elements of the source stream that
// do _not_ match the given check.
func Remove[T, U any](
	srcSt ef.Stream[ef.Pair[T, U]],
	removeOp func(T, U) bool,
) ef.Stream[ef.Pair[T, U]] {
	return stream.StreamRemove(srcSt, func(p ef.Pair[T, U]) bool {
		return removeOp(p.Get())
	})
}

// EachPair will perform the given function on each element of the input pair
// stream.
func EachPair[T, U any, S ef.Streamable[ef.Pair[T, U]]](srcSt S, eachOp func(T, U)) {
	// ques [bs]: this is one of an increasing number of cases where restricting
	// pair to have a comparable would be really convenient.
	stream.Each(srcSt, func(p ef.Pair[T, U]) {
		eachOp(p.First, p.Second)
	})
}
