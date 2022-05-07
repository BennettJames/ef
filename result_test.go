package ef

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRes(t *testing.T) {

	t.Run("Val", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			assert.Equal(t, "hello", ResOfVal("hello").Val())
		})

		t.Run("OnErr", func(t *testing.T) {
			assert.Panics(t, func() {
				ResOfErr[string](fmt.Errorf("error")).Val()
			})
		})
	})

	t.Run("Err", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			assert.Equal(
				t,
				fmt.Errorf("error"),
				ResOfErr[string](fmt.Errorf("error")).Err())
		})

		t.Run("OnErr", func(t *testing.T) {
			assert.Panics(t, func() {
				ResOfVal("hello").Err()
			})
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(
				t,
				"<val='hello' err='<nil>'>",
				ResOfVal("hello").String())
		})

		t.Run("Err", func(t *testing.T) {
			assert.Equal(
				t,
				"<val='' err='hello'>",
				ResOfErr[string](fmt.Errorf("hello")).String())
		})
	})

	t.Run("Try", func(t *testing.T) {
		// ques [bs]: is there an idiomatic "void" type in go? If not, should I add
		// one?

		t.Run("success", func(t *testing.T) {
			r := ResOfVal(22)
			assert.Equal(
				t,
				ResOfVal("value: 22"),
				ResTry(r, func(v int) string {
					return fmt.Sprintf("value: %v", v)
				}))
		})

		t.Run("failErr", func(t *testing.T) {
			r := ResOfVal(22)
			panicVal := fmt.Errorf("error")
			assert.Equal(
				t,
				ResOfErr[string](panicVal),
				ResTry(r, func(v int) string {
					panic(panicVal)
				}))
		})

		t.Run("failOtherType", func(t *testing.T) {
			r := ResOfVal(22)
			var panicVal any = "error"
			assert.Equal(
				t,
				ResOfErr[string](&ResultRecoverError{panicVal}),
				ResTry(r, func(v int) string {
					panic(panicVal)
				}))
		})
	})

	t.Run("Flatten", func(t *testing.T) {
		t.Run("outerErr", func(t *testing.T) {
			val := fmt.Errorf("error")
			assert.Equal(t, ResOfErr[string](val), FlattenRes(ResOfVal(ResOfErr[string](val))))
		})

		t.Run("innerErr", func(t *testing.T) {
			val := fmt.Errorf("error")
			assert.Equal(t, ResOfErr[string](val), FlattenRes(ResOfVal(ResOfErr[string](val))))
		})

		t.Run("success", func(t *testing.T) {
			rVal := "value"
			assert.Equal(t, ResOfVal(rVal), FlattenRes(ResOfVal(ResOfVal(rVal))))
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
