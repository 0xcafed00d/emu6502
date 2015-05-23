package core6502

import (
	"errors"
	"strconv"
	"strings"
)

func ParseInt(s string, bitSize int) (int64, error) {
	s = strings.Trim(s, " \t")

	if len(s) == 0 {
		return 0, errors.New("Empty input")
	}

	// hex value
	if s[0] == '$' {
		i, err := strconv.ParseInt(s[1:], 16, bitSize)
		return i, err
	} else {
		i, err := strconv.ParseInt(s, 10, bitSize)
		return i, err
	}
}

func ParseUint(s string, bitSize int) (uint64, error) {
	s = strings.Trim(s, " \t")

	if len(s) == 0 {
		return 0, errors.New("Empty input")
	}

	// hex value
	if s[0] == '$' {
		i, err := strconv.ParseUint(s[1:], 16, bitSize)
		return i, err
	} else {
		i, err := strconv.ParseUint(s, 10, bitSize)
		return i, err
	}
}
