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

type multiStream[T any] struct {
	streams     []Stream[T]
	streamIndex int
}

func (i *multiStream[T]) Next() Opt[T] {
	// ques [bs]: I feel like some of the internal iterator patterns I've used here
	// are a bit sloppy / inconvenient. Is that a sign of bad / weird design, or more
	// a consequence of how this is trying to abstract ugly access patterns?
	for ; i.streamIndex < len(i.streams); i.streamIndex++ {
		nextStream := i.streams[i.streamIndex]
		nextVal := nextStream.src.Next()
		if !nextVal.IsEmpty() {
			return nextVal
		}
	}
	return OptEmpty[T]()
}
