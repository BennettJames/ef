package ef

import (
	"fmt"
)

type (
	Iter[T any] interface {
		Next() Opt[T]
	}

	Pair[T1 any, T2 any] struct {
		First  T1
		Second T2
	}
)

// MapKeys returns a slice of all the keys in m.
// The keys are not returned in any particular order.
func MapKeys[Key comparable, Val any](m map[Key]Val) []Key {
	s := make([]Key, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

// MapList applies a function to every item in a given list, and
// returns the combined
func MapList[T any, U any](
	input []T,
	fn func(T) U,
) []U {
	ret := make([]U, len(input))
	for i, v := range input {
		ret[i] = fn(v)
	}
	return ret
}

func MapMap[T1 comparable, U1 any, T2 comparable, U2 any](
	input map[T1]U1,
	fn func(k T1, v U1) (k2 T2, v2 U2),
) map[T2]U2 {
	ret := make(map[T2]U2, len(input))
	for k1, v1 := range input {
		k2, v2 := fn(k1, v1)
		ret[k2] = v2
	}
	return ret
}

func NewPair[T1 any, T2 any](
	left T1,
	right T2,
) Pair[T1, T2] {
	return Pair[T1, T2]{
		First:  left,
		Second: right,
	}
}

func (p Pair[T1, T2]) Unpack() (T1, T2) {
	return p.First, p.Second
}

func (p Pair[T1, T2]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

type iterFn[T any] struct {
	fn func() Opt[T]
}

func newFnStream[T any](fn func() Opt[T]) Stream[T] {
	return Stream[T]{
		src: &iterFn[T]{
			fn: fn,
		},
	}
}

func (i *iterFn[T]) Next() Opt[T] {
	return i.fn()
}

type listIter[T any] struct {
	list      []T
	nextIndex int
}

func newListStream[T any](values []T) Stream[T] {
	return Stream[T]{
		src: &listIter[T]{
			list: values,
		},
	}
}

func (l *listIter[T]) Next() Opt[T] {
	if l.nextIndex >= len(l.list) {
		return Opt[T]{}
	}
	v := l.list[l.nextIndex]
	l.nextIndex++
	return NewOpt(v)
}
