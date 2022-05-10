package ef

import (
	"fmt"
)

type (
	Iter[T any] interface {
		Next() Opt[T]
	}

	Pair[T1, T2 any] struct {
		First  T1
		Second T2
	}
)

func PairOf[T1, T2 any](
	left T1,
	right T2,
) Pair[T1, T2] {
	return Pair[T1, T2]{
		First:  left,
		Second: right,
	}
}

func (p Pair[T1, T2]) Get() (T1, T2) {
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
	return OptOf(v)
}

// Ptr wraps the provided value as a . Mostly useful for primitives in contexts
// where you'd otherwise have to declare an extra variable.
//
// Example:
//
//     // without Ptr:
//     value := "a string value"
//     fnThatTakesAStringPointer(&value)
//
//     // with Ptr:
//     fnThatTakesAStringPointer(ef.Ptr("a string value"))
//
func Ptr[V any](val V) *V {
	return &val
}

// Deref does a "safe dereferencing" of a pointer. If the pointer points to a
// value, the value is returned; if it is null, it returns a zero-value for the
// underlying type.
//
// Example:
//
//     ef.Deref(nil)             // == ""
//     ef.Deref(ef.Ptr("hello")) // == "hello"
//
func Deref[V any](val *V) V {
	if val == nil {
		return *new(V)
	}
	return *val
}
