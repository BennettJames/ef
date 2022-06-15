package streamiter

type (
	Stream[T any] struct {
		// todo [bs]: consider potentially holding an "eachable" instead - that is,
		// something that can be passed a fn that will be called on each source
		// value. Better for lists and other cases as well; can still convert an
		// iter-interface into that when it's appropriate.

		src iter[T]
	}

	Stream4[T any] struct {
		forEach func(func(val T) (advance bool))
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

	Opt[T any] struct {
		value   T
		present bool
	}
)

func StreamOfSlice[T any](values []T) Stream[T] {
	return Stream[T]{
		src: &listIter[T]{
			list: values,
		},
	}
}

func (s Stream[V]) Each(fn func(v V)) {
	for next := s.src.Next(); next.present; next = s.src.Next() {
		fn(next.UnsafeGet())
	}
}

type listIter[T any] struct {
	list      []T
	nextIndex int
}

func (l *listIter[T]) Next() Opt[T] {
	if l.nextIndex >= len(l.list) {
		return Opt[T]{}
	}
	v := l.list[l.nextIndex]
	l.nextIndex++
	return Opt[T]{
		value:   v,
		present: true,
	}
}

func (o Opt[T]) UnsafeGet() T {
	if !o.present {
		// todo [bs]: I still sorta suspect that the nil error should be able
		// to contain and communicate some amount of context.
		panic("nil")
	}
	return o.value
}

func genericIter[T any](vals []T, fn func(T)) {
	for _, v := range vals {
		fn(v)
	}
}

func intIter(vals []int, fn func(int)) {
	for _, v := range vals {
		fn(v)
	}
}

func Stream4OfSlice[T any](vals []T) Stream4[T] {
	return Stream4[T]{
		forEach: func(fn func(T) bool) {
			for _, v := range vals {
				if !fn(v) {
					break
				}
			}
		},
	}
}

func (st Stream4[T]) Each(fn func(T) bool) {
	st.forEach(fn)
}

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
