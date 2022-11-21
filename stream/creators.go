package stream

import "github.com/BennettJames/ef"

// Of creates a stream out of several types that can be converted to a stream.
func Of[T any, S ef.Streamable[T]](s S) ef.Stream[T] {
	switch narrowed := any(s).(type) {
	case []T:
		return OfSlice(narrowed)
	case *T:
		if narrowed == nil {
			return OfSlice([]T{})
		} else {
			return OfSlice([]T{*narrowed})
		}
	case ef.Opt[T]:
		return OfSlice(narrowed.ToList())
	case ef.Stream[T]:
		return narrowed
	default:
		panic(&ef.UnreachableError{})
	}
}

// OfSlice returns a stream of the values in the provided slice.
func OfSlice[T any](values []T) ef.Stream[T] {
	return ef.NewStream[T](&ef.SliceIter[T]{
		Vals: values,
	})
}

// OfVals returns a stream consisting of all elements passed to it.
func OfVals[T any](vals ...T) ef.Stream[T] {
	return OfSlice(ef.Slice(vals...))
}

// OfIndexedSlice returns a stream of the values in the provided slice.
func OfIndexedSlice[T any](values []T) ef.Stream[ef.Pair[int, T]] {
	return ef.NewStream[ef.Pair[int, T]](&ef.IndexedSliceIter[T]{
		Vals: values,
	})
}

func OfFn[T any](iterFn func(func(T) bool)) ef.Stream[T] {
	return ef.NewStream[T](&ef.FnIter[T]{
		Fn: iterFn,
	})
}

// Empty returns an empty stream.
func Empty[T any]() ef.Stream[T] {
	return OfSlice([]T{})
}

// OfMap creates a stream out of a map, where each entry in the stream is
// a pair of key->values foudn in the original map.
func OfMap[T comparable, U any](m map[T]U) ef.Stream[ef.Pair[T, U]] {
	return ef.NewStream[ef.Pair[T, U]](&ef.MapIter[T, U]{
		Vals: m,
	})
}

// Concat combines any number of streams into a single stream.
func Concat[T any](srcStreams ...ef.Stream[T]) ef.Stream[T] {
	return ef.NewStream[T](&ef.MultiStream[T]{
		Streams: srcStreams,
	})
}
