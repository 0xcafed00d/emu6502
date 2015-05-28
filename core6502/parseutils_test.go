package core6502

import (
	"github.com/simulatedsimian/testbuddy"
	"testing"
)

type T testing.T

func (t *T) AssertEqualsNoErr(val interface{}, err error) func(val interface{}) {
	return func(val interface{}) {

	}
}

func (t *T) AssertErr(val interface{}, err error) {
	testbuddy.AssertEqual((*testing.T)(t), val, nil)
}

func TestSplit(t *testing.T) {
	tt := (*T)(t)

	testbuddy.AssertEqual(t, Split("", " \t"), []string{})
	testbuddy.AssertEqual(t, Split("   \t\t\t", " \t"), []string{})
	testbuddy.AssertEqual(t, Split("hello world", " \t"), []string{"hello", "world"})
	testbuddy.AssertEqual(t, Split("    hello     world     ", " \t"), []string{"hello", "world"})

	tt.AssertEqualsNoErr(ParseUint("123", 16))(1)
}
