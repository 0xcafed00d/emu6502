package core6502

func SignExtend8To16(val uint8) uint16 {
	return uint16(int16(int8(val)))
}

func AddWithCarry8(a, b uint8, carry bool) (uint8, bool) {
	res := SignExtend8To16(a) + SignExtend8To16(b)
	if carry {
		res++
	}

	if res > 0xff {
		carry = true
	}

	return uint8(res), carry
}

func SubWithCarry8(a, b uint8, carry bool) (uint8, bool) {
	res := SignExtend8To16(a) - SignExtend8To16(b)
	if carry {
		res--
	}

	if res > 0xff {
		carry = true
	}

	return uint8(res), carry
}

// carry << val << carry
func LogicalShiftLeft8(a uint8, carry bool) (uint8, bool) {
	res := uint16(a) << 1
	if carry {
		res |= 1
	}

	if res > 0xff {
		carry = true
	}

	return uint8(res), carry
}

// carry >> val >> carry
func LogicalShiftRight8(a uint8, carry bool) (uint8, bool) {
	res := a >> 1
	if carry {
		res |= 0x80
	}

	if a&1 != 0 {
		carry = true
	}

	return res, carry
}

// msb(val) >> val >> carry
func ArithmeticShiftRight8(a uint8) (uint8, bool) {
	res := int8(a) >> 1

	carry := false
	if a&1 != 0 {
		carry = true
	}

	return uint8(res), carry
}
