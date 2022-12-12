package ef

type (
	streamTransform[T, U any] struct {
		srcStream Stream[T]
		transform func(T, func(U) bool) bool
	}
)

func (s *streamTransform[T, U]) Next(opFn func(U) bool) {
	s.srcStream.srcIter.Next(func(val T) bool {
		return s.transform(val, opFn)
	})
}

// StreamTransform is a generic helper that can be used to inject an operator in
// a stream, and allow for composition.
func StreamTransform[T, U any](
	srcSt Stream[T],
	op func(val T, nextOp func(U) bool) (advance bool),
) Stream[U] {
	return Stream[U]{
		srcIter: &streamTransform[T, U]{
			srcStream: srcSt,
			transform: op,
		},
	}
}
