package stream

import "github.com/BennettJames/ef"

// StreamMap transforms each value in the input stream into a new value with the
// provided function, and returns a new stream with the result.
func StreamMap[T, U any](srcSt ef.Stream[T], mapOp func(v T) U) ef.Stream[U] {
	return ef.StreamTransform(srcSt, func(val T, nextOp func(U) bool) bool {
		return nextOp(mapOp(val))
	})
}

// StreamPeek will call the function on each element in the stream, but without
// any other side effects on the stream.
func StreamPeek[T any](srcSt ef.Stream[T], peekOp func(v T)) ef.Stream[T] {
	return ef.StreamTransform(srcSt, func(val T, nextOp func(T) bool) bool {
		peekOp(val)
		return nextOp(val)
	})
}

// StreamKeep returns a stream consisting of all elements of the source stream
// that match the given check.
func StreamKeep[T any](srcSt ef.Stream[T], keepOp func(T) bool) ef.Stream[T] {
	return ef.StreamTransform(srcSt, func(val T, nextOp func(T) bool) bool {
		// note [bs]: let's think through the bool flow here. I think it's
		// ok.
		if keepOp(val) {
			return nextOp(val)
		}
		return true
	})
}

// StreamRemove returns a stream consisting of all elements of the source stream
// that do _not_ match the given check.
func StreamRemove[T any](srcSt ef.Stream[T], removeOp func(T) bool) ef.Stream[T] {
	return ef.StreamTransform(srcSt, func(val T, nextOp func(T) bool) bool {
		if !removeOp(val) {
			return nextOp(val)
		}
		return true
	})
}

// Each will perform the given function on each element of the input.
//
// Note this takes any streamable value as input - e.g. a stream or list can
// be given. Examples -
//
//	Each(Slice(1, 2, 3), func(val int) {
//	   fmt.Println("value is - ", val)
//	})
//
//	Each(StreamOfVals(1, 2, 3), func(val int) {
//	   fmt.Println("value is - ", val)
//	})
func Each[T any, S ef.Streamable[T]](srcSt S, eachOp func(T)) {
	Of[T](srcSt).Each(eachOp)
}

// EachPair will perform the given function on each element of the input pair
// stream.
func EachPair[T, U any, S ef.Streamable[ef.Pair[T, U]]](srcSt S, eachOp func(T, U)) {
	// ques [bs]: this is one of an increasing number of cases where restricting
	// pair to have a comparable would be really convenient.
	Each(srcSt, func(p ef.Pair[T, U]) {
		eachOp(p.First, p.Second)
	})
}
