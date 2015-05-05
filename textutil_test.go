package main

import (
	"testing"
	"unicode"
)

type testtrim struct {
	data     string
	expected string
	f        func([]rune, func(rune) bool) []rune
}

var testdata_trim = []testtrim{
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
	for _, tst := range testdata_trim {
		got := tst.f([]rune(tst.data), unicode.IsSpace)
		if string(got) != tst.expected {
			t.Fatalf("Trim Test: [%s] Expected: [%s] Got: [%s]",
				tst.data, tst.expected, string(got))
		}
	}
}

type testinsertdel struct {
	data            string
	insert_expected string
	delete_expected string
	index           int
}

var testdata_id = []testinsertdel{
	{"123", "x123", "23", 0},
	{"123", "1x23", "13", 1},
	{"123", "12x3", "12", 2},
	{"123", "123x", "ignore", 3},
}

func TestInsertDelete(t *testing.T) {
	for i, tst := range testdata_id {
		data := CloneRuneSlice([]rune(tst.data))

		if tst.insert_expected != "ignore" {
			res := string(InsertRuneAt(data, 'x', tst.index))
			if res != tst.insert_expected {
				t.Fatalf("Insert Rune Test: [%d] Expected: [%s] Got: [%s]",
					i, tst.insert_expected, res)
			}
		}

		data = CloneRuneSlice([]rune(tst.data))
		if tst.delete_expected != "ignore" {
			res := string(DeleteRuneAt(data, tst.index))
			if res != tst.delete_expected {
				t.Fatalf("Delete Rune Test: [%d] Expected: [%s] Got: [%s]",
					i, tst.delete_expected, res)
			}
		}
	}
}
