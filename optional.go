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

func (o Opt[T]) Get() T {
	// todo: consider renaming this some like "DangerouslyGet".
	if !o.present {
		panic("'Get' called on empty optional")
	}
	return o.value
}

func (o Opt[T]) IsPresent() bool {
	return o.present
}

func (o Opt[T]) IfSet(fn func(v T)) {
	if o.present {
		fn(o.value)
	}
}

func (o Opt[T]) Try() (value *T, isSet bool) {
	return &o.value, o.present
}

func (o Opt[T]) Or(altVal T) T {
	if o.IsPresent() {
		return o.value
	}
	return altVal
}

func (o Opt[T]) OrCalc(orFn func() T) T {
	if o.IsPresent() {
		return o.value
	}
	return orFn()
}

func (o Opt[T]) ToList() []T {
	if o.IsPresent() {
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

// OptFlatten reduces a nested optional down to one.
func OptFlatten[T any](o Opt[Opt[T]]) Opt[T] {
	if o.IsPresent() {
		return o.Get()
	}
	return Opt[T]{}
}
