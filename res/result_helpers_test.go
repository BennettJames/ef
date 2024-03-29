package res

import (
	"fmt"
	"testing"

	"github.com/BennettJames/ef"
	"github.com/stretchr/testify/assert"
)

func TestRes(t *testing.T) {

	t.Run("Of", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			r := Of(passthrough("value", nil))
			assert.Equal(t, Val("value"), r)
		})

		t.Run("Err", func(t *testing.T) {
			r := Of(passthrough("value", fmt.Errorf("error")))
			assert.Equal(t, Err[string](fmt.Errorf("error")), r)
		})
	})

	t.Run("Of2", func(t *testing.T) {
		t.Run("Vals", func(t *testing.T) {
			res := Of2("a", 22, nil)
			assert.Equal(t, Val(ef.PairOf("a", 22)), res)
		})

		t.Run("Err", func(t *testing.T) {
			err := fmt.Errorf("error")
			res := Of2("a", 22, err)
			assert.Equal(t, Err[ef.Pair[string, int]](err), res)

		})
	})

	t.Run("OfPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			val, err := OfPtr(ef.Ptr("value"), nil).Get()
			assert.Equal(t, "value", val)
			assert.Nil(t, err)
		})

		t.Run("NilVal", func(t *testing.T) {
			val, err := OfPtr[string](nil, nil).Get()
			assert.Equal(t, "", val)
			assert.Equal(t, &ef.UnexpectedNilError{}, err)
		})

		t.Run("Err", func(t *testing.T) {
			val, err := OfPtr(passthrough[*string](nil, fmt.Errorf("error"))).Get()
			assert.Equal(t, fmt.Errorf("error"), err)
			assert.Equal(t, "", val)
		})
	})

	t.Run("OfOpt", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			res := OfOpt(ef.NewOptValue("hello"))
			assert.Equal(t, Val("hello"), res)
		})

		t.Run("Nil", func(t *testing.T) {
			res := OfOpt(ef.Opt[string]{})
			assert.Equal(t, Err[string](&ef.UnexpectedNilError{}), res)
		})
	})

	t.Run("OfErr", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			val, err := Err[*string](fmt.Errorf("error")).Get()
			assert.Nil(t, val)
			assert.Equal(t, fmt.Errorf("error"), err)
		})

		t.Run("NilErr", func(t *testing.T) {
			val, err := Err[string](nil).Get()
			assert.Nil(t, err)
			assert.Equal(t, "", val)
		})
	})

	t.Run("Deref", func(t *testing.T) {
		t.Run("WithValue", func(t *testing.T) {
			assert.Equal(
				t,
				Val("hello"),
				Deref(Val(ef.Ptr("hello"))))
		})

		t.Run("WithNil", func(t *testing.T) {
			assert.Equal(
				t,
				Val(""),
				Deref(Val[*string](nil)))

		})
	})

	t.Run("Map", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			in := 22
			ret := "hello"
			assert.Equal(
				t,
				Val(ret),
				Map(Val(in), func(val int) string {
					assert.Equal(t, in, val)
					return ret
				}))
		})

		t.Run("Err", func(t *testing.T) {
			err := fmt.Errorf("error")
			assert.Equal(
				t,
				Err[string](err),
				Map(Err[int](err), func(int) string {
					panic("unreachable")
				}))
		})
	})

	t.Run("FlatMap", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			in := 22
			ret := "hello"
			assert.Equal(
				t,
				Val(ret),
				FlatMap(Val(in), func(val int) ef.Res[string] {
					assert.Equal(t, in, val)
					return Val(ret)
				}))
		})

		t.Run("OuterErr", func(t *testing.T) {
			err := fmt.Errorf("error")
			assert.Equal(
				t,
				Err[string](err),
				FlatMap(Err[int](err), func(int) ef.Res[string] {
					panic("unreachable")
				}))
		})

		t.Run("InnerErr", func(t *testing.T) {
			in := 22
			err := fmt.Errorf("error")
			assert.Equal(
				t,
				Err[string](err),
				FlatMap(Err[int](err), func(val int) ef.Res[string] {
					assert.Equal(t, in, val)
					return Err[string](err)
				}))
		})
	})

	t.Run("Recover", func(t *testing.T) {
		t.Run("NoPanic", func(t *testing.T) {
			res := Val("value")
			Recover(&res)
			assert.Equal(t, "value", res.Val())
		})

		t.Run("NilRes", func(t *testing.T) {
			assert.Panics(t, func() {
				Recover[string](nil)
			})
		})
	})

	t.Run("Try", func(t *testing.T) {
		t.Run("NoPanicVal", func(t *testing.T) {
			r := Val(22)
			assert.Equal(
				t,
				Val("value: 22"),
				TryMap(r, func(v int) string {
					return fmt.Sprintf("value: %v", v)
				}))
		})

		t.Run("NoPanicErr", func(t *testing.T) {
			r := Err[string](fmt.Errorf("error"))
			assert.Equal(
				t,
				r,
				TryMap(r, func(v string) string {
					return "value"
				}))
		})

		t.Run("PanicErr", func(t *testing.T) {
			r := Val(22)
			panicVal := fmt.Errorf("error")
			assert.Equal(
				t,
				Err[string](panicVal),
				TryMap(r, func(v int) string {
					panic(panicVal)
				}))
		})

		t.Run("PanicOther", func(t *testing.T) {
			r := Val(22)
			var panicVal any = "error"
			assert.Equal(
				t,
				Err[string](ef.NewRecoverError(panicVal)),
				TryMap(r, func(v int) string {
					panic(panicVal)
				}))
		})
	})

	t.Run("FlatTry", func(t *testing.T) {
		t.Run("NoPanicVal", func(t *testing.T) {
			r := Val(22)
			assert.Equal(
				t,
				Val("value: 22"),
				TryFlatMap(r, func(v int) ef.Res[string] {
					return Val(fmt.Sprintf("value: %v", v))
				}))
		})

		t.Run("NoPanicErr", func(t *testing.T) {
			r := Err[string](fmt.Errorf("error"))
			assert.Equal(
				t,
				r,
				TryFlatMap(r, func(v string) ef.Res[string] {
					return Val("value")
				}))
		})

		t.Run("PanicErr", func(t *testing.T) {
			r := Val(22)
			panicVal := fmt.Errorf("error")
			assert.Equal(
				t,
				Err[string](panicVal),
				TryFlatMap(r, func(v int) ef.Res[string] {
					panic(panicVal)
				}))
		})
	})

	t.Run("Flatten", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			rVal := "value"
			assert.Equal(t,
				Val(rVal),
				Flatten(Val(Val(rVal))))
		})

		t.Run("ErrOuter", func(t *testing.T) {
			val := fmt.Errorf("error")
			assert.Equal(
				t,
				Err[string](val),
				Flatten(Val(Err[string](val))))
		})

		t.Run("ErrInner", func(t *testing.T) {
			val := fmt.Errorf("error")
			assert.Equal(
				t,
				Err[string](val),
				Flatten(Err[ef.Res[string]](val)))
		})
	})
}

func passthrough[V any](v V, e error) (V, error) {
	return v, e
}
