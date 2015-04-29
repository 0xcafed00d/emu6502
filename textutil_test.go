package main

import (
	"testing"
	"unicode"
)

type test struct {
	data     string
	expected string
	f        func([]rune, func(rune) bool) []rune
}

var testdata = []test{
	{"123", "123", TrimRunesLeft},
	{" 123", "123", TrimRunesLeft},
	{"  123", "123", TrimRunesLeft},
	{"   ", "", TrimRunesLeft},
	{"", "", TrimRunesLeft},

	{"456", "456", TrimRunesRight},
	{"456 ", "456", TrimRunesRight},
	{"456  ", "456", TrimRunesRight},
	{"   ", "", TrimRunesRight},
	{"", "", TrimRunesRight},

	{"789", "789", TrimRunes},
	{"789 ", "789", TrimRunes},
	{"789  ", "789", TrimRunes},
	{" 789", "789", TrimRunes},
	{"  789", "789", TrimRunes},
	{" 789 ", "789", TrimRunes},
	{"  789  ", "789", TrimRunes},
	{"   ", "", TrimRunes},
	{"", "", TrimRunes},
}

func TestTrim(t *testing.T) {
	for _, tst := range testdata {
		got := tst.f([]rune(tst.data), unicode.IsSpace)
		if string(got) != tst.expected {
			t.Fatalf("Test: [%s] Expected: [%s] Got: [%s]",
				tst.data, tst.expected, string(got))
		}
	}

}
