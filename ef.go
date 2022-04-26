package ef

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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

type Res[V any, E error] struct {
	val *V
	err *E
}

func NewResValue[V any](v V) Res[V, error] {
	// note [bs]: this might have some weird nested nullability issues I would
	// have to be mindful of.
	return Res[V, error]{
		val: &v,
	}
}

func NewResErr[V any, E error](e E) Res[V, E] {
	// note [bs]: this might have some weird nested nullability issues I would
	// have to be mindful of.
	return Res[V, E]{
		err: &e,
	}
}

func OptTry[V any](o Opt[V]) (*V, bool) {
	return o.val, o.IsPresent()
}

func OptToResult[V any](o Opt[V]) Res[V, error] {
	return Res[V, error]{}
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

// coarse proposal here: whenever you
//
// many questions immediately come from that - I'll note to start that the
// src iter will _not_ be the same type as the stream (duh). So, need to
// carefully think about how to layer transforms here. Note that some targeted
// type unsafety here is fine, but if you have to bust out reflection / type
// switches where just straight logic exist.
//
// I would _guess_ that you'd want some measure of inter-pluggability in
// the transforms themselves.

type transform struct {
}

//  filter, peek, map, collect, each, reduce

type filterTransform[V any] struct {
	src Iter[V]
	fn  func(v V) bool
}

type mapTransform[V any, U any] struct {
	src Iter[V]
	fn  func(v V) U
}

type peekTransform[V any] struct {
	src Iter[V]
	fn  func(v V)
}

type reduceTransform[V any, U any] struct {
	src      Iter[V]
	fn       func(total U, val V) U
	totalVal U
}

// so - not convinced by this approach, but let's play around a bit.

// let's think a little more about reduction in particular. Is it the same as a
// collect/each in that it's terminal in of itself? Sort of.
//
// yes, it is. Yeah, I'm leaning towards thinking of transforms a little more
// statically, or even just implicitly.
//
// part of me still sorta thinks there should be a struct definition of them one
// way or another. I kinda feel like pure function composition can be a
// little... wrong for go.
//
// so then, what do I want?

// so for a simple d

// I don't _think_ that each / collect would need a transform as they
// are end states.

// It's possible I should standardize the transform to avoid the need
// for top level variance. Then you could treat the

type Stream[V any] struct {
	src Iter[V]

	transforms []transform
}

// a proposed way of handling entry to a stream. Note I'd like to take a sec
// to consider the stream api itself here - I think the usage of values
// as an array might have run it's course.

type Streamable[V any] interface {
	[]V | *V | Opt[V] | Stream[V] | func() Opt[V]
}

func NewStream[V any, S Streamable[V]](s S) Stream[V] {
	// note [bs]: for some of these, may be better to custom define an iterator
	// rather than
	switch narrowed := (interface{})(s).(type) {
	case []V:
		return Stream[V]{
			src: &listIter[V]{
				list: narrowed,
			},
		}
	case *V:
		if narrowed != nil {
			return Stream[V]{
				src: &listIter[V]{
					list: []V{*narrowed},
				},
			}
		} else {
			return Stream[V]{
				src: &listIter[V]{},
			}
		}
	case Opt[V]:
		if narrowed.IsPresent() {
			return Stream[V]{
				src: &listIter[V]{
					list: []V{*narrowed.val},
				},
			}
		} else {
			return Stream[V]{
				src: &listIter[V]{},
			}
		}
	case Stream[V]:
		return narrowed
	case func() Opt[V]:
		return Stream[V]{
			src: &iterFn[V]{
				fn: narrowed,
			},
		}
	default:
		panic("unreachable")
	}
}

// general note - the lack of additional method parameters does kind of mean
// that you're fairly limited w/ any sort of higher-order composition. Which,
// honestly, I think is for the best. I'd still like to dink around here, but
// the limitations mean you probably can't go too crazy.
//
// A stream is still sort of an interesting one since it can be more efficient
// for processing intermediaries - let's still play around a bit with that.

func StreamMap[V any, U any](s Stream[V], fn func(v V) U) Stream[U] {
	// fixme
	return Stream[U]{
		// values: MapList(s.values, fn),
	}
}

func StreamToPairs[V any, U comparable, T any](
	s Stream[V],
	fn func(v V) (U, T),
) Stream[Pair[U, T]] {
	// fixme
	return Stream[Pair[U, T]]{
		// values: MapList(s.values, func(v V) Pair[U, T] {
		// 	u, t := fn(v)
		// 	return Pair[U, T]{u, t}
		// }),
	}
}

// so - I'm not convinced keypair needs to exist as such. Seems like there would
// be cases were you could just restrict the pair to have a comparable key and
// that'd be it.
//
// Yeah, I'm thinking I agree. Let's kill it off for now.
//
// So - one other question about pairs. I want to support making it easy to map
// pair-to-pair. Which means not having to redo the pair declaration each time.
// I'd still need easy to-pair and from-pair options, but I think default pair
// behavior should be pair-to-pair.
//
// Another question worth mentioning: should you even have to return the key? It
// does come with problems - like if a user modifies the key. That said that's
// always a risk from collection - and it's worth emphasizing pair streaming is
// _not_ the same as key/value manipulation it's just optimized to cooperate
// with it.

func PStreamToMap[K comparable, V any](
	s Stream[Pair[K, V]],
) map[K]V {
	m := make(map[K]V)
	// todo [bs]: implement
	return m
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

type Iter[V any] interface {
	Next() Opt[V]
}

func IterEach[V any](iter Iter[V], fn func(v V)) {
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

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type Float interface {
	float32 | float64
}

type Complex interface {
	complex64 | complex128
}

type Number interface {
	Integer | Float
}

type AllNumber interface {
	Number | Complex
}

type iterInts[I Integer] struct {
	start, end I
	index      I
}

func (i *iterInts[I]) Next() Opt[I] {
	if i.index > i.end {
		return Opt[I]{}
	}
	v := i.start + i.index
	i.index += 1
	return Opt[I]{val: &v}
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

func order[N Number](v1, v2 N) (N, N) {

	// note [bs]: generally I don't think these sorts of number types really
	// belong in the same package, but for the sake of experimentation eh why not.
	if v1 < v2 {
		return v1, v2
	} else {
		return v2, v1
	}
}

func NewRange[I Integer](start, end I) Stream[I] {
	start, end = order(start, end)
	return Stream[I]{
		src: &iterInts[I]{
			start: start,
			end:   end,
		},
	}
}

func (l *listIter[V]) Next() Opt[V] {
	if l.nextIndex >= len(l.list) {
		return Opt[V]{}
	}
	v := l.list[l.nextIndex]
	l.nextIndex++
	return Opt[V]{&v}
}

type Readable interface {
	[]byte | string | *string | []rune | func(p []byte) (n int, err error)
}

type explicitReader struct {
	fn func(p []byte) (n int, err error)
}

func (sr *explicitReader) Read(p []byte) (n int, err error) {
	return sr.fn(p)
}

func AutoReader[R Readable](r R) io.Reader {
	switch narrowed := (interface{})(r).(type) {
	case []byte:
		return bytes.NewReader(narrowed)
	case string:
		return strings.NewReader(narrowed)
	case *string:
		if narrowed != nil {
			return strings.NewReader(*narrowed)
		} else {
			return strings.NewReader("")
		}
	case []rune:
		// ques [bs]: is this particularly inefficient? I kinda suspect so.
		return strings.NewReader(string(narrowed))
	case func(p []byte) (n int, err error):
		return &explicitReader{
			fn: narrowed,
		}
	default:
		panic("unreachable")
	}
}

// so interestingly, the mostly-original design doc suggests the following as a
// legal constraint. I'm guessing they decided along the way to restrain
// generics further. A little limiting but oh well.
//
// This may explain the issue -
// https://github.com/golang/go/issues/45346#issuecomment-862505803

func JSONReader[S Readable, V any](s S) func() V {
	return nil
}

func ReadJSON[Out any](s io.Reader) (out *Out, err error) {
	out = new(Out)
	err = json.NewDecoder(s).Decode(out)
	return out, err
}
