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

func LogicalShiftRight8(a uint8, carry bool) (uint8, bool) {
	res := uint16(a) << 1
	if carry {
		res |= 1
	}

	if res > 0xff {
		carry = true
	}

	return uint8(res), carry
}
