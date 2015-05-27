package core6502

import (
	"errors"
	"strconv"
	"strings"
)

func Split(s, charset string) []string {
	res := []string{}
	tokenStart := -1

	for i, r := range s {
		if strings.ContainsRune(charset, r) {
			if tokenStart != -1 {
				res = append(res, s[tokenStart:i])
				tokenStart = -1
			}
		} else {
			if tokenStart == -1 {
				tokenStart = i
			}
		}
	}
	if tokenStart != -1 {
		res = append(res, s[tokenStart:])
	}
	return res
}

func ParseUint(s string, bitSize int) (uint64, error) {
	s = strings.Trim(s, " \t")

	if len(s) == 0 {
		return 0, errors.New("Empty input")
	}

	var ui uint64
	var err error

	switch s[0] {
	case '$':
		ui, err = strconv.ParseUint(s[1:], 16, bitSize)
	case '-':
		i, err := strconv.ParseInt(s, 10, bitSize)
		return uint64(i), err
	default:
		ui, err = strconv.ParseUint(s, 10, bitSize)
	}
	return ui, err
}
