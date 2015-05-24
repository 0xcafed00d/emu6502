package core6502

import (
	"github.com/simulatedsimian/testbuddy"
	"testing"
)

func TestSplit(t *testing.T) {
	testbuddy.AssertEqual(t, Split("", " \t"), []string{})
	testbuddy.AssertEqual(t, Split("   \t\t\t", " \t"), []string{})
	testbuddy.AssertEqual(t, Split("hello world", " \t"), []string{"hello", "world"})
	testbuddy.AssertEqual(t, Split("    hello     world     ", " \t"), []string{"hello", "world"})
}
