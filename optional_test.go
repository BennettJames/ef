package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {

	t.Run("OptOf", func(t *testing.T) {
		assert.Equal(t,
			Opt[string]{
				value:   "hello",
				present: true,
			},
			OptOf("hello"))
	})

	t.Run("OptEmpty", func(t *testing.T) {
		assert.Equal(t,
			Opt[string]{},
			OptEmpty[string]())
	})

	t.Run("OptOfPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			strVal := "hello"
			var strPtr *string = &strVal
			assert.Equal(t,
				Opt[string]{
					value:   strVal,
					present: true,
				},
				OptOfPtr(strPtr))
		})

		t.Run("Empty", func(t *testing.T) {
			var strPtr *string
			assert.Equal(t,
				Opt[string]{
					value:   "",
					present: false,
				},
				OptOfPtr(strPtr))
		})
	})

	t.Run("OptOfOk", func(t *testing.T) {
		t.Run("Ok", func(t *testing.T) {
			assert.Equal(t,
				Opt[string]{
					value:   "hello",
					present: true,
				},
				OptOfOk("hello", true))
		})

		t.Run("NotOk", func(t *testing.T) {
			assert.Equal(t,
				Opt[string]{
					value:   "",
					present: false,
				},
				OptOfOk("hello", false))
		})
	})

	t.Run("OptSliceGet", func(t *testing.T) {
		s := []string{"hello", "world"}
		t.Run("InBounds", func(t *testing.T) {
			assert.Equal(t,
				OptOf("hello"),
				OptSliceGet(s, 0))
			assert.Equal(t,
				OptOf("world"),
				OptSliceGet(s, 1))
		})

		t.Run("OutOfBounds", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[string](),
				OptSliceGet(s, -1))
			assert.Equal(t,
				OptEmpty[string](),
				OptSliceGet(s, 2))
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(t,
				"hello",
				OptOf("hello").Get())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.PanicsWithError(t, (&UnexpectedNilError{}).Error(), func() {
				OptEmpty[string]().Get()
			})
		})
	})

	t.Run("GetPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.Equal(t,
				Ptr("hello"),
				OptOf("hello").GetPtr())
		})

		t.Run("Empty", func(t *testing.T) {
			var strPtr *string
			assert.Equal(t,
				strPtr,
				OptEmpty[string]().GetPtr())
		})
	})

	t.Run("OptMapGet", func(t *testing.T) {
		m := map[string]string{
			"hello": "world",
		}
		t.Run("Ok", func(t *testing.T) {
			assert.Equal(t,
				Opt[string]{
					value:   "world",
					present: true,
				},
				OptMapGet(m, "hello"))
		})

		t.Run("NotOk", func(t *testing.T) {
			assert.Equal(t,
				Opt[string]{
					value:   "",
					present: false,
				},
				OptMapGet(m, "badkey"))
		})
	})

	t.Run("IsVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.True(t, OptOf("hello").IsVal())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.False(t, OptEmpty[string]().IsVal())
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			assert.False(t, OptOf("hello").IsEmpty())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.True(t, OptEmpty[string]().IsEmpty())
		})
	})

	t.Run("IfVal", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			run := false
			OptOf("hello").IfVal(func(v string) {
				assert.Equal(t, "hello", v)
				run = true
			})
			assert.True(t, run)
		})

		t.Run("Empty", func(t *testing.T) {
			run := false
			OptEmpty[string]().IfVal(func(v string) {
				run = true
			})
			assert.False(t, run)
		})
	})

	t.Run("OptFlatten", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			// hrm, so the use of "optlike" here does add some inconvenience -
			// possibly enough for me to give up on it in this context.
			//
			// I have to specify the type here if I use it. Bleh. That means
			// both cases require a bit of extra typing. That is biasing me in
			// favor of just making this two different methods.
			innerO := OptOf("hello")
			outerO := OptOf(innerO)
			assert.Equal(t,
				OptOf("hello"),
				OptFlatten[string](outerO),
			)
		})

		t.Run("OuterEmpty", func(t *testing.T) {
			outerO := OptEmpty[Opt[string]]()
			assert.Equal(t,
				OptEmpty[string](),
				OptFlatten[string](outerO),
			)
		})

		t.Run("InnerEmpty", func(t *testing.T) {
			innerO := OptEmpty[string]()
			outerO := OptOf(innerO)
			assert.Equal(t,
				OptEmpty[string](),
				OptFlatten[string](outerO),
			)
		})

		t.Run("InnerPtr", func(t *testing.T) {
			var innerVal *string = Ptr("hello")
			outerO := OptOf(innerVal)
			assert.Equal(t,
				OptOf("hello"),
				OptFlatten[string](outerO),
			)
		})

		t.Run("InnerNilPtr", func(t *testing.T) {
			var innerVal *string
			outerO := OptOf(innerVal)
			assert.Equal(t,
				OptEmpty[string](),
				OptFlatten[string](outerO),
			)
		})
	})
}
