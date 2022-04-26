package ef

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

// ques [bs]: should the simpler functions here just be methods, provided they
// don't require external types? Hmm. Possibly; need to think on that one a bit
// more.

func OptOr[V any](o Opt[V], or V) V {
	if o.val == nil {
		return *o.val
	}
	return or
}

func OptOrCalc[V any](o Opt[V], orFn func() V) V {
	if o.val == nil {
		return *o.val
	}
	return orFn()
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
