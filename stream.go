package ef

import (
	"fmt"
)

type (
	Stream[T any] struct {
		// todo [bs]: consider potentially holding an "eachable" instead - that is,
		// something that can be passed a fn that will be called on each source
		// value. Better for lists and other cases as well; can still convert an
		// iter-interface into that when it's appropriate.

		src iter[T]
	}

	iter[T any] interface {
		Next() Opt[T]
	}

	Streamable[T any] interface {
		~[]T | ~*T | Opt[T] | Stream[T] | ~func() Opt[T]
	}
)

// StreamOf creates a stream out of several types that can be converted to a
// stream.
func StreamOf[T any, S Streamable[T]](s S) Stream[T] {
	// note [bs]: for some of these, may be better to custom define an iterator
	// rather than re-using the list iterator. Also, some of the indirection here
	// is probably silly and a bit inefficient.
	switch narrowed := (interface{})(s).(type) {
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
		panic("unreachable")
	}
}

// todo [bs]: need to add some other methods, like the ability to combine
// streams.

// StreamOfSlice returns a stream of the values in the provided slice.
func StreamOfSlice[T any](values []T) Stream[T] {
	return Stream[T]{
		src: &listIter[T]{
			list: values,
		},
	}
}

// StreamOfSliceIndexed returns a stream of the values in the provided slice.
func StreamOfSliceIndexed[T any](values []T) Stream[Pair[int, T]] {
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

func NewPStream[T comparable, U any](m map[T]U) Stream[Pair[T, U]] {
	// note [bs]: this is inefficient. it has to create a list of the pairs.
	// Ideally there'd be a way to implement iter w/out having to reify the
	// full set of pairs.
	return StreamOfSlice(mapToList(m))
}

func mapToList[T comparable, U any](m map[T]U) []Pair[T, U] {
	l := make([]Pair[T, U], 0, len(m))
	for k, v := range m {
		l = append(l, PairOf(k, v))
	}
	return l
}

func (s Stream[V]) ToList() []V {
	// todo [bs]: should add facility so streams of known size can use
	// that to seed size here.
	l := make([]V, 0)
	StreamEach(s, func(v V) {
		l = append(l, v)
	})
	return l
}

func (s Stream[V]) Each(fn func(v V)) {
	StreamEach(s, fn)
}

func StreamToMap[T comparable, U any](s Stream[Pair[T, U]]) map[T]U {
	m := make(map[T]U)
	PStreamEach(s, func(t T, u U) {
		if existing, exists := m[t]; !exists {
			m[t] = u
		} else {
			panic(fmt.Sprintf(
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
	PStreamEach(s, func(t T, u U) {
		if existing, exists := m[t]; !exists {
			m[t] = u
		} else {
			m[t] = merge(t, existing, u)
		}
		m[t] = u
	})
	return m
}

func StreamMap[T, U any](s Stream[T], fn func(v T) U) Stream[U] {

	return StreamOfFn(func() Opt[U] {
		return OptMap(s.src.Next(), func(v T) U {
			return fn(v)
		})
	})
}

func PStreamMap[T, U, V, W any](
	s Stream[Pair[T, U]],
	fn func(t T, u U) (V, W),
) Stream[Pair[V, W]] {
	return StreamMap(s, func(p Pair[T, U]) Pair[V, W] {
		return PairOf(fn(p.Get()))
	})
}

func StreamPeek[T any](s Stream[T], fn func(v T)) Stream[T] {
	return StreamOfFn(func() Opt[T] {
		next := s.src.Next()
		next.IfVal(func(v T) {
			fn(v)
		})
		return next
	})
}

func PStreamPeek[T, U any](
	s Stream[Pair[T, U]],
	fn func(v T, u U),
) Stream[Pair[T, U]] {
	return StreamPeek(s, func(p Pair[T, U]) {
		fn(p.Get())
	})
}

func StreamFilter[T any](s Stream[T], fn func(v T) bool) Stream[T] {
	return StreamOfFn(func() Opt[T] {
		for {
			next := s.src.Next()
			if next.IsEmpty() || fn(next.UnsafeGet()) {
				return next
			}
		}
	})
}

func PStreamFilter[T, U any](
	s Stream[Pair[T, U]],
	fn func(t T, u U) bool,
) Stream[Pair[T, U]] {
	return StreamFilter(s, func(p Pair[T, U]) bool {
		return fn(p.Get())
	})
}

func StreamToPairs[T any, U any, V any](
	s Stream[T],
	fn func(t T) (U, V),
) Stream[Pair[U, V]] {
	// fixme - reimplement
	return Stream[Pair[U, V]]{
		// values: MapList(s.values, func(v V) Pair[U, T] {
		// 	u, t := fn(v)
		// 	return Pair[U, T]{u, t}
		// }),
	}
}

func StreamEach[T any](s Stream[T], fn func(v T)) {
	for {
		next := s.src.Next()
		if next.IsEmpty() {
			return
		}
		// note [bs]: in this case, I sorta suspect the unsafe get is probably
		// a bit more efficient and should be used instead.
		next.IfVal(fn)
		// fn(next.UnsafeGet())
	}

}

func PStreamEach[T, U any](s Stream[Pair[T, U]], fn func(t T, v U)) {
	StreamEach(s, func(p Pair[T, U]) {
		fn(p.First, p.Second)
	})
}

func StreamReduce[T any, U any](
	s Stream[T],
	merge func(total U, val T) U,
) U {
	// todo [bs]: workshop the initialization here
	u := new(U)
	return StreamReduceInit(s, *u, merge)
}

func StreamReduceInit[T any, U any](
	s Stream[T],
	initVal U,
	merge func(total U, val T) U,
) U {
	StreamEach(s, func(v T) {
		initVal = merge(initVal, v)
	})
	return initVal
}

func StreamAverage[N Number](s Stream[N]) N {
	// note [bs]: this is not mathematically sound w.r.t to overflow, among other
	// things. Let's just rip off a superior implementation like in java.

	var total N
	cnt := int64(0)

	StreamEach(s, func(v N) {
		total += v
	})

	// note [bs]: not really convinced this is legal.
	return total / N(cnt)
}

// so - let's do a bit of research into how to do average correctly.
//
// I think total simply has no guarantees about overflow, unless you go
// out of your way to check it.

// SummaryStats
//
// Note that this is not safe with overflow - if the sum exceeds the number
// type, then overflow will occur and total / average will not be accurate.
type SummaryStats[N Number] struct {
	Average  float64
	Size     int
	Total    N
	Min, Max N
}

func StreamStats[N Number](s Stream[N]) SummaryStats[N] {
	stats := SummaryStats[N]{
		Min: MaxNumber[N](),
		Max: MinNumber[N](),
	}
	StreamEach(s, func(v N) {
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

// todo [bs]: let's add a few simple things like find, any, first, etc.
