package ef

// coarse proposal here: whenever you
//
// many questions immediately come from that - I'll note to start that the
// src iter will _not_ be the same type as the stream (duh). So, need to
// carefully think about how to layer transforms here. Note that some targeted
// type unsafety here is fine, but if you have to bust out reflection / type
// switches where just straight logic exist.
//
// I would _guess_ that you'd want some measure of inter-pluggability in
// the transforms themselves.

type transform struct {
}

//  filter, peek, map, collect, each, reduce

type filterTransform[V any] struct {
	src Iter[V]
	fn  func(v V) bool
}

type mapTransform[V any, U any] struct {
	src Iter[V]
	fn  func(v V) U
}

type peekTransform[V any] struct {
	src Iter[V]
	fn  func(v V)
}

type reduceTransform[V any, U any] struct {
	src      Iter[V]
	fn       func(total U, val V) U
	totalVal U
}

// so - not convinced by this approach, but let's play around a bit.

// let's think a little more about reduction in particular. Is it the same as a
// collect/each in that it's terminal in of itself? Sort of.
//
// yes, it is. Yeah, I'm leaning towards thinking of transforms a little more
// statically, or even just implicitly.
//
// part of me still sorta thinks there should be a struct definition of them one
// way or another. I kinda feel like pure function composition can be a
// little... wrong for go.
//
// so then, what do I want?

// so for a simple d

// I don't _think_ that each / collect would need a transform as they
// are end states.

// It's possible I should standardize the transform to avoid the need
// for top level variance. Then you could treat the

type Stream[V any] struct {
	src Iter[V]

	// transforms []transform
}

// a proposed way of handling entry to a stream. Note I'd like to take a sec
// to consider the stream api itself here - I think the usage of values
// as an array might have run it's course.

type Streamable[V any] interface {
	[]V | *V | Opt[V] | Stream[V] | func() Opt[V]
}

func NewStream[V any, S Streamable[V]](s S) Stream[V] {
	// note [bs]: for some of these, may be better to custom define an iterator
	// rather than
	switch narrowed := (interface{})(s).(type) {
	case []V:
		return Stream[V]{
			src: &listIter[V]{
				list: narrowed,
			},
		}
	case *V:
		if narrowed != nil {
			return Stream[V]{
				src: &listIter[V]{
					list: []V{*narrowed},
				},
			}
		} else {
			return Stream[V]{
				src: &listIter[V]{},
			}
		}
	case Opt[V]:
		if narrowed.IsPresent() {
			return Stream[V]{
				src: &listIter[V]{
					list: []V{*narrowed.val},
				},
			}
		} else {
			return Stream[V]{
				src: &listIter[V]{},
			}
		}
	case Stream[V]:
		return narrowed
	case func() Opt[V]:
		return Stream[V]{
			src: &iterFn[V]{
				fn: narrowed,
			},
		}
	default:
		panic("unreachable")
	}
}

// general note - the lack of additional method parameters does kind of mean
// that you're fairly limited w/ any sort of higher-order composition. Which,
// honestly, I think is for the best. I'd still like to dink around here, but
// the limitations mean you probably can't go too crazy.
//
// A stream is still sort of an interesting one since it can be more efficient
// for processing intermediaries - let's still play around a bit with that.

func StreamMap[V any, U any](s Stream[V], fn func(v V) U) Stream[U] {
	// fixme - reimplement
	//
	// let's try to do the chained iter approach and see how it works out.
	return Stream[U]{
		src: &iterFn[U]{
			fn: func() Opt[U] {
				next := s.src.Next()
				if !next.IsPresent() {
					return Opt[U]{}
				}
				var val U = fn(next.Get())
				return Opt[U]{val: &val}
			},
		},

		// values: MapList(s.values, fn),
	}
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

// so - I'm not convinced keypair needs to exist as such. Seems like there would
// be cases were you could just restrict the pair to have a comparable key and
// that'd be it.
//
// Yeah, I'm thinking I agree. Let's kill it off for now.
//
// So - one other question about pairs. I want to support making it easy to map
// pair-to-pair. Which means not having to redo the pair declaration each time.
// I'd still need easy to-pair and from-pair options, but I think default pair
// behavior should be pair-to-pair.
//
// Another question worth mentioning: should you even have to return the key? It
// does come with problems - like if a user modifies the key. That said that's
// always a risk from collection - and it's worth emphasizing pair streaming is
// _not_ the same as key/value manipulation it's just optimized to cooperate
// with it.

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
		}
		fn(*next.val)
	}
}
