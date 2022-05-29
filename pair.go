package ef

import "fmt"

type (
	// Pair is a combination of two values of any type.
	Pair[T1, T2 any] struct {
		First  T1
		Second T2
	}
)

// PairOf creates a pair of two different values.
func PairOf[T1, T2 any](
	left T1,
	right T2,
) Pair[T1, T2] {
	return Pair[T1, T2]{
		First:  left,
		Second: right,
	}
}

// Get unpacks the two values in the pair.
func (p Pair[T1, T2]) Get() (T1, T2) {
	return p.First, p.Second
}

func (p Pair[T1, T2]) String() string {
	return fmt.Sprintf("(`%v`, `%v`)", p.First, p.Second)
}
