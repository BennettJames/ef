package ef

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRes(t *testing.T) {

	t.Run("Val", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			assert.Equal(t, "hello", ResVal("hello").Val())
		})

		t.Run("OnErr", func(t *testing.T) {
			assert.Panics(t, func() {
				ResErr[string](fmt.Errorf("error")).Val()
			})
		})
	})

	t.Run("Err", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			assert.Equal(
				t,
				fmt.Errorf("error"),
				ResErr[string](fmt.Errorf("error")).Err())
		})

		t.Run("OnErr", func(t *testing.T) {
			assert.Panics(t, func() {
				ResVal("hello").Err()
			})
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			val, err := ResOf("hello", nil).Get()
			assert.Nil(t, err)
			assert.Equal(t, "hello", val)
		})

		t.Run("Err", func(t *testing.T) {
			val, err := ResOf("", fmt.Errorf("error")).Get()
			assert.Equal(t, fmt.Errorf("error"), err)
			assert.Equal(t, "", val)
		})
	})

	t.Run("GetPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			val, err := ResOf("hello", nil).GetPtr()
			assert.Nil(t, err)
			assert.Equal(t, Ptr("hello"), val)
		})

		t.Run("Err", func(t *testing.T) {
			val, err := ResOf("", fmt.Errorf("error")).GetPtr()
			assert.Nil(t, val)
			assert.Equal(t, fmt.Errorf("error"), err)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(
				t,
				"<val='hello' err='<nil>'>",
				ResVal("hello").String())
		})

		t.Run("Err", func(t *testing.T) {
			assert.Equal(
				t,
				"<val='' err='hello'>",
				ResErr[string](fmt.Errorf("hello")).String())
		})
	})

	t.Run("ResOf", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			r := ResOf(passthrough("value", nil))
			assert.Equal(t, ResVal("value"), r)
		})

		t.Run("Err", func(t *testing.T) {
			r := ResOf(passthrough("value", fmt.Errorf("error")))
			assert.Equal(t, ResErr[string](fmt.Errorf("error")), r)
		})
	})

	t.Run("ResOf2", func(t *testing.T) {
		t.Run("Vals", func(t *testing.T) {
			res := ResOf2("a", 22, nil)
			assert.Equal(t, ResVal(PairOf("a", 22)), res)
		})

		t.Run("Err", func(t *testing.T) {
			err := fmt.Errorf("error")
			res := ResOf2("a", 22, err)
			assert.Equal(t, ResErr[Pair[string, int]](err), res)

		})
	})

	t.Run("ResOfPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			val, err := ResOfPtr(Ptr("value"), nil).Get()
			assert.Equal(t, "value", val)
			assert.Nil(t, err)
		})

		t.Run("NilVal", func(t *testing.T) {
			val, err := ResOfPtr[string](nil, nil).Get()
			assert.Equal(t, "", val)
			assert.Equal(t, &UnexpectedNilError{}, err)
		})

		t.Run("Err", func(t *testing.T) {
			val, err := ResOfPtr(passthrough[*string](nil, fmt.Errorf("error"))).Get()
			assert.Equal(t, fmt.Errorf("error"), err)
			assert.Equal(t, "", val)
		})
	})

	t.Run("ResOfOpt", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			res := ResOfOpt(OptOf("hello"))
			assert.Equal(t, ResVal("hello"), res)
		})

		t.Run("Nil", func(t *testing.T) {
			res := ResOfOpt(OptOfPtr[string](nil))
			assert.Equal(t, ResErr[string](&UnexpectedNilError{}), res)
		})
	})

	t.Run("ResOfErr", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			val, err := ResErr[*string](fmt.Errorf("error")).Get()
			assert.Nil(t, val)
			assert.Equal(t, fmt.Errorf("error"), err)
		})

		t.Run("NilErr", func(t *testing.T) {
			val, err := ResErr[string](nil).Get()
			assert.Nil(t, err)
			assert.Equal(t, "", val)
		})
	})

	t.Run("IfVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			set := false
			ResVal("value").IfVal(func(val string) {
				set = true
			})
			assert.True(t, set)
		})

		t.Run("Err", func(t *testing.T) {
			set := false
			ResErr[string](fmt.Errorf("error")).IfVal(func(val string) {
				set = true
			})
			assert.False(t, set)
		})
	})

	t.Run("IfErr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			set := false
			ResVal("hello").IfErr(func(e error) {
				set = true
			})
			assert.False(t, set)
		})

		t.Run("Err", func(t *testing.T) {
			set := false
			ResErr[string](fmt.Errorf("error")).IfErr(func(e error) {
				set = true
			})
			assert.True(t, set)
		})
	})

	t.Run("ResDeref", func(t *testing.T) {
		t.Run("WithValue", func(t *testing.T) {
			assert.Equal(
				t,
				ResVal("hello"),
				ResDeref(ResVal(Ptr("hello"))))
		})

		t.Run("WithNil", func(t *testing.T) {
			assert.Equal(
				t,
				ResVal(""),
				ResDeref(ResVal[*string](nil)))

		})
	})

	t.Run("ResMap", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			in := 22
			ret := "hello"
			assert.Equal(
				t,
				ResVal(ret),
				ResMap(ResVal(in), func(val int) string {
					assert.Equal(t, in, val)
					return ret
				}))
		})

		t.Run("Err", func(t *testing.T) {
			err := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](err),
				ResMap(ResErr[int](err), func(int) string {
					panic("unreachable")
				}))
		})
	})

	t.Run("ResFlatMap", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			in := 22
			ret := "hello"
			assert.Equal(
				t,
				ResVal(ret),
				ResFlatMap(ResVal(in), func(val int) Res[string] {
					assert.Equal(t, in, val)
					return ResVal(ret)
				}))
		})

		t.Run("OuterErr", func(t *testing.T) {
			err := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](err),
				ResFlatMap(ResErr[int](err), func(int) Res[string] {
					panic("unreachable")
				}))
		})

		t.Run("InnerErr", func(t *testing.T) {
			in := 22
			err := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](err),
				ResFlatMap(ResErr[int](err), func(val int) Res[string] {
					assert.Equal(t, in, val)
					return ResErr[string](err)
				}))
		})
	})

	t.Run("ResRecover", func(t *testing.T) {
		t.Run("NoPanic", func(t *testing.T) {
			res := ResVal("value")
			ResRecover(&res)
			assert.Equal(t, "value", res.Val())
		})

		t.Run("NilRes", func(t *testing.T) {
			assert.Panics(t, func() {
				ResRecover[string](nil)
			})
		})
	})

	t.Run("ResTry", func(t *testing.T) {
		t.Run("NoPanicVal", func(t *testing.T) {
			r := ResVal(22)
			assert.Equal(
				t,
				ResVal("value: 22"),
				ResTryMap(r, func(v int) string {
					return fmt.Sprintf("value: %v", v)
				}))
		})

		t.Run("NoPanicErr", func(t *testing.T) {
			r := ResErr[string](fmt.Errorf("error"))
			assert.Equal(
				t,
				r,
				ResTryMap(r, func(v string) string {
					return "value"
				}))
		})

		t.Run("PanicErr", func(t *testing.T) {
			r := ResVal(22)
			panicVal := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](panicVal),
				ResTryMap(r, func(v int) string {
					panic(panicVal)
				}))
		})

		t.Run("PanicOther", func(t *testing.T) {
			r := ResVal(22)
			var panicVal any = "error"
			assert.Equal(
				t,
				ResErr[string](&RecoverError{panicVal}),
				ResTryMap(r, func(v int) string {
					panic(panicVal)
				}))
		})
	})

	t.Run("ResFlatTry", func(t *testing.T) {
		t.Run("NoPanicVal", func(t *testing.T) {
			r := ResVal(22)
			assert.Equal(
				t,
				ResVal("value: 22"),
				ResTryFlatMap(r, func(v int) Res[string] {
					return ResVal(fmt.Sprintf("value: %v", v))
				}))
		})

		t.Run("NoPanicErr", func(t *testing.T) {
			r := ResErr[string](fmt.Errorf("error"))
			assert.Equal(
				t,
				r,
				ResTryFlatMap(r, func(v string) Res[string] {
					return ResVal("value")
				}))
		})

		t.Run("PanicErr", func(t *testing.T) {
			r := ResVal(22)
			panicVal := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](panicVal),
				ResTryFlatMap(r, func(v int) Res[string] {
					panic(panicVal)
				}))
		})
	})

	t.Run("Flatten", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			rVal := "value"
			assert.Equal(t,
				ResVal(rVal),
				ResFlatten(ResVal(ResVal(rVal))))
		})

		t.Run("ErrOuter", func(t *testing.T) {
			val := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](val),
				ResFlatten(ResVal(ResErr[string](val))))
		})

		t.Run("ErrInner", func(t *testing.T) {
			val := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](val),
				ResFlatten(ResErr[Res[string]](val)))
		})
	})
}

func passthrough[V any](v V, e error) (V, error) {
	return v, e
}
