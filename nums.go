package ef

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
	return NewOpt(v)
}
