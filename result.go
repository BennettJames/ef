package ef

import "fmt"

// Res is a "result" that tracks a value that can be either an error
// or a value, but not both. It comes with a set of methods and function
// helpers to make it easier to correctly work with errors.
//
// A note on nil - if a result wraps a nilable type (e.g. a pointer),
// then it is possible that the value and the error can be nil. This is
// still considered a "value typed result".
type Res[V any] struct {
	val V
	err error
}

func (r Res[V]) Get() (V, error) {
	return r.val, r.err
}

// Val returns the underlying value from the result, or panics if the result is
// not a value type.
func (r Res[V]) Val() V {
	if !r.IsVal() {
		panic("res.Val() called on non-value result")
	}
	return r.val
}

// Err returns the underlying error from the result, or panics if the result
// is not an error type.
func (r Res[V]) Err() error {
	if r.IsVal() {
		panic("res.Err() called on non-error result")
	}
	return r.err
}

// String is just a simple string representation of the result for debugging.
func (r Res[V]) String() string {
	// ques [bs]: should this have more of a structural difference between
	// value / error?
	return fmt.Sprintf("<val='%v' err='%v'>", r.val, r.err)
}

// ResOf creates a result from a pair of values
func ResOf[V any](v V, e error) Res[V] {
	if e == nil {
		return ResOfVal(v)
	}
	return ResOfErr[V](e)
}

func ResOfPtr[V any](v *V, e error) Res[V] {
	if e != nil {
		return ResOfErr[V](e)
	}
	if v != nil {
		return ResOfVal(*v)
	}

	// so - what do I do if error is nil, but value also is?
	//
	// I see two main options:
	//
	// - create a nil pointer error, return it
	//
	// - return a zero-value for the v and use it
	//
	// I think I prefer the latter? The former

	panic("undecided on what to do here")
}

func ResOfVal[V any](val V) Res[V] {
	return Res[V]{
		val: val,
	}
}

func ResOfErr[V any](e error) Res[V] {
	// ques [bs]: what do I want to do here if the error is nil?
	// I should either panic or create a default value;
	return Res[V]{
		err: e,
	}
}

// IsVal indicates if the result has a value (and is not an error).
func (r Res[V]) IsVal() bool {
	return r.err == nil
}

// IsErr indicates if the result is an error (and does not have a value).
func (r Res[V]) IsErr() bool {
	return r.err != nil
}

// IfVal will execute the passed function if the result is a value.
func (r Res[V]) IfVal(fn func(v V)) {
	if r.IsVal() {
		fn(r.val)
	}
}

// IfErr will execute the passed function if the result is an error.
func (r Res[V]) IfErr(fn func(e error)) {
	if !r.IsVal() {
		fn(r.err)
	}
}

// ResMap will execute the passed function if the result is a
func ResMap[T, U any](r Res[T], fn func(val T) U) Res[U] {
	if r.IsVal() {
		return ResOfVal(fn(r.val))
	}
	return ResOfErr[U](r.err)
}

// so - I'm not sure this is "good design", but I'd like to look into the
// possibility of "auto recovery" from nil / err results. Note that I will
// need to create some new error types here.

func RecoverRes[T any]() {
	// so - I don't think you can necessarily do this via
}

func ResTry[V, U any](r Res[V], fn func(v V) U) (res Res[U]) {
	defer func() {
		if r := recover(); r != nil {
			if err, isErr := r.(error); isErr {
				res = ResOfErr[U](err)
			} else {
				res = ResOfErr[U](&ResultRecoverError{recovered: r})
			}
		}
	}()

	if !r.IsVal() {
		return ResOfErr[U](r.err)
	}
	return ResOfVal(fn(r.val))
}

func FlattenRes[T any](r Res[Res[T]]) Res[T] {
	if r.IsErr() {
		return ResOfErr[T](r.err)
	}
	return r.val
}

type ResultRecoverError struct {
	recovered any
}

func (e *ResultRecoverError) Error() string {
	// note [bs]: not super happy with this text value; let's workshop it.
	return fmt.Sprintf("Recovered try with value: '%v'", e.recovered)
}
