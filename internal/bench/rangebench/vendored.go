// Code that's been inlined from the base project for point-in-time testing.
package rangebench2

type (
	opt[T any] struct {
		value   T
		present bool
	}

	iter[T any] interface {
		next() opt[T]
	}

	stream[T any] struct {
		src iter[T]
	}

	singedInteger interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}

	unsignedInteger interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}

	integer interface {
		singedInteger | unsignedInteger
	}

	iterFn[T any] struct {
		fn func() opt[T]
	}
)

func (o opt[T]) IsVal() bool {
	return o.present
}

func (o opt[T]) Val() T {
	if !o.present {
		// todo [bs]: let's use a standard null pointer exception type here.
		// Want that in a few different places.
		panic("'Get' called on empty optional")
	}
	return o.value
}

func IterEach[T any](iter iter[T], fn func(v T)) {
	// todo - consider whether this method even should exist.

	// ques [bs]: can I do a quick extended test on that?
	for {
		next := iter.next()
		if !next.IsVal() {
			return
		} else {
			fn(next.Val())
		}
	}
}

func streamEach[T any](s stream[T], fn func(v T)) {
	IterEach(s.src, fn)
}

func (s stream[V]) toList() []V {
	// todo [bs]: should add facility so streams of known size can use
	// that to seed size here.
	l := make([]V, 0)
	streamEach(s, func(v V) {
		l = append(l, v)
	})
	return l
}

func optOf[T any](val T) opt[T] {
	return opt[T]{
		value:   val,
		present: true,
	}
}

func (i *iterFn[T]) next() opt[T] {
	return i.fn()
}

func newFnStream[T any](fn func() opt[T]) stream[T] {
	return stream[T]{
		src: &iterFn[T]{
			fn: fn,
		},
	}
}
