package ef

import (
	"fmt"
)

type Iter[V any] interface {
	Next() Opt[V]
}

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
func MapList[In any, Out any](
	input []In,
	fn func(In) Out,
) []Out {
	ret := make([]Out, len(input))
	for i, v := range input {
		ret[i] = fn(v)
	}
	return ret
}

func MapMap[K1 comparable, V1 any, K2 comparable, V2 any](
	input map[K1]V1,
	fn func(k K1, v V1) (k2 K2, v2 V2),
) map[K2]V2 {
	ret := make(map[K2]V2, len(input))
	for k1, v1 := range input {
		k2, v2 := fn(k1, v1)
		ret[k2] = v2
	}
	return ret
}

type Pair[T1 any, T2 any] struct {
	First  T1
	Second T2
}

func NewPair[LeftType any, RightType any](
	left LeftType,
	right RightType,
) Pair[LeftType, RightType] {
	return Pair[LeftType, RightType]{
		First:  left,
		Second: right,
	}
}

func (p Pair[T1, T2]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}

// func PStreamToDict[K comparable, V any](s PStream[K, V]) Dict[K, V] {}

type Dict[KeyType comparable, ValueType any] interface {
	Get(key KeyType) Opt[ValueType]
	Iter() Iter[Pair[KeyType, ValueType]]
}

type MutDict[KeyType comparable, ValueType any] interface {
	Dict[KeyType, ValueType]
	// Get(key KeyType) Opt[KeyType]
	Set(key KeyType, value ValueType)
}

func DictForEach[KeyType comparable, ValueType any](
	dict Dict[KeyType, ValueType],
	fn func(key KeyType, value ValueType),
) {
	// note [bs]: not saying this would be better or worse, but this could be
	// implemented as a method on dict.
	iter := dict.Iter()
	for {
		next := iter.Next()
		if !next.IsPresent() {
			return
		}
		fn(next.val.First, next.val.Second)
	}
}

// ques - what's the best way to do a map-based iterator atop a go map? Might
// require writing it to an intermediary first. That'd be... unfortunate.
//
// Yeah, I think that's the case. Oh well, I'd just proceed w/ a bit of relevant
// hacks then revisit some assumptions later.
//
// Let's do a detour to understand what exactly union types are like in go
// generics.

type gmapDict[KeyType comparable, ValueType any] struct {
	m map[KeyType]ValueType
}

func (m *gmapDict[KeyType, ValueType]) Get(key KeyType) Opt[ValueType] {
	value, ok := m.m[key]
	if !ok {
		return Opt[ValueType]{}
	}
	return Opt[ValueType]{&value}
}

type iterFn[V any] struct {
	fn func() Opt[V]
}

func (i *iterFn[V]) Next() Opt[V] {
	return i.fn()
}

type listIter[V any] struct {
	list      []V
	nextIndex int
}

func (l *listIter[V]) Next() Opt[V] {
	if l.nextIndex >= len(l.list) {
		return Opt[V]{}
	}
	v := l.list[l.nextIndex]
	l.nextIndex++
	return Opt[V]{&v}
}
