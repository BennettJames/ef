package ef

type (
	Stream[T any] struct {
		srcIter Iter[T]
	}

	// Iter represents something that can repeatedly yield values until
	// exhaustion. To obtain values, pass a function that will take a value,
	// until it returns false.
	Iter[T any] interface {
		Next(operatorFn func(val T) (advance bool))
	}

	// Streamable represents something that resembles a stream, and thus can be
	// easily converted to one.
	Streamable[T any] interface {
		~[]T | ~*T | Opt[T] | Stream[T]
	}

	// SummaryStats contains a set of data about the values in a stream of numbers.
	//
	// Note that this is not safe with overflow - if the sum exceeds the number
	// type, then overflow will occur and total / average will not be accurate.
	SummaryStats[N Number] struct {
		Average  float64
		Size     int
		Total    N
		Min, Max N
	}
)

// Creates a new stream who's source is provided by the given iterator.
func NewStream[T any](iter Iter[T]) Stream[T] {
	return Stream[T]{
		srcIter: iter,
	}
}

// Each performs the provided fn on each element in the stream.
func (s Stream[V]) Each(eachOp func(V)) {
	s.srcIter.Next(func(val V) (advance bool) {
		eachOp(val)
		return true
	})
}

// ExitableEach performs the provided fn on each element in the stream, but will
// exit early and stop iteration if the operator returns false.
func (s Stream[V]) ExitableEach(eachOp func(V) bool) {
	s.srcIter.Next(func(val V) (advance bool) {
		return eachOp(val)
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
