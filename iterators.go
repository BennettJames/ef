package ef

type (
	SliceIter[T any] struct {
		Vals []T
	}

	IndexedSliceIter[T any] struct {
		Vals []T
	}

	MapIter[T comparable, V any] struct {
		Vals map[T]V
	}

	FnIter[T any] struct {
		Fn func(func(val T) (advance bool))
	}

	MultiStream[T any] struct {
		Streams []Stream[T]
	}
)

func (si *SliceIter[T]) Next(opFn func(T) (advance bool)) {
	for _, v := range si.Vals {
		if !opFn(v) {
			break
		}
	}
}

func (si *IndexedSliceIter[T]) Next(opFn func(Pair[int, T]) (advance bool)) {
	for i, v := range si.Vals {
		advance := opFn(PairOf(i, v))
		if !advance {
			break
		}
	}
}

func (si *MapIter[T, U]) Next(opFn func(Pair[T, U]) (advance bool)) {
	for k, v := range si.Vals {
		if !opFn(PairOf(k, v)) {
			break
		}
	}
}

func (fi *FnIter[T]) Next(opFn func(T) bool) {
	fi.Fn(opFn)
}

func (ms *MultiStream[T]) Next(opFn func(T) (advance bool)) {
	for _, st := range ms.Streams {
		st.srcIter.Next(opFn)
	}
}
