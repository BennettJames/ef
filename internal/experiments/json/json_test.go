package json

import (
	"testing"

	"github.com/BennettJames/ef"
	"github.com/stretchr/testify/assert"
)

func Test_scratch(t *testing.T) {

	// note [bs]: this is a little convoluted, but is interesting. Wonder if you
	// could make an "ef flavored testing suite" for fun after, but let's not
	// worry about that quite yet (also - might be able to make some nifty
	// table-style tests that make use )

	type Foo struct {
		Value int
		Key   string
	}

	t.Run("Valid", func(t *testing.T) {
		const input = `
		{
			"Value": 22,
			"Key": "key-one"
		}
	`

		value := ef.Try(Parse[Foo](ToReader(input)))
		assert.Equal(t, Foo{Value: 22, Key: "key-one"}, value)
	})

	t.Run("Invalid", func(t *testing.T) {
		const input = `
		{
			"Value": 22,
			"Key": "key-one
		}
	`

		_, err := Parse[Foo](ToReader(input))
		assert.NotNil(t, err)
	})
}

type TableTest[T, U any] struct {
	Name     string
	Input    T
	Expected U
}

func RunTableTest[T, U any](t *testing.T, cases []TableTest[T, U], fn func(T) U) {
	// note [bs]: this is interesting. Feels like the matching / equivalence is a
	// little narrow, esp. for cases involving errors. Still, an interesting path to
	// go down at a later point.

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			out := fn(c.Input)
			assert.Equal(t, c.Expected, out)
		})
	}
}
