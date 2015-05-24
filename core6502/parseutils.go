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

func ParseInt(s string, bitSize int) (int64, error) {
	s = strings.Trim(s, " \t")

	if len(s) == 0 {
		return 0, errors.New("Empty input")
	}

	if s[0] == '$' {
		// hex value
		i, err := strconv.ParseInt(s[1:], 16, bitSize)
		return i, err
	} else {
		// decimal value
		i, err := strconv.ParseInt(s, 10, bitSize)
		return i, err
	}
}

func ParseUint(s string, bitSize int) (uint64, error) {
	s = strings.Trim(s, " \t")

	if len(s) == 0 {
		return 0, errors.New("Empty input")
	}

	if s[0] == '$' {
		// hex value
		i, err := strconv.ParseUint(s[1:], 16, bitSize)
		return i, err
	} else {
		// decimal value
		i, err := strconv.ParseUint(s, 10, bitSize)
		return i, err
	}
}
