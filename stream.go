package ef

import (
	"fmt"
	"strings"
)

type (
	Stream[T any] struct {
		// todo [bs]: consider potentially holding an "eachable" instead - that is,
		// something that can be passed a fn that will be called on each source
		// value. Better for lists and other cases as well; can still convert an
		// iter-interface into that when it's appropriate.

		src iter[T]
	}

	Stream5[T any] struct {
		iter stream5Iter[T]
	}

	stream5Iter[T any] interface {
		forEach(func(val T) (advance bool))
	}

	stream5Slice[T any] struct {
		vals []T
	}

	iter[T any] interface {
		Next() Opt[T]
	}

	Streamable[T any] interface {
		~[]T | ~*T | Opt[T] | Stream[T] | ~func() Opt[T]
	}
)

func Stream5OfSlice[T any](vals []T) Stream5[T] {
	return Stream5[T]{
		iter: &stream5Slice[T]{
			vals: vals,
		},
	}
}

func (st Stream5[T]) Each(fn func(T) bool) {
	st.iter.forEach(fn)
}

func (ss *stream5Slice[T]) forEach(fn func(val T) (advance bool)) {
	for _, v := range ss.vals {
		if !fn(v) {
			break
		}
	}
}

// StreamOf creates a stream out of several types that can be converted to a
// stream.
func StreamOf[T any, S Streamable[T]](s S) Stream[T] {
	switch narrowed := any(s).(type) {
	case []T:
		return StreamOfSlice(narrowed)
	case *T:
		return StreamOfSlice(OptOfPtr(narrowed).ToList())
	case Opt[T]:
		return StreamOfSlice(narrowed.ToList())
	case Stream[T]:
		return narrowed
	case func() Opt[T]:
		// fixme [bs]: I don't think this works certain nil patterns
		return StreamOfFn(narrowed)
	default:
		panic(&UnreachableError{})
	}
}

// StreamOfSlice returns a stream of the values in the provided slice.
func StreamOfSlice[T any](values []T) Stream[T] {
	return Stream[T]{
		src: &listIter[T]{
			list: values,
		},
	}
}

// StreamOfVals returns a stream consisting of all elements passed to it.
func StreamOfVals[T any](vals ...T) Stream[T] {
	return StreamOfSlice(Slice(vals...))
}

// StreamOfIndexedSlice returns a stream of the values in the provided slice.
func StreamOfIndexedSlice[T any](values []T) Stream[Pair[int, T]] {
	return Stream[Pair[int, T]]{
		src: &listIterIndexed[T]{
			list: values,
		},
	}
}

// StreamOfFn takes a function that yields an optional, and builds a stream
// around it. The stream will repeatedly call the function until it yield an
// empty optional, at which point the source will be considered exhausted.
func StreamOfFn[T any](fnSrc func() Opt[T]) Stream[T] {
	return Stream[T]{
		src: &iterFn[T]{
			fn: fnSrc,
		},
	}
}

// StreamEmpty returns an empty stream.
func StreamEmpty[T any]() Stream[T] {
	return StreamOfSlice([]T{})
}

// StreamOfMap creates a stream out of a map, where each entry in the stream is
// a pair of key->values foudn in the original map.
func StreamOfMap[T comparable, U any](m map[T]U) Stream[Pair[T, U]] {
	// note [bs]: this is inefficient. it has to create a list of the pairs.
	// Ideally there'd be a way to implement iter w/out having to reify the
	// full set of pairs.
	l := make([]Pair[T, U], 0, len(m))
	for k, v := range m {
		l = append(l, PairOf(k, v))
	}
	return StreamOfSlice(l)
}

// Each performs the provided fn on each element in the stream.
func (s Stream[V]) Each(fn func(v V)) {
	for next := s.src.Next(); !next.IsEmpty(); next = s.src.Next() {
		fn(next.UnsafeGet())
	}
}

// ToSlice puts every value of the stream into a slice.
func (s Stream[V]) ToSlice() []V {
	// todo [bs]: should add facility so streams of known size can use
	// that to seed size here.
	l := make([]V, 0)
	s.Each(func(v V) {
		l = append(l, v)
	})
	return l
}

// StreamToMap takes each value in a pair-stream, and turns it into a map where
// the keys are the first value in the pairs, and the values the second.
//
// Note that this cannot handle key collisions - if two pairs have the same `T`
// value, this will panic. Use `StreamToMapMerge` to resolve collisions.
func StreamToMap[T comparable, U any](s Stream[Pair[T, U]]) map[T]U {
	m := make(map[T]U)
	EachPair(s, func(t T, u U) {
		if existing, exists := m[t]; !exists {
			m[t] = u
		} else {
			// todo [bs]: probably want a custom error for this
			panic(fmt.Errorf(
				"StreamToMap: duplicate values found for key '%v' - ['%v', '%v']",
				t, u, existing))
		}
		m[t] = u
	})
	return m
}

