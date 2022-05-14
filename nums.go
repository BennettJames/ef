package ef

import "math"

type SingedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	SingedInteger | UnsignedInteger
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

func order[N Number](v1, v2 N) (N, N) {

	// note [bs]: generally I don't think these sorts of number types really
	// belong in the same package, but for the sake of experimentation eh why not.
	if v1 < v2 {
		return v1, v2
	} else {
		return v2, v1
	}
}

func Range[I Integer](start, end I) Stream[I] {
	return Stream[I]{
		src: &iterInts[I]{
			start: start,
			end:   end,
		},
	}
}

func (i *iterInts[I]) Next() Opt[I] {
	v := i.start + i.index
	if v > i.end {
		return Opt[I]{}
	}
	i.index += 1
	return OptOf(v)
}

func Min[N Number](v1, v2 N) N {
	if v1 <= v2 {
		return v1
	}
	return v2
}

func Max[N Number](v1, v2 N) N {
	if v1 >= v2 {
		return v1
	}
	return v2
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
