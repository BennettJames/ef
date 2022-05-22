package rangebench2

type iterInts[I integer] struct {
	start, end I
	index      I
}

func (i *iterInts[I]) next() opt[I] {
	v := i.start + i.index
	if v >= i.end {
		return opt[I]{}
	}
	i.index += 1
	return optOf(v)
}

// rangeStruct creates a range via a struct to manage start / end values.
func rangeStruct[I integer](start, end I) stream[I] {
	return stream[I]{
		src: &iterInts[I]{
			start: start,
			end:   end,
		},
	}
}

// rangeCloser creates a range stream via a closure function that inlines state
// behavior.
func rangeClosure[I integer](start, end I) stream[I] {
	offset := I(0)
	return newFnStream(func() opt[I] {
		val := start + offset
		if val > end {
			return opt[I]{}
		}
		offset++
		return opt[I]{value: val, present: true}
	})
}
