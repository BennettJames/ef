package ef

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
