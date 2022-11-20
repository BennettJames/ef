package opt

import (
	"testing"

	"github.com/BennettJames/ef"
	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {

	t.Run("Of", func(t *testing.T) {
		assert.Equal(t,
			ef.NewOptValue("hello"),
			Of("hello"))
	})

	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t,
			ef.Opt[string]{},
			Empty[string]())
	})

	t.Run("OfPtr", func(t *testing.T) {
		t.Run("Val", func(t *testing.T) {
			strVal := "hello"
			var strPtr *string = &strVal
			assert.Equal(t,
				ef.NewOptValue(strVal),
				OfPtr(strPtr))
		})

		t.Run("Empty", func(t *testing.T) {
			var strPtr *string
			assert.Equal(t,
				ef.Opt[string]{},
				OfPtr(strPtr))
		})
	})

	t.Run("OfOk", func(t *testing.T) {
		t.Run("Ok", func(t *testing.T) {
			assert.Equal(t,
				ef.Opt[string]{},
				OfOk("hello", true))
		})

		t.Run("NotOk", func(t *testing.T) {
			assert.Equal(t,
				ef.Opt[string]{},
				OfOk("hello", false))
		})
	})

	t.Run("SliceGet", func(t *testing.T) {
		s := []string{"hello", "world"}
		t.Run("InBounds", func(t *testing.T) {
			assert.Equal(t,
				Of("hello"),
				SliceGet(s, 0))
			assert.Equal(t,
				Of("world"),
				SliceGet(s, 1))
		})

		t.Run("OutOfBounds", func(t *testing.T) {
			assert.Equal(t,
				Empty[string](),
				SliceGet(s, -1))
			assert.Equal(t,
				Empty[string](),
				SliceGet(s, 2))
		})
	})

	t.Run("MapGet", func(t *testing.T) {
		m := map[string]string{
			"hello": "world",
		}
		t.Run("Ok", func(t *testing.T) {
			assert.Equal(t,
				ef.NewOptValue("world"),
				MapGet(m, "hello"))
		})

		t.Run("NotOk", func(t *testing.T) {
			assert.Equal(t,
				ef.Opt[string]{},
				MapGet(m, "badkey"))
		})
	})

	t.Run("MapGet", func(t *testing.T) {
		m := map[string]string{
			"hello": "world",
		}
		t.Run("Ok", func(t *testing.T) {
			assert.Equal(t,
				ef.NewOptValue("world"),
				MapGet(m, "hello"))
		})

		t.Run("NotOk", func(t *testing.T) {
			assert.Equal(t,
				ef.Opt[string]{},
				MapGet(m, "badkey"))
		})
	})

	t.Run("Map", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				Of(22),
				Map(Of("hello"), func(v string) int {
					assert.Equal(t, "hello", v)
					return 22
				}))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				Empty[int](),
				Map(Empty[string](), func(v string) int {
					panic("unreachable")
				}))
		})
	})

	t.Run("FlatMap", func(t *testing.T) {
		t.Run("ValueToValue", func(t *testing.T) {
			assert.Equal(t,
				Of(22),
				FlatMap(Of("hello"), func(v string) ef.Opt[int] {
					assert.Equal(t, "hello", v)
					return Of(22)
				}))
		})

		t.Run("ValueToEmpty", func(t *testing.T) {
			assert.Equal(t,
				Empty[int](),
				FlatMap(Of("hello"), func(v string) ef.Opt[int] {
					assert.Equal(t, "hello", v)
					return Empty[int]()
				}))
		})

		t.Run("Empty", func(t *testing.T) {
			assert.Equal(t,
				Empty[int](),
				FlatMap(Empty[string](), func(v string) ef.Opt[int] {
					panic("unreachable")
				}))
		})
	})

	t.Run("Flatten", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				Of("hello"),
				Flatten(Of(Of("hello"))))
		})

		t.Run("InnerEmpty", func(t *testing.T) {
			assert.Equal(t,
				Empty[string](),
				Flatten(Of(Empty[string]())))
		})

		t.Run("OuterEmpty", func(t *testing.T) {
			assert.Equal(t,
				Empty[string](),
				Flatten(Empty[ef.Opt[string]]()),
			)
		})
	})

	t.Run("Flatten", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			assert.Equal(t,
				Of("hello"),
				Deref(Of(ef.Ptr("hello"))),
			)
		})

		t.Run("InnerNil", func(t *testing.T) {
			assert.Equal(t,
				Empty[string](),
				Deref(Of[*string](nil)),
			)
		})

		t.Run("OuterEmpty", func(t *testing.T) {
			assert.Equal(t,
				Empty[string](),
				Deref(Empty[*string]()),
			)
		})
	})
}
