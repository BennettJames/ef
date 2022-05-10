package ef

import (
	"fmt"
	"math"
)

// ques [bs]: do I want stream to be a concrete or abstract type?

type Stream[T any] struct {
	src Iter[T]
}

type Streamable[T any] interface {
	~[]T | ~*T | Opt[T] | Stream[T] | ~func() Opt[T]
}

func StreamOf[T any, S Streamable[T]](s S) Stream[T] {
	// note [bs]: for some of these, may be better to custom define an iterator
	// rather than re-using the list iterator. Also, some of the indirection here
	// is probably silly and a bit inefficient.
	switch narrowed := (interface{})(s).(type) {
	case []T:
		return newListStream(narrowed)
	case *T:
		return newListStream(OptOfPtr(narrowed).ToList())
	case Opt[T]:
		return newListStream(narrowed.ToList())
	case Stream[T]:
		return narrowed
	case func() Opt[T]:
		return newFnStream(narrowed)
	default:
		panic("unreachable")
	}
}

func NewPStream[T comparable, U any](m map[T]U) Stream[Pair[T, U]] {
	// note [bs]: this is inefficient. it has to create a list of the pairs.
	// Ideally there'd be a way to implement iter w/out having to reify the
	// full set of pairs.
	return newListStream(mapToList(m))
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
	return newFnStream(func() Opt[U] {
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
	return newFnStream(func() Opt[T] {
		next := s.src.Next()
		next.IfSet(func(v T) {
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
	return newFnStream(func() Opt[T] {
		// note [bs]: I know this isn't actually a risk, but this makes me a
		// little uncomfortable. Let's think about safer conditions.
		for {
			next := s.src.Next()
			if !next.IsPresent() || fn(next.Get()) {
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
	IterEach(s.src, fn)
}

func PStreamEach[T, U any](s Stream[Pair[T, U]], fn func(t T, v U)) {
	IterEach(s.src, func(p Pair[T, U]) {
		fn(p.Get())
	})
}

func IterEach[T any](iter Iter[T], fn func(v T)) {
	// todo - consider whether this method even should exist.
	for {
		next := iter.Next()
		if !next.IsPresent() {
			return
		} else {
			fn(next.Get())
		}
	}
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

type SummaryStats[N Number] struct {
	Average  float64
	Size     int
	Total    N
	Min, Max N
}

func streamStatsInt(s Stream[int]) SummaryStats[int] {
	return streamStatsInner(math.MaxInt, math.MinInt, s)
}

func streamStatsInner[N Number](min, max N, s Stream[N]) SummaryStats[N] {
	stats := SummaryStats[N]{
		Min: min,
		Max: max,
	}
	StreamEach(s, func(v N) {
		stats.Size++
		stats.Total += v
		stats.Min = Min(stats.Min, v)
		stats.Max = Max(stats.Max, v)
	})
	return stats
}
