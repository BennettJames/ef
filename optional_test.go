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
				OptOf("hello").UnsafeGet())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.PanicsWithError(t, (&UnexpectedNilError{}).Error(), func() {
				OptEmpty[string]().UnsafeGet()
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
			assert.True(t, OptOf("hello").HasVal())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.False(t, OptEmpty[string]().HasVal())
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
			assert.True(t, OptOf("hello").HasVal())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.False(t, OptEmpty[string]().HasVal())
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

	t.Run("Or", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				"hello",
				OptOf("hello").Or("world"))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				"world",
				OptEmpty[string]().Or("world"))
		})
	})

	t.Run("OrCalc", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				"hello",
				OptOf("hello").OrCalc(func() string {
					panic("unreachable")
				}))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				"world",
				OptEmpty[string]().OrCalc(func() string {
					return "world"
				}))
		})
	})

	t.Run("ToList", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				[]string{"hello"},
				OptOf("hello").ToList())
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				[]string{},
				OptEmpty[string]().ToList())
		})
	})

	t.Run("OptMap", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				OptOf(22),
				OptMap(OptOf("hello"), func(v string) int {
					assert.Equal(t, "hello", v)
					return 22
				}))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[int](),
				OptMap(OptEmpty[string](), func(v string) int {
					panic("unreachable")
				}))
		})
	})

	t.Run("OptFlatMap", func(t *testing.T) {
		t.Run("ValueToValue", func(t *testing.T) {
			assert.Equal(t,
				OptOf(22),
				OptFlatMap(OptOf("hello"), func(v string) Opt[int] {
					assert.Equal(t, "hello", v)
					return OptOf(22)
				}))
		})

		t.Run("ValueToEmpty", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[int](),
				OptFlatMap(OptOf("hello"), func(v string) Opt[int] {
					assert.Equal(t, "hello", v)
					return OptEmpty[int]()
				}))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[int](),
				OptFlatMap(OptEmpty[string](), func(v string) Opt[int] {
					panic("unreachable")
				}))
		})
	})

	t.Run("OptFlatten", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				OptOf("hello"),
				OptFlatten(OptOf(OptOf("hello"))))
		})

		t.Run("InnerEmpty", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[string](),
				OptFlatten(OptOf(OptEmpty[string]())))
		})

		t.Run("OuterEmpty", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[string](),
				OptFlatten(OptEmpty[Opt[string]]()),
			)
		})
	})

	t.Run("OptFlatten", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				OptOf("hello"),
				OptFlattenPtr(OptOf(Ptr("hello"))),
			)
		})

		t.Run("InnerNil", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[string](),
				OptFlattenPtr(OptOf[*string](nil)),
			)
		})

		t.Run("OuterEmpty", func(t *testing.T) {
			assert.Equal(t,
				OptEmpty[string](),
				OptFlattenPtr(OptEmpty[*string]()),
			)
		})
	})
}
