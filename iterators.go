package ef

type (
	sliceIter[T any] struct {
		vals []T
	}

	indexedSliceIter[T any] struct {
		vals []T
	}

	mapIter[T comparable, V any] struct {
		vals map[T]V
	}

	fnIter[T any] struct {
		fn func(func(val T) (advance bool))
	}

	multiStream[T any] struct {
		streams []Stream[T]
	}
)

func (si *sliceIter[T]) iterate(opFn func(T) (advance bool)) {
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

func (si *mapIter[T, U]) iterate(opFn func(Pair[T, U]) (advance bool)) {
	for k, v := range si.vals {
		if !opFn(PairOf(k, v)) {
			break
		}
	}
}

func (fi *fnIter[T]) iterate(opFn func(T) bool) {
	fi.fn(opFn)
}

func (ms *multiStream[T]) iterate(opFn func(T) (advance bool)) {
	for _, st := range ms.streams {
		st.srcIter.iterate(opFn)
	}
}
