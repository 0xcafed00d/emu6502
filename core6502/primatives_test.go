package core6502

import (
	"testing"
)

type bintestdata struct {
	dataVal1  uint8
	dataVal2  uint8
	dataCarry bool

	expectedVal   uint8
	expectedCarry bool
}

type unarytestdata struct {
	dataVal   uint8
	dataCarry bool

	expectedVal   uint8
	expectedCarry bool
}

type BinaryFunction func(a, b uint8, carry bool) (uint8, bool)
type UnaryFunction func(a uint8, carry bool) (uint8, bool)

func testBinaryFunction(t *testing.T, data []bintestdata, binFunc BinaryFunction) {
	for i, tst := range data {
		val, carry := binFunc(tst.dataVal1, tst.dataVal2, tst.dataCarry)
		if val != tst.expectedVal || carry != tst.expectedCarry {
			t.Fatalf("Failed: [%d] Expected: [%v,%v] Got: [%v,%v]",
				i, tst.expectedVal, tst.expectedCarry, val, carry)
		}
	}
}

func testUnaryFunction(t *testing.T, data []unarytestdata, unaryFunc UnaryFunction) {
	for i, tst := range data {
		val, carry := unaryFunc(tst.dataVal, tst.dataCarry)
		if val != tst.expectedVal || carry != tst.expectedCarry {
			t.Fatalf("Failed: [%d] Expected: [%v,%v] Got: [%v,%v]",
				i, tst.expectedVal, tst.expectedCarry, val, carry)
		}
	}
}

func TestAddWithCarry(t *testing.T) {
	var test = []bintestdata{
		{0x10, 0x10, false, 0x20, false},
		{0xff, 0x00, false, 0xff, false},
		{0xfe, 0x00, true, 0xff, false},
		{0xff, 0x00, true, 0x00, true},
		{0xff, 0x01, false, 0x00, true},
		{0xfd, 0x01, true, 0xff, false},
		{0xfe, 0x01, true, 0x00, true},
		{0xff, 0x01, true, 0x01, true},
	}
	testBinaryFunction(t, test, AddWithCarry8)
}

func TestSubWithCarry(t *testing.T) {
	var test = []bintestdata{
		{0x10, 0x10, false, 0x00, false},
	}
	testBinaryFunction(t, test, SubWithCarry8)
}

func TestLSL(t *testing.T) {
	var test = []unarytestdata{
		{B_01000111, false, B_10001110, false},
		{B_01000111, true, B_10001111, false},
		{B_11000111, false, B_10001110, true},
		{B_11000111, true, B_10001111, true},
	}

	testUnaryFunction(t, test, LogicalShiftLeft8)
}

func TestLSR(t *testing.T) {
	var test = []unarytestdata{
		{B_10001110, false, B_01000111, false},
		{B_10001110, true, B_11000111, false},
		{B_11000111, false, B_01100011, true},
		{B_11000111, true, B_11100011, true},
	}

	testUnaryFunction(t, test, LogicalShiftRight8)
}
