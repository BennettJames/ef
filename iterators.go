package ef

type iterFn[T any] struct {
	fn func() Opt[T]
}

func (i *iterFn[T]) Next() Opt[T] {
	return i.fn()
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
	return OptOf(v)
}

type listIterIndexed[T any] struct {
	list      []T
	nextIndex int
}

func (l *listIterIndexed[T]) Next() Opt[Pair[int, T]] {
	index := l.nextIndex
	if index >= len(l.list) {
		return OptEmpty[Pair[int, T]]()
	}
	v := l.list[index]
	l.nextIndex++
	return OptOf(PairOf(index, v))
}
