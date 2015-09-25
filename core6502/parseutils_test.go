package core6502

import (
	"github.com/simulatedsimian/assert"
	"testing"
)

func TestSplit(t *testing.T) {
	assert.Equal(t, Split("", " \t"), []string{})
	assert.Equal(t, Split("   \t\t\t", " \t"), []string{})
	assert.Equal(t, Split("hello world", " \t"), []string{"hello", "world"})
	assert.Equal(t, Split("    hello     world     ", " \t"), []string{"hello", "world"})
}

func TestParseUint(t *testing.T) {
//	t.AssertNoErr(ParseUint("1234", 16)).Equal(uint64(1234))
//	t.AssertErr(ParseUint("1234", 8)).Equal(uint64(0xffffffffffffffff))
}
