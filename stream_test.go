package ef

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	t.Run("StreamOf", func(t *testing.T) {
		t.Run("Slice", func(t *testing.T) {
			st := StreamOf[int]([]int{1, 2, 3})
			assert.Equal(t, []int{1, 2, 3}, st.ToList())
		})

		t.Run("Pointer", func(t *testing.T) {
			t.Run("Value", func(t *testing.T) {
				st := StreamOf[int](Ptr(1))
				assert.Equal(t, []int{1}, st.ToList())
			})

			t.Run("Nil", func(t *testing.T) {
				var value *int
				st := StreamOf[int](value)
				assert.Equal(t, []int{}, st.ToList())
			})
		})

		t.Run("Optional", func(t *testing.T) {
			t.Run("Value", func(t *testing.T) {
				st := StreamOf[int](OptOf(1))
				assert.Equal(t, []int{1}, st.ToList())
			})

			t.Run("Empty", func(t *testing.T) {
				st := StreamOf[int](OptEmpty[int]())
				assert.Equal(t, []int{}, st.ToList())
			})
		})

		t.Run("Stream", func(t *testing.T) {
			st1 := StreamOf[int]([]int{1, 2, 3})
			st2 := StreamOf[int](st1)
			assert.Equal(t, []int{1, 2, 3}, st2.ToList())
		})

		// note [bs]: now seems like a good time to dive a littler deeper into
		// certain nil behavior w/ functions.
		t.Run("Func", func(t *testing.T) {
			counter := 1
			fn := func() Opt[int] {
				if counter > 3 {
					return OptEmpty[int]()
				}
				val := counter
				counter++
				return OptOf(val)
			}
			st := StreamOf[int](fn)
			assert.Equal(t, []int{1, 2, 3}, st.ToList())
		})
	})
}
