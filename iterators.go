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
		// note [bs]: I'm not 100% sold on this. I'm tempted to just have a (T, bool)
		// return value instead.
		fn func() Opt[T]
	}

	multiStream[T any] struct {
		streams     []Stream[T]
		streamIndex int
	}
)

func (si *sliceIter[T]) iterate(opFn func(val T) (advance bool)) {
	for _, v := range si.vals {
		if !opFn(v) {
			break
		}
	}
}

func (si *indexedSliceIter[T]) iterate(opFn func(Pair[int, T]) (advance bool)) {
	for i, v := range si.vals {
		advance := opFn(PairOf(i, v))
		if !advance {
			break
		}
	}
}

func (fi *fnIter[T]) iterate(opFn func(T) bool) {
	fi.fn(opFn)
}

func (fi *optFnIter[T]) iterate(opFn func(T) (advance bool)) {
	for v := fi.fn(); v.HasVal(); v = fi.fn() {
		advance := opFn(v.UnsafeGet())
		if !advance {
			break
		}
	}
}

func (ms *multiStream[T]) iterate(opFn func(T) (advance bool)) {
	for _, st := range ms.streams {
		st.srcIter.iterate(opFn)
	}
}
