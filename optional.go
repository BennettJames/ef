package ef

type Opt[V any] struct {
	val     V
	present bool
}

func NewOpt[V any](val V) Opt[V] {
	return Opt[V]{
		val:     val,
		present: true,
	}
}

func NewNullableOpt[V any](val *V) Opt[V] {
	if val == nil {
		return Opt[V]{
			present: false,
		}
	}
	return NewOpt(*val)
}

func (o Opt[V]) Get() V {
	// todo: consider renaming this some like "DangerouslyGet".
	if !o.present {
		panic("'Get' called on empty optional")
	}
	return o.val
}

func (o Opt[V]) IsPresent() bool {
	return o.present
}

func (o Opt[V]) IfSet(fn func(v V)) {
	if o.present {
		fn(o.val)
	}
}

func (o Opt[V]) Try() (value *V, isSet bool) {
	return &o.val, o.present
}

func (o Opt[V]) Or(altVal V) V {
	if o.IsPresent() {
		return o.val
	}
	return altVal
}

func (o Opt[V]) OrCalc(orFn func() V) V {
	if o.IsPresent() {
		return o.val
	}
	return orFn()
}

func (o Opt[V]) ToList() []V {
	if o.IsPresent() {
		return []V{o.val}
	} else {
		return nil
	}
}

// so - I believe the following could *not* be turned into a method,
// as methods have fairly restrictive generic options - basically
// can just operate on a root.

func OptMap[V any, T any](o Opt[V], fn func(v V) T) Opt[T] {
	if !o.present {
		return Opt[T]{}
	}
	return NewOpt(fn(o.val))
}

// OptFlatten reduces a nested optional down to one.
func OptFlatten[V any](o Opt[Opt[V]]) Opt[V] {
	if o.IsPresent() {
		return o.Get()
	}
	return Opt[V]{}
}
