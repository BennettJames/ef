package ef

type (
	// Opt represents an optional of the given type. An optional can either have a value,
	// or be empty.
	//
	// Used well, this works as a kind of typesafe pointer - a series of helper functions
	// make it easier to compose around and use in a safe way, and discourage any behavior
	// that can result in a nil exception.
	Opt[T any] struct {
		value   T
		present bool
	}
)

func NewOptValue[T any](val T) Opt[T] {
	return Opt[T]{
		value:   val,
		present: true,
	}
}

// UnsafeGet returns the value if present, and panics if it does not exist. Note
// this is a dangerous method to use - generally it's best to an alternative to
// safely process, like `IfVal`, `Or`, or `OptMap`. Aim to structure the usage
// of the optional so the code can't err and makes no assumptions.
func (o Opt[T]) UnsafeGet() T {
	if !o.present {
		// todo [bs]: I still sorta suspect that the nil error should be able
		// to contain and communicate some amount of context.
		panic(&UnexpectedNilError{})
	}
	return o.value
}

// GetPtr returns a pointer of the inner value (and is nil if the optional
// is empty).
func (o Opt[T]) GetPtr() *T {
	if !o.present {
		return nil
	}
	return &o.value
}

// HasVal indicates if the optional has a value.
func (o Opt[T]) HasVal() bool {
	return o.present
}

// IsEmpty indicates if the optional lacks a value.
func (o Opt[T]) IsEmpty() bool {
	return !o.present
}

// todo [bs]: add a matching IsEmpty pure function once this is moved to
// a subpackage for filtering.

// IfVal executes the provided function with the stored value if the optional
// has a value; otherwise does nothing. Returns itself for chaining.
func (o Opt[T]) IfVal(fn func(v T)) Opt[T] {
	if o.present {
		fn(o.value)
	}
	return o
}

// IfEmpty calls the passed function if the optional is empty, otherwise does
// nothing. Returns itself for chaining.
func (o Opt[T]) IfEmpty(fn func()) Opt[T] {
	if !o.present {
		fn()
	}
	return o
}

// Or returns the provided value if the optional is empty, or the value if it
// has one.
func (o Opt[T]) Or(altVal T) T {
	if o.HasVal() {
		return o.value
	}
	return altVal
}

// OrCalc calls and returns the value from the function if the optional is
// empty, the the value if the optional has one.
//
// This is an alternative to `Or` for cases where it may be undesirable to
// unnecessarily compute the alternative value - for instance, if the
// calculation is expensive.
func (o Opt[T]) OrCalc(fn func() T) T {
	if o.HasVal() {
		return o.value
	}
	return fn()
}

// ToList converts the optional to a list. If the optional is empty, then the
// list is; otherwise it consists of just the single value held by the optional.
func (o Opt[T]) ToList() []T {
	if o.HasVal() {
		return []T{o.value}
	} else {
		return []T{}
	}
}
