package core6502

import (
	"testing"
)

type testdata struct {
	dataVal1  uint8
	dataVal2  uint8
	dataCarry bool

	expectedVal   uint8
	expectedCarry bool
}

var testAddWithCarry = []testdata{
	{10, 10, false, 20, true},
}

func TestAddWithCarry(t *testing.T) {
	for i, tst := range testAddWithCarry {
		val, carry := AddWithCarry8(tst.dataVal1, tst.dataVal2, tst.dataCarry)
		if val != tst.expectedVal || carry != tst.expectedCarry {
			t.Fatalf("Insert Rune Test: [%d] Expected: [%v,%v] Got: [%v,%v]",
				i, tst.expectedVal, tst.expectedCarry, val, carry)
		}
	}
}
