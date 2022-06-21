package ef

import "math"

// SingedInteger is an interface union of all signed integer primitives.
type SingedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger is an interface union of all unsigned integer primitives.
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is an interface union of all integer primitives.
type Integer interface {
	SingedInteger | UnsignedInteger
}

// Float is an interfaces union of all float primitives.
type Float interface {
	float32 | float64
}

// Complex is an interface union of all explicitly complex primitives.
type Complex interface {
	complex64 | complex128
}

// Number is an interface union of all real numbers.
//
// Authors note: this is called "Number" rather than "RealNumber" mostly as it
// is shorter, and order-able numbers are used much more often than complex.
type Number interface {
	Integer | Float
}

// AllNumber is an interface union of all complex and real number types.
type AllNumber interface {
	Number | Complex
}

// Order returns the two given numbers so the first is the lower value.
func Order[N Number](v1, v2 N) (low N, high N) {
	if v1 < v2 {
		return v1, v2
	} else {
		return v2, v1
	}
}

// Range returns a stream that consists of values from the start until the end.
// The end is exclusive - that is, the stream consists of integers less than the
// end.
func Range[I Integer](start, end I) Stream[I] {
	return Stream[I]{
		src: &rangeStruct[I]{
			start: start,
			end:   end,
		},
	}
}

// rangeStruct is a simple iterator that supports the range
type rangeStruct[I Integer] struct {
	// note [bs]: might iterate a bit more on this to see what pattern is most
	// efficient / flexible.
	start, end I
	offset     I
}

func (ri *rangeStruct[T]) forEach(fn func(T) bool) {
	// note [bs]: double check the bound behavior here
	for v := ri.start; v < ri.end; v++ {
		if !fn(v) {
			break
		}
	}
}

// RangeIncl creates a stream that goes from start to end, including the end
// value.
func RangeIncl[I Integer](start, end I) Stream[I] {
	return Stream[I]{
		src: &rangeInclStruct[I]{
			start: start,
			end:   end,
		},
	}
}

type rangeInclStruct[I Integer] struct {
	// note [bs]: there's some duplication here, but I sorta suspect (with light
	// evidence) that just enumerating the struct types for slightly different
	// behavior is more appropriate, if tedious, and appropriate for a library.
	// Still would like to quantify that a bit better.
	start, end I
	offset     I
}

func (ri *rangeInclStruct[T]) forEach(fn func(T) bool) {
	// note [bs]: double check the bound behavior here
	for v := ri.start; v <= ri.end; v++ {
		if !fn(v) {
			break
		}
	}
}

// RangeRev produces the same values as Range, but in reverse. Note it is still
// exclusive on the last value - so the first value is `end - 1``, and the last
// value is `start`.
func RangeRev[I Integer](start, end I) Stream[I] {
	return Stream[I]{
		src: &rangeRevStruct[I]{
			start: start,
			end:   end,
		},
	}
}

type rangeRevStruct[I Integer] struct {
	start, end I
	offset     I
}

func (ri *rangeRevStruct[T]) forEach(fn func(T) bool) {
	// note [bs]: double check the bound behavior here
	for v := ri.end - 1; v >= ri.start; v-- {
		if !fn(v) {
			break
		}
	}
}

// RangeRevIncl produces the same values as RangeIncl, but in reverse.
func RangeRevIncl[I Integer](start, end I) Stream[I] {
	return Stream[I]{
		src: &rangeReverseInclStruct[I]{
			start: start,
			end:   end,
		},
	}
}

type rangeReverseInclStruct[I Integer] struct {
	start, end I
	offset     I
}

func (ri *rangeReverseInclStruct[T]) forEach(fn func(T) bool) {
	// note [bs]: double check the bound behavior here
	for v := ri.end; v >= ri.start; v-- {
		if !fn(v) {
			break
		}
	}
}

// Min returns the lower of the two values.
func Min[N Number](v1, v2 N) N {
	if v1 <= v2 {
		return v1
	}
	return v2
}

// Max returns the higher of the two values.
func Max[N Number](v1, v2 N) N {
	if v1 >= v2 {
		return v1
	}
	return v2
}

// Add adds two numbers of the same type.
func Add[N AllNumber](v1, v2 N) N {
	// note [bs]: I added this as an experiment for function helpers for reduce.
	// not convinced by it - I think I should move to a system where
	//
	// Not 100% sure that'd be the right approach either, but it'd be interested
	// to try.
	return v1 + v2
}

// Mult multiplies two number of the same type together.
func Mult[N AllNumber](v1, v2 N) N {
	return v1 * v2
}

// MinNumber returns the minimum possible number for the given number type.
func MinNumber[N Number]() N {
	switch any(*new(N)).(type) {
	case int:
		x := int(math.MinInt)
		return N(x)
	case int8:
		x := math.MinInt8
		return N(x)
	case int16:
		x := (math.MinInt16)
		return N(x)
	case int32:
		x := math.MinInt32
		return N(x)
	case int64:
		x := math.MinInt64
		return N(x)
	case uint:
		return N(0)
	case uint8:
		return N(0)
	case uint16:
		return N(0)
	case uint32:
		return N(0)
	case uint64:
		return N(0)
	case uintptr:
		return N(0)
	case float32:
		x := float32(math.Inf(-1))
		return N(x)
	case float64:
		x := float64(math.Inf(-1))
		return N(x)
	default:
		panic("unreachable")
	}
}

// MaxNumber returns the maximum possible value for the given number type.
func MaxNumber[N Number]() N {
	switch any(*new(N)).(type) {
	case int:
		x := int(math.MaxInt)
		return N(x)
	case int8:
		x := int8(math.MaxInt8)
		return N(x)
	case int16:
		x := int16(math.MaxInt16)
		return N(x)
	case int32:
		x := int32(math.MaxInt32)
		return N(x)
	case int64:
		x := int64(math.MaxInt64)
		return N(x)
	case uint:
		x := uint(math.MaxUint)
		return N(x)
	case uint8:
		x := uint8(math.MaxUint8)
		return N(x)
	case uint16:
		x := uint16(math.MaxUint16)
		return N(x)
	case uint32:
		x := uint32(math.MaxUint32)
		return N(x)
	case uint64:
		x := uint64(math.MaxUint64)
		return N(x)
	case uintptr:
		x := uintptr(math.MaxUint64)
		return N(x)
	case float32:
		x := float32(math.Inf(1))
		return N(x)
	case float64:
		x := float64(math.Inf(1))
		return N(x)
	default:
		panic("unreachable")
	}
}
