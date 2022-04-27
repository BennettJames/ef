package ef

type Stream[V any] struct {
	src Iter[V]
}

// a proposed way of handling entry to a stream. Note I'd like to take a sec
// to consider the stream api itself here - I think the usage of values
// as an array might have run it's course.

type Streamable[V any] interface {
	[]V | *V | Opt[V] | Stream[V] | func() Opt[V]
}

func NewStream[V any, S Streamable[V]](s S) Stream[V] {
	// note [bs]: for some of these, may be better to custom define an iterator
	// rather than re-using the list iterator. Also, some of the indirection here
	// is probably silly and a bit inefficient.
	switch narrowed := (interface{})(s).(type) {
	case []V:
		return newListStream(narrowed)
	case *V:
		return newListStream(NewNullableOpt(narrowed).ToList())
	case Opt[V]:
		return newListStream(narrowed.ToList())
	case Stream[V]:
		return narrowed
	case func() Opt[V]:
		return newFnStream(narrowed)
	default:
		panic("unreachable")
	}
}

func StreamMap[V any, U any](s Stream[V], fn func(v V) U) Stream[U] {
	return newFnStream(func() Opt[U] {
		return OptMap(s.src.Next(), func(v V) U {
			return fn(v)
		})
	})
}

func StreamPeek[V any](s Stream[V], fn func(v V)) Stream[V] {
	return newFnStream(func() Opt[V] {
		next := s.src.Next()
		next.IfSet(func(v V) {
			fn(v)
		})
		return next
	})
}

func StreamFilter[V any](s Stream[V], fn func(v V) bool) Stream[V] {
	return newFnStream(func() Opt[V] {
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

func StreamToPairs[V any, U comparable, T any](
	s Stream[V],
	fn func(v V) (U, T),
) Stream[Pair[U, T]] {
	// fixme - reimplement
	return Stream[Pair[U, T]]{
		// values: MapList(s.values, func(v V) Pair[U, T] {
		// 	u, t := fn(v)
		// 	return Pair[U, T]{u, t}
		// }),
	}
}

func PStreamToMap[K comparable, V any](
	s Stream[Pair[K, V]],
) map[K]V {
	m := make(map[K]V)
	// todo [bs]: implement
	return m
}

func StreamEach[V any](s Stream[V], fn func(v V)) {
	IterEach(s.src, fn)
}

func IterEach[V any](iter Iter[V], fn func(v V)) {
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
