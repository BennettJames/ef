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

	t.Run("ResOfPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			val, err := ResOfPtr(Ptr("value"), nil).Get()
			assert.Equal(t, "value", val)
			assert.Nil(t, err)
		})

		t.Run("NilVal", func(t *testing.T) {
			val, err := ResOfPtr[string](nil, nil).Get()
			assert.Equal(t, "", val)
			assert.Equal(t, &ResultNilError{}, err)
		})

		t.Run("Err", func(t *testing.T) {
			val, err := ResOfPtr(passthrough[*string](nil, fmt.Errorf("error"))).Get()
			assert.Equal(t, fmt.Errorf("error"), err)
			assert.Equal(t, "", val)
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

	t.Run("ResTry", func(t *testing.T) {
		// ques [bs]: is there an idiomatic "void" type in go? If not, should I add
		// one?

		t.Run("NoPanicVal", func(t *testing.T) {
			r := ResVal(22)
			assert.Equal(
				t,
				ResVal("value: 22"),
				ResTry(r, func(v int) string {
					return fmt.Sprintf("value: %v", v)
				}))
		})

		t.Run("NoPanicErr", func(t *testing.T) {
			r := ResErr[string](fmt.Errorf("error"))
			assert.Equal(
				t,
				r,
				ResTry(r, func(v string) string {
					return "value"
				}))
		})

		t.Run("PanicErr", func(t *testing.T) {
			r := ResVal(22)
			panicVal := fmt.Errorf("error")
			assert.Equal(
				t,
				ResErr[string](panicVal),
				ResTry(r, func(v int) string {
					panic(panicVal)
				}))
		})

		t.Run("PanicOther", func(t *testing.T) {
			r := ResVal(22)
			var panicVal any = "error"
			assert.Equal(
				t,
				ResErr[string](&ResultRecoverError{panicVal}),
				ResTry(r, func(v int) string {
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

func genericReturn() (string, error) {
	return "", fmt.Errorf("(error)")
}

func Create[T any]() T {
	return *new(T)
}

type Hello interface {
	SayHi() string
}

type helloImpl struct {
}

func (hi *helloImpl) SayHi() string {
	return "hi"
}

func passthrough[V any](v V, e error) (V, error) {
	return v, e
}

func TestPtrAssignment(t *testing.T) {

	t.Run("intTest", func(t *testing.T) {
		setInt := func(addr *int, value int) {
			*addr = value
		}

		value := 0
		fmt.Println("value is:", value)

		setInt(&value, 10)
		fmt.Println("value is:", value)
	})

	t.Run("resTest", func(t *testing.T) {
		setRes := func(addr *Res[string], value Res[string]) {
			*addr = value
		}

		value := ResVal("hello")
		fmt.Println("value is:", value)

		setRes(&value, ResVal("there"))
		fmt.Println("value is:", value)

	})
}

func TestScratchRes(t *testing.T) {

	t.Run("panic experiment", func(t *testing.T) {
		var res Res[string]
		defer ResRecover(&res)

		panic("test")
	})

	t.Run("deref experiment", func(t *testing.T) {
		nilRef := ResOf[*string](nil, nil)
		defaultRef := ResMap(nilRef, DerefFn[string]())
		fmt.Println("@@@ default ref - ", defaultRef)
	})
}

func DerefFn[V any]() func(*V) V {
	// so - I'm not yet sure if this is worth keeping, but it does reveal an
	// interesting constraint. A generic function cannot be be used as a "bare"
	// value; it must be wrapped. Making a function like this and forcing the caller
	// to specify the generic parameter does work, however.
	return func(val *V) V {
		if val == nil {
			return *new(V)
		}
		return *val
	}
}