// StreamToMapMerge gathers a pair stream into a map, and resolves any duplicate
// keys using the merge function to combine values.
func StreamToMapMerge[T comparable, U any](
	s Stream[Pair[T, U]],
	merge func(key T, val1, val2 U) U,
) map[T]U {
	m := make(map[T]U)
	EachPair(s, func(key T, value U) {
		if existing, exists := m[key]; !exists {
			m[key] = value
		} else {
			m[key] = merge(key, existing, value)
		}
		m[key] = value
	})
	return m
}

// StreamMap transforms each value in the input stream into a new value with
// the provided function, and returns a new stream with the result.
func StreamMap[T, U any](s Stream[T], fn func(v T) U) Stream[U] {
	return StreamOfFn(func() Opt[U] {
		return OptMap(s.src.Next(), func(v T) U {
			return fn(v)
		})
	})
}

// PStreamMap transforms both the values in each pair in the stream with the
// provided function.
func PStreamMap[K1, V1, K2, V2 any](
	s Stream[Pair[K1, V1]],
	fn func(K1, V1) (K2, V2),
) Stream[Pair[K2, V2]] {
	// ques [bs]: how convinced am I that this is the right interface for this
	// behavior? particularly - could add function to map value or map key,
	// possibly in addition to or as a replacement for this.
	return StreamMap(s, func(p Pair[K1, V1]) Pair[K2, V2] {
		return PairOf(fn(p.Get()))
	})
}

// PStreamMapValue transforms the value in each pair in the stream with the
// provided function.
func PStreamMapValue[K, V1, V2 any](
	s Stream[Pair[K, V1]],
	fn func(K, V1) V2,
) Stream[Pair[K, V2]] {
	return StreamMap(s, func(p Pair[K, V1]) Pair[K, V2] {
		return PairOf(p.First, fn(p.Get()))
	})
}

// PStreamMapKey transforms the key in each pair in the stream with the provided
// function.
func PStreamMapKey[K1, K2, V any](
	s Stream[Pair[K1, V]],
	fn func(K1, V) K2,
) Stream[Pair[K2, V]] {
	return StreamMap(s, func(p Pair[K1, V]) Pair[K2, V] {
		return PairOf(fn(p.Get()), p.Second)
	})
}

// StreamPeek will call the function on each element in the stream, but without
// any other side effects on the stream.
func StreamPeek[T any](s Stream[T], fn func(v T)) Stream[T] {
	return StreamOfFn(func() Opt[T] {
		next := s.src.Next()
		next.IfVal(func(v T) {
			fn(v)
		})
		return next
	})
}

// StreamPeek will call the function on each pair in the stream, but without any
// other side effects on the stream.
func PStreamPeek[T, U any](
	s Stream[Pair[T, U]],
	fn func(v T, u U),
) Stream[Pair[T, U]] {
	return StreamPeek(s, func(p Pair[T, U]) {
		fn(p.Get())
	})
}

// StreamKeep returns a stream consisting of all elements of the source stream
// that match the given check.
func StreamKeep[T any](src Stream[T], check func(T) bool) Stream[T] {
	return StreamOfFn(func() Opt[T] {
		for {
			next := src.src.Next()
			if next.IsEmpty() || check(next.UnsafeGet()) {
				return next
			}
		}
	})
}

// PStreamKeep returns a stream consisting of all elements of the source stream
// that match the given check.
func PStreamKeep[T, U any](
	src Stream[Pair[T, U]],
	check func(T, U) bool,
) Stream[Pair[T, U]] {
	return StreamOfFn(func() Opt[Pair[T, U]] {
		for {
			next := src.src.Next()
			if next.IsEmpty() || check(next.UnsafeGet().Get()) {
				return next
			}
		}
	})
}

// StreamRemove returns a stream consisting of all elements of the source stream
// that do _not_ match the given check.
func StreamRemove[T any](src Stream[T], check func(T) bool) Stream[T] {
	return StreamOfFn(func() Opt[T] {
		for {
			next := src.src.Next()
			if next.IsEmpty() || !check(next.UnsafeGet()) {
				return next
			}
		}
	})
}

