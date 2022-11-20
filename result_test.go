package ef

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRes(t *testing.T) {

	t.Run("Val", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			assert.Equal(t, "hello", NewResValue("hello").Val())
		})

		t.Run("OnErr", func(t *testing.T) {
			assert.Panics(t, func() {
				NewResError[string](fmt.Errorf("error")).Val()
			})
		})
	})

	t.Run("Err", func(t *testing.T) {
		t.Run("Basic", func(t *testing.T) {
			assert.Equal(
				t,
				fmt.Errorf("error"),
				NewResError[string](fmt.Errorf("error")).Err())
		})

		t.Run("OnErr", func(t *testing.T) {
			assert.Panics(t, func() {
				NewResValue("hello").Err()
			})
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			val, err := NewResValue("hello").Get()
			assert.Nil(t, err)
			assert.Equal(t, "hello", val)
		})

		t.Run("Err", func(t *testing.T) {
			val, err := NewResError[string](fmt.Errorf("error")).Get()
			assert.Equal(t, fmt.Errorf("error"), err)
			assert.Equal(t, "", val)
		})
	})

	t.Run("GetPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			val, err := NewResValue("hello").GetPtr()
			assert.Nil(t, err)
			assert.Equal(t, Ptr("hello"), val)
		})

		t.Run("Err", func(t *testing.T) {
			val, err := NewResError[string](fmt.Errorf("error")).GetPtr()
			assert.Nil(t, val)
			assert.Equal(t, fmt.Errorf("error"), err)
		})
	})

	t.Run("String", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(
				t,
				"<val='hello' err='<nil>'>",
				NewResValue("hello").String())
		})

		t.Run("Err", func(t *testing.T) {
			assert.Equal(
				t,
				"<val='' err='hello'>",
				NewResError[string](fmt.Errorf("hello")).String())
		})
	})

	t.Run("IsErr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.False(t, NewResValue("val").IsErr())
		})

		t.Run("Err", func(t *testing.T) {
			assert.True(t, NewResError[string](fmt.Errorf("")).IsErr())
		})
	})

	t.Run("IfVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			set := false
			NewResValue("value").IfVal(func(val string) {
				set = true
			})
			assert.True(t, set)
		})

		t.Run("Err", func(t *testing.T) {
			set := false
			NewResError[string](fmt.Errorf("error")).IfVal(func(val string) {
				set = true
			})
			assert.False(t, set)
		})
	})

	t.Run("IfErr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			set := false
			NewResValue("hello").IfErr(func(e error) {
				set = true
			})
			assert.False(t, set)
		})

		t.Run("Err", func(t *testing.T) {
			set := false
			NewResError[string](fmt.Errorf("error")).IfErr(func(e error) {
				set = true
			})
			assert.True(t, set)
		})
	})
}

func passthrough[V any](v V, e error) (V, error) {
	return v, e
}
