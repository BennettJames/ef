package ef

type (
	Stream[T any] struct {
		// todo [bs]: consider potentially holding an "eachable" instead - that is,
		// something that can be passed a fn that will be called on each source
		// value. Better for lists and other cases as well; can still convert an
		// iter-interface into that when it's appropriate.

		srcIter iter[T]
	}

	iter[T any] interface {
		iterate(operatorFn func(val T) (advance bool))
	}

	Streamable[T any] interface {
		~[]T | ~*T | Opt[T] | Stream[T] | ~func() Opt[T]
	}
)

// StreamOf creates a stream out of several types that can be converted to a
// stream.
func StreamOf[T any, S Streamable[T]](s S) Stream[T] {
	switch narrowed := any(s).(type) {
	case []T:
		return StreamOfSlice(narrowed)
	case *T:
		return StreamOfSlice(OptOfPtr(narrowed).ToList())
	case Opt[T]:
		return StreamOfSlice(narrowed.ToList())
	case Stream[T]:
		return narrowed
	case func() Opt[T]:
		// fixme [bs]: I don't think this works certain nil patterns
		return StreamOfFn(narrowed)
	default:
		panic(&UnreachableError{})
	}
}

// StreamOfSlice returns a stream of the values in the provided slice.
func StreamOfSlice[T any](values []T) Stream[T] {
	return Stream[T]{
		srcIter: &sliceIter[T]{
			vals: values,
		},
	}
}

// StreamOfVals returns a stream consisting of all elements passed to it.
func StreamOfVals[T any](vals ...T) Stream[T] {
	return StreamOfSlice(Slice(vals...))
}

// StreamOfIndexedSlice returns a stream of the values in the provided slice.
func StreamOfIndexedSlice[T any](values []T) Stream[Pair[int, T]] {
	return Stream[Pair[int, T]]{
		srcIter: &indexedSliceIter[T]{
			vals: values,
		},
	}
}

// StreamOfFn takes a function that yields an optional, and builds a stream
// around it. The stream will repeatedly call the function until it yield an
// empty optional, at which point the source will be considered exhausted.
func StreamOfFn[T any](fnSrc func() Opt[T]) Stream[T] {
	return Stream[T]{
		srcIter: &optFnIter[T]{
			fn: fnSrc,
		},
	}
}

func StreamOfFn2[T any](iterFn func(func(T) bool)) Stream[T] {
	// I feel like I'm still struggling with the precisely correct intuition here.
	// Need to think a little more about the function composition here.
	//
	// Roughly speaking: there are two levels of functions here, though I don't
	// necessarily want to think about both all the time. Which I think may be the
	// issue - I want a better level and composition for the two types.
	//
	// One level is performing the base level iteration on the source, and the
	// next is performing
	//
	// A stream transform function is often just more-or-less layering another fn
	// into the core stream (indeed, it might be worth thinking about the data
	// structure representation for that). Sometimes you would in fact want to
	// affect the
	return Stream[T]{
		srcIter: &fnIter[T]{
			fn: iterFn,
		},
	}
}

// StreamEmpty returns an empty stream.
func StreamEmpty[T any]() Stream[T] {
	return StreamOfSlice([]T{})
}

// StreamOfMap creates a stream out of a map, where each entry in the stream is
// a pair of key->values foudn in the original map.
func StreamOfMap[T comparable, U any](m map[T]U) Stream[Pair[T, U]] {
	// note [bs]: this is inefficient. it has to create a list of the pairs.
	// Ideally there'd be a way to implement iter w/out having to reify the
	// full set of pairs.
	l := make([]Pair[T, U], 0, len(m))
	for k, v := range m {
		l = append(l, PairOf(k, v))
	}
	return StreamOfSlice(l)
}

// Each performs the provided fn on each element in the stream.
func (s Stream[V]) Each(eachOp func(V)) {
	s.srcIter.iterate(func(val V) (advance bool) {
		eachOp(val)
		return true
	})
}

// ToSlice puts every value of the stream into a slice.
func (s Stream[V]) ToSlice() []V {
	// todo [bs]: should add facility so streams of known size can use
	// that to seed size here.
	l := make([]V, 0)
	s.Each(func(v V) {
		l = append(l, v)
	})
	return l
}

// StreamConcat combines any number of streams into a single stream.
func StreamConcat[T any](srcStreams ...Stream[T]) Stream[T] {
	// todo [bs]: consider inspecting / flattening any concat-ed streams here
	return Stream[T]{
		srcIter: &multiStream[T]{
			streams: srcStreams,
		},
	}
}
