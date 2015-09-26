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
	pack := assert.Pack

	assert.Equal(t, pack(ParseUint("1234", 16))[0], uint64(1234))

	assert.NoError(t, pack(ParseUint("1234", 16)))
	assert.HasError(t, pack(ParseUint("1234", 8)))
}
