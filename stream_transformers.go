package ef

type (
	streamTransform[T, U any] struct {
		srcStream Stream[T]
		transform func(T, func(U) bool) bool
	}
)

func (s *streamTransform[T, U]) iterate(opFn func(U) bool) {
	s.srcStream.srcIter.iterate(func(val T) bool {
		return s.transform(val, opFn)
	})
}

// StreamTransform is a generic helper that can be used to inject an operator
// in a stream, and allow for composition.
func StreamTransform[T, U any](
	srcSt Stream[T],
	op func(val T, nextOp func(U) bool) (advance bool),
) Stream[U] {
	return Stream[U]{
		srcIter: &streamTransform[T, U]{
			srcStream: srcSt,
			transform: op,
		},
	}
}

// StreamMap transforms each value in the input stream into a new value with
// the provided function, and returns a new stream with the result.
func StreamMap[T, U any](srcSt Stream[T], mapOp func(v T) U) Stream[U] {
	return StreamTransform(srcSt, func(val T, nextOp func(U) bool) bool {
		return nextOp(mapOp(val))
	})
}

func StreamMap2[T, U any, S Streamable[T]](srcSt S, mapOp func(v T) U) Stream[U] {
	return StreamTransform(StreamOf[T](srcSt), func(val T, nextOp func(U) bool) bool {
		return nextOp(mapOp(val))
	})
}

// PStreamMap transforms both the values in each pair in the stream with the
// provided function.
func PStreamMap[K1, V1, K2, V2 any](
	srcSt Stream[Pair[K1, V1]],
	mapOp func(K1, V1) (K2, V2),
) Stream[Pair[K2, V2]] {
	// ques [bs]: how convinced am I that this is the right interface for this
	// behavior? particularly - could add function to map value or map key,
	// possibly in addition to or as a replacement for this.
	return StreamMap(srcSt, func(p Pair[K1, V1]) Pair[K2, V2] {
		return PairOf(mapOp(p.Get()))
	})
}

// PStreamMapValue transforms the value in each pair in the stream with the
// provided function.
func PStreamMapValue[K, V1, V2 any](
	srcSt Stream[Pair[K, V1]],
	mapOp func(K, V1) V2,
) Stream[Pair[K, V2]] {
	return StreamMap(srcSt, func(p Pair[K, V1]) Pair[K, V2] {
		return PairOf(p.First, mapOp(p.Get()))
	})
}

// PStreamMapKey transforms the key in each pair in the stream with the provided
// function.
func PStreamMapKey[K1, K2, V any](
	srcSt Stream[Pair[K1, V]],
	mapOp func(K1, V) K2,
) Stream[Pair[K2, V]] {
	return StreamMap(srcSt, func(p Pair[K1, V]) Pair[K2, V] {
		return PairOf(mapOp(p.Get()), p.Second)
	})
}

// StreamPeek will call the function on each element in the stream, but without
// any other side effects on the stream.
func StreamPeek[T any](srcSt Stream[T], peekOp func(v T)) Stream[T] {
	return StreamTransform(srcSt, func(val T, nextOp func(T) bool) bool {
		peekOp(val)
		return nextOp(val)
	})
}

// StreamPeek will call the function on each pair in the stream, but without any
// other side effects on the stream.
func PStreamPeek[T, U any](
	srcSt Stream[Pair[T, U]],
	peekOp func(v T, u U),
) Stream[Pair[T, U]] {
	return StreamPeek(srcSt, func(p Pair[T, U]) {
		peekOp(p.Get())
	})
}

// StreamKeep returns a stream consisting of all elements of the source stream
// that match the given check.
func StreamKeep[T any](srcSt Stream[T], keepOp func(T) bool) Stream[T] {
	return StreamTransform(srcSt, func(val T, nextOp func(T) bool) bool {
		// note [bs]: let's think through the bool flow here. I think it's
		// ok.
		if keepOp(val) {
			return nextOp(val)
		}
		return true
	})
}

// PStreamKeep returns a stream consisting of all elements of the source stream
// that match the given check.
func PStreamKeep[T, U any](
	srcSt Stream[Pair[T, U]],
	keepOp func(T, U) bool,
) Stream[Pair[T, U]] {
	return StreamKeep(srcSt, func(p Pair[T, U]) bool {
		return keepOp(p.Get())
	})
}

// StreamRemove returns a stream consisting of all elements of the source stream
// that do _not_ match the given check.
func StreamRemove[T any](srcSt Stream[T], removeOp func(T) bool) Stream[T] {
	return StreamTransform(srcSt, func(val T, nextOp func(T) bool) bool {
		if !removeOp(val) {
			return nextOp(val)
		}
		return true
	})
}

// PStreamRemove returns a stream consisting of all elements of the source stream
// that do _not_ match the given check.
func PStreamRemove[T, U any](
	srcSt Stream[Pair[T, U]],
	removeOp func(T, U) bool,
) Stream[Pair[T, U]] {
	return StreamRemove(srcSt, func(p Pair[T, U]) bool {
		return removeOp(p.Get())
	})
}

// Each will perform the given function on each element of the input.
//
// Note this takes any streamable value as input - e.g. a stream or list can
// be given. Examples -
//
//    Each(Slice(1, 2, 3), func(val int) {
//       fmt.Println("value is - ", val)
//    })
//
//    Each(StreamOfVals(1, 2, 3), func(val int) {
//       fmt.Println("value is - ", val)
//    })
//
func Each[T any, S Streamable[T]](srcSt S, eachOp func(T)) {
	StreamOf[T](srcSt).Each(eachOp)
}

// EachPair will perform the given function on each element of the input pair
// stream.
func EachPair[T, U any, S Streamable[Pair[T, U]]](srcSt S, eachOp func(T, U)) {
	// ques [bs]: this is one of an increasing number of cases where restricting
	// pair to have a comparable would be really convenient.
	Each(srcSt, func(p Pair[T, U]) {
		eachOp(p.First, p.Second)
	})
}
