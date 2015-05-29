package core6502

import (
	"github.com/simulatedsimian/testbuddy"
	"testing"
)

type T testing.T

func TestSplit(tst *testing.T) {
	t := (*testbuddy.T)(tst)

	t.Assert(Split("", " \t")).Equal([]string{})
	t.Assert(Split("   \t\t\t", " \t")).Equal([]string{})
	t.Assert(Split("hello world", " \t")).Equal([]string{"hello", "world"})
	t.Assert(Split("    hello     world     ", " \t")).Equal([]string{"hello", "world"})
}

func TestParseUint(tst *testing.T) {
	t := (*testbuddy.T)(tst)

	t.AssertNoErr(ParseUint("1234", 16)).Equal(uint64(1234))
	t.AssertErr(ParseUint("1234", 8)).Equal(uint64(0xffffffffffffffff))
}
