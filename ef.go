package ef

import (
	"encoding/json"
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

type Opt[V any] struct {
	val *V
}

func (o Opt[V]) Get() V {
	return *o.val
}

func (o Opt[V]) IsPresent() bool {
	// todo [bs]: look into weird golang nil behavior here
	return o.val != nil
}

type Monad[V any, E error] struct {
}

func OptToMonad[V any](o Opt[V]) Monad[V, error] {
	return Monad[V, error]{}
}

// so - I believe the following could *not* be turned into a method,
// as methods have fairly restrictive generic options - basically
// can just operate on a root.

func MapOpt[V any, T any](o Opt[V], fn func(v V) T) Opt[T] {
	// ques [bs]: does this handle interface nullability correctly?
	if o.val == nil {
		return Opt[T]{}
	}
	val := fn(*o.val)
	return Opt[T]{
		val: &val,
	}
}

// tempted to try out a full "java stream" api here. well, "full" might
// be a little aspirational, but you get my point. Just a loose way
// to kick the tires some.

type Stream[V any] struct {
	// ques [bs]: to start, could I kind of fake this and just do an
	// inner transformation each time? By that I mean just modify the
	// parameters as-is and maintain an inner set.

	values []V
}

// general note - the lack of additional method parameters does kind of mean
// that you're fairly limited w/ any sort of higher-order composition. Which,
// honestly, I think is for the best. I'd still like to dink around here, but
// the limitations mean you probably can't go too crazy.
//
// A stream is still sort of an interesting one since it can be more efficient
// for processing intermediaries - let's still play around a bit with that.

func StreamMap[V any, U any](s Stream[V], fn func(v V) U) Stream[U] {
	return Stream[U]{
		values: MapList(s.values, fn),
	}
}

// idea - auto, multisource json deserializer. Mostly as a way to experiment w/
// generic pseudo-unions that I've seen - not sure they're applicable here, but
// only one way to find out.

func StreamToPairs[V any, U comparable, T any](
	s Stream[V],
	fn func(v V) (U, T),
) PStream[U, T] {
	return PStream[U, T]{
		values: MapList(s.values, func(v V) KeyPair[U, T] {
			u, t := fn(v)
			return KeyPair[U, T]{u, t}
		}),
	}
}

type PStream[K comparable, V any] struct {
	values []KeyPair[K, V]
}

type KeyPair[K comparable, V any] struct {
	k K
	v V
}

type Pair[LeftType any, RightType any] struct {
	Left  LeftType
	Right RightType
}

func PairOf[LeftType any, RightType any](
	left LeftType,
	right RightType,
) Pair[LeftType, RightType] {
	return Pair[LeftType, RightType]{
		Left:  left,
		Right: right,
	}
}

type Iterator[V any] interface {
	Next() Opt[V]
}

func IteratorForEach[V any](iter Iterator[V], fn func(v V)) {
	for {
		next := iter.Next()
		if !next.IsPresent() {
			return
		}
		fn(*next.val)
	}
}

// func PStreamToDict[K comparable, V any](s PStream[K, V]) Dict[K, V] {}

type Dict[KeyType comparable, ValueType any] interface {
	Get(key KeyType) Opt[ValueType]
	Iter() Iterator[KeyPair[KeyType, ValueType]]
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
		fn(next.val.k, next.val.v)
	}
}

// does golang have a notion of "default interface fn"? I'd guess not -
// and honestly, it wouldn't feel entirely appropriate for my use
// case. S

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

// ques - what's the best way to do a map-based iterator atop a go map? Might
// require writing it to an intermediary first. That'd be... unfortunate.
//
// Yeah, I think that's the case. Oh well, I'd just proceed w/ a bit of relevant
// hacks then revisit some assumptions later.
//
// Let's do a detour to understand what exactly union types are like in go
// generics.

// so - at first blush I suspect the authors do not like this use case.
// which is fair - I'm not 100% convinced it ought be allowed either.

type Stringlike interface {
	[]byte | string | []rune
}

// so interestingly, the mostly-original design doc suggests the following
// as a legal constraint. I'm guessing they decided along the way to restrain
// generics further. A little limiting but oh well.
//
// This may explain the issue - https://github.com/golang/go/issues/45346#issuecomment-862505803

func JSONReader[S Stringlike, V any](s S) func() V {
	return nil
}

func ReadJSON[S Stringlike, Out any](s S) (*Out, error) {
	v := new(Out)
	// alright, so this is actually a little more generous than I originally
	// thought - seems like you need to "get away" from the base constraint
	// type, so to speak. Still seems to have problems with interfaces in
	// interfaces that's not _super_ well documented

	switch narrowed := (interface{})(s).(type) {
	case []byte:
		err := json.Unmarshal(narrowed, &v)
		return v, err
	case string:
		err := json.Unmarshal([]byte(narrowed), &v)
		return v, err
	case []rune:
		err := json.Unmarshal([]byte(string(narrowed)), &v)
		return v, err
	default:
		panic("unreachable")
	}
}

func MustReadJSON[S Stringlike, Out any](s S) *Out {
	v := new(Out)
	// alright, so this is actually a little more generous than I originally
	// thought - seems like you need to "get away" from the base constraint
	// type, so to speak. Still seems to have problems with interfaces in
	//

	switch narrowed := (interface{})(s).(type) {
	case []byte:
		json.Unmarshal(narrowed, v)
		return v
	case string:
		json.Unmarshal([]byte(narrowed), v)
		return v
	case []rune:
		json.Unmarshal([]byte(string(narrowed)), v)
		return v
	default:
		panic("unreachable")
	}
}