// PStreamRemove returns a stream consisting of all elements of the source stream
// that do _not_ match the given check.
func PStreamRemove[T, U any](
	src Stream[Pair[T, U]],
	check func(T, U) bool,
) Stream[Pair[T, U]] {
	return StreamOfFn(func() Opt[Pair[T, U]] {
		for {
			next := src.src.Next()
			if next.IsEmpty() || !check(next.UnsafeGet().Get()) {
				return next
			}
		}
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
func Each[T any, S Streamable[T]](s S, fn func(v T)) {
	StreamOf[T](s).Each(fn)
}

// EachPair will perform the given function on each element of the input pair
// stream.
func EachPair[T, U any, S Streamable[Pair[T, U]]](s S, fn func(t T, u U)) {
	// ques [bs]: this is one of an increasing number of cases where restricting
	// pair to have a comparable would be really convenient.
	Each(s, func(p Pair[T, U]) {
		fn(p.First, p.Second)
	})
}

// StreamReduce combines all the values in the stream down to one of type `U`. A
// value of `U` type is initialized, and merge is called repeatedly on each
// element in the stream with the current value of `U`. Once the stream is
// finished, the final value is returned.
//
// Example:
//
//   st := StreamOfVals(1, 2, 3)
//   sum := StreamReduce(st, func(v1, v2 int) int { return v1 + v2 })
//
func StreamReduce[T any, U any](
	s Stream[T],
	merge func(total U, val T) U,
) U {
	// todo [bs]: workshop the initialization here
	u := new(U)
	return StreamReduceInit(s, *u, merge)
}

// StreamReduceInit combines all the values in the stream down to one of type
// `U`. `merge`` is called repeatedly on each element in the stream with the
// current value of `U`, starting with the provided `initVal`. Once the stream
// is finished, the final value is returned.
//
// Example:
//
//   st := StreamOfVals(1, 2, 3)
//   product := StreamReduceInit(st, 1, func(v1, v2 int) int { return v1 * v2 })
//
func StreamReduceInit[T any, U any](
	s Stream[T],
	initVal U,
	merge func(total U, val T) U,
) U {
	s.Each(func(v T) {
		initVal = merge(initVal, v)
	})
	return initVal
}

// StreamReduce combines all the values in the pair-stream down to one of type `V`. A
// value of `V` type is initialized, and merge is called repeatedly on each
// element in the stream with the current value of `V`. Once the stream is
// finished, the final value is returned.
//
// Example:
//
//   st := StreamOfMap(map[string]int{
//     "a": 1,
//     "b": 2,
//     "c": 3,
//   })
//   sum := StreamReduce(st, func(total int, k string, v int) int { return total + v })
//
func PStreamReduce[T any, U any, V any](
	s Stream[Pair[T, U]],
	merge func(total V, first T, second U) V,
) V {
	v := new(V)
	return PStreamReduceInit(s, *v, merge)
}

// PStreamReduceInit combines all the values in the pair-stream down to one of
// type `V`. `merge`` is called repeatedly on each element in the stream with
// the current value of `V`, starting with the provided `initVal`. Once the
// stream is finished, the final value is returned.
//
// Example:
//
//   st := StreamOfMap(map[string]int{
//     "a": 1,
//     "b": 2,
//     "c": 3,
//   })
//   product := PStreamReduceInit(st, 1, func(total int, key string, val int) int {
//     return total * val
//   })
//
func PStreamReduceInit[T any, U any, V any](
	s Stream[Pair[T, U]],
	initVal V,
	merge func(total V, first T, second U) V,
) V {
	EachPair(s, func(v1 T, v2 U) {
		initVal = merge(initVal, v1, v2)
	})
	return initVal
}

// SummaryStats contains a set of data about the values in a stream of numbers.
//
// Note that this is not safe with overflow - if the sum exceeds the number
// type, then overflow will occur and total / average will not be accurate.
type SummaryStats[N Number] struct {
	Average  float64
	Size     int
	Total    N
	Min, Max N
}

// StreamStats calculates the SummaryStats object for a stream of numbers.
func StreamStats[N Number](s Stream[N]) SummaryStats[N] {
	// note [bs]: possible it'd just be better to make these optionals rather than
	// have default vals for them.
	stats := SummaryStats[N]{
		Min: MaxNumber[N](),
		Max: MinNumber[N](),
	}
	s.Each(func(v N) {
		stats.Size++
		stats.Total += v
		stats.Min = Min(stats.Min, v)
		stats.Max = Max(stats.Max, v)
	})
	if stats.Size > 0 {
		stats.Average = float64(stats.Total) / float64(stats.Size)
	}
	return stats
}

// StreamJoinString combines a stream of strings to a single string, adding
// `sep` between each string.
func StreamJoinString(st Stream[string], sep string) string {
	var sb strings.Builder
	first := true
	st.Each(func(v string) {
		if !first {
			sb.WriteString(sep)
		} else {
			first = false
		}
		sb.WriteString(v)
	})
	return sb.String()
}

// StreamConcat combines any number of streams into a single stream.
func StreamConcat[T any](streams ...Stream[T]) Stream[T] {
	// todo [bs]: consider inspecting / flattening any concat-ed streams here
	return Stream[T]{
		src: &multiStream[T]{
			streams: streams,
		},
	}
}

// todo [bs]: let's add a few simple things like find, any, first, etc.
