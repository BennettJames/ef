package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {

	t.Run("Get", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(t,
				"hello",
				NewOptValue("hello").UnsafeGet())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.PanicsWithError(t, (&UnexpectedNilError{}).Error(), func() {
				Opt[string]{}.UnsafeGet()
			})
		})
	})

	t.Run("GetPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(t,
				Ptr("hello"),
				NewOptValue("hello").GetPtr())
		})

		t.Run("Empty", func(t *testing.T) {
			var strPtr *string
			assert.Equal(t,
				strPtr,
				Opt[string]{}.GetPtr())
		})
	})

	t.Run("IsVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.True(t, NewOptValue("hello").HasVal())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.False(t, Opt[string]{}.HasVal())
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.False(t, NewOptValue("hello").IsEmpty())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.True(t, Opt[string]{}.IsEmpty())
		})
	})

	t.Run("IfVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			run := false
			opt := NewOptValue("hello")
			ret := opt.IfVal(func(v string) {
				assert.Equal(t, "hello", v)
				run = true
			})
			assert.Equal(t, opt, ret)
			assert.True(t, run)
		})

		t.Run("Empty", func(t *testing.T) {
			opt := Opt[string]{}
			ret := opt.IfVal(func(v string) {
				panic(&UnreachableError{})
			})
			assert.Equal(t, opt, ret)
		})
	})

	t.Run("IfEmpty", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			opt := NewOptValue("hello")
			ret := opt.IfEmpty(func() {
				panic(&UnreachableError{})
			})
			assert.Equal(t, opt, ret)
		})

		t.Run("Empty", func(t *testing.T) {
			run := false
			opt := Opt[string]{}
			ret := opt.IfEmpty(func() {
				run = true
			})
			assert.Equal(t, opt, ret)
			assert.True(t, run)
		})
	})

	t.Run("GetPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(t,
				Ptr("hello"),
				NewOptValue("hello").GetPtr())
		})

		t.Run("Empty", func(t *testing.T) {
			var strPtr *string
			assert.Equal(t,
				strPtr,
				Opt[string]{}.GetPtr())
		})
	})

	t.Run("IsVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.True(t, NewOptValue("hello").HasVal())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.False(t, Opt[string]{}.HasVal())
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.False(t, NewOptValue("hello").IsEmpty())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.True(t, Opt[string]{}.IsEmpty())
		})
	})

	t.Run("IfVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			run := false
			NewOptValue("hello").IfVal(func(v string) {
				assert.Equal(t, "hello", v)
				run = true
			})
			assert.True(t, run)
		})

		t.Run("Empty", func(t *testing.T) {
			run := false
			Opt[string]{}.IfVal(func(v string) {
				run = true
			})
			assert.False(t, run)
		})
	})

	t.Run("Or", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				"hello",
				NewOptValue("hello").Or("world"))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				"world",
				Opt[string]{}.Or("world"))
		})
	})

	t.Run("OrCalc", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				"hello",
				NewOptValue("hello").OrCalc(func() string {
					panic("unreachable")
				}))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				"world",
				Opt[string]{}.OrCalc(func() string {
					return "world"
				}))
		})
	})

	t.Run("ToList", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				[]string{"hello"},
				NewOptValue("hello").ToList())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				[]string{},
				Opt[string]{}.ToList())
		})
	})
}
