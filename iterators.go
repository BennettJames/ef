package ef

type (
	sliceIter[T any] struct {
		vals []T
	}

	indexedSliceIter[T any] struct {
		vals []T
	}

	mapIter[T any] struct {
		vals []T
	}

	fnIter[T any] struct {
		fn func(func(val T) (advance bool))
	}

	optFnIter[T any] struct {
		fn func() Opt[T]
	}

	multiStream[T any] struct {
		streams     []Stream[T]
		streamIndex int
	}
)

func (si *sliceIter[T]) forEach(fn func(val T) (advance bool)) {
	for _, v := range si.vals {
		if !fn(v) {
			break
		}
	}
}

func (si *indexedSliceIter[T]) forEach(fn func(Pair[int, T]) (advance bool)) {
	for i, v := range si.vals {
		advance := fn(PairOf(i, v))
		if !advance {
			break
		}
	}
}

func (fi *fnIter[T]) forEach(fn func(T) bool) {
	fi.fn(fn)
}

func (fi *optFnIter[T]) forEach(fn func(T) (advance bool)) {
	for v := fi.fn(); v.HasVal(); v = fi.fn() {
		advance := fn(v.UnsafeGet())
		if !advance {
			break
		}
	}
}

func (ms *multiStream[T]) forEach(fn func(T) (advance bool)) {
	for _, st := range ms.streams {
		st.src.forEach(fn)
	}
}
