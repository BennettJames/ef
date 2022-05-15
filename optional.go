package ef

type Opt[T any] struct {
	value   T
	present bool
}

func OptOf[T any](val T) Opt[T] {
	return Opt[T]{
		value:   val,
		present: true,
	}
}

func OptOfPtr[T any](val *T) Opt[T] {
	if val == nil {
		return Opt[T]{}
	}
	return OptOf(*val)
}

func (o Opt[T]) Val() T {
	if !o.present {
		// todo [bs]: let's use a standard null pointer exception type here.
		// Want that in a few different places.
		panic("'Get' called on empty optional")
	}
	return o.value
}

// GetPtr will
func (o Opt[T]) GetPtr() *T {
	if !o.present {
		return nil
	}
	return &o.value
}

func (o Opt[T]) IsVal() bool {
	return o.present
}

func (o Opt[T]) IfVal(fn func(v T)) {
	if o.present {
		fn(o.value)
	}
}

func (o Opt[T]) Get() (value *T, isSet bool) {
	// todo [bs]: let's rename this. "Try" has taken on a life of
	// it's own here.
	return &o.value, o.present
}

func (o Opt[T]) Or(altVal T) T {
	if o.IsVal() {
		return o.value
	}
	return altVal
}

func (o Opt[T]) OrCalc(orFn func() T) T {
	if o.IsVal() {
		return o.value
	}
	return orFn()
}

func (o Opt[T]) ToList() []T {
	if o.IsVal() {
		return []T{o.value}
	} else {
		return nil
	}
}

func OptMap[T any, U any](o Opt[T], fn func(v T) U) Opt[U] {
	if !o.present {
		return Opt[U]{}
	}
	return OptOf(fn(o.value))
}

func OptFlatMap[T any, U any](o Opt[T], fn func(v T) Opt[U]) Opt[U] {
	if !o.present {
		return Opt[U]{}
	}
	return fn(o.value)
}

// OptFlatten reduces a nested optional down to one.
func OptFlatten[T any](o Opt[Opt[T]]) Opt[T] {
	if o.IsVal() {
		return o.Val()
	}
	return Opt[T]{}
}
