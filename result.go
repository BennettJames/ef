package ef

import "fmt"

// Res is a "result" that tracks a value that can be either an error
// or a value, but not both. It comes with a set of methods and function
// helpers to make it easier to correctly work with errors.
//
// A note on nil - if a result wraps a nilable type (e.g. a pointer),
// then it is possible that the value and the error can be nil. This is
// still considered a "value typed result".
type Res[T any] struct {
	val T
	err error
}

func NewResValue[T any](val T) Res[T] {
	return Res[T]{val: val}
}

func NewResError[T any](err error) Res[T] {
	return Res[T]{err: err}
}

// Get returns both the value and the error of the result.
func (r Res[T]) Get() (T, error) {
	return r.val, r.err
}

// GetPtr returns the value / error pair, but with the value wrapped as a
// pointer.
func (r Res[T]) GetPtr() (*T, error) {
	if r.err != nil {
		return nil, r.err
	}
	// ques [bs]: does this expose any weird undesirable mutability?
	return &r.val, nil
}

// Val returns the underlying value from the result, or panics if the result is
// not a value type.
func (r Res[T]) Val() T {
	if !r.IsVal() {
		panic("res.Val() called on non-value result")
	}
	return r.val
}

// Err returns the underlying error from the result, or panics if the result
// is not an error type.
func (r Res[T]) Err() error {
	if r.IsVal() {
		panic("res.Err() called on non-error result")
	}
	return r.err
}

// IsVal indicates if the result has a value (and is not an error).
func (r Res[T]) IsVal() bool {
	return r.err == nil
}

// IsErr indicates if the result is an error (and does not have a value).
func (r Res[T]) IsErr() bool {
	return r.err != nil
}

// IfVal will execute the passed function if the result is a value.
func (r Res[T]) IfVal(fn func(val T)) {
	if r.IsVal() {
		fn(r.val)
	}
}

// IfErr will execute the passed function if the result is an error.
func (r Res[T]) IfErr(fn func(e error)) {
	if !r.IsVal() {
		fn(r.err)
	}
}

// String is just a simple string representation of the result for debugging.
func (r Res[T]) String() string {
	// ques [bs]: should this have more of a structural difference between
	// value / error?
	return fmt.Sprintf("<val='%v' err='%v'>", r.val, r.err)
}
