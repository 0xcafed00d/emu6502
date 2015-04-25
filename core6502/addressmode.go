package core6502

type AddrModeFunc func(ctx CPUContext, val uint8) (uint8, int)

// #$ff
func ReadImmediate(ctx CPUContext, _ uint8) (uint8, int) {
	return ctx.Peek(ctx.RegPC() + 1), 0
}

// $ff
func ReadAbsoluteZeroPage(ctx CPUContext, _ uint8) (uint8, int) {
	return ctx.Peek(uint16(ctx.Peek(ctx.RegPC() + 1))), 0
}

// $ff
func WriteAbsoluteZeroPage(ctx CPUContext, val uint8) (uint8, int) {
	ctx.Poke(uint16(ctx.Peek(ctx.RegPC()+1)), val)
	return 0, 0
}

// $ffff
func ReadAbsolute(ctx CPUContext, _ uint8) (uint8, int) {
	return ctx.Peek(ctx.PeekWord(ctx.RegPC() + 1)), 0
}

// $ffff
func WriteAbsolute(ctx CPUContext, val uint8) (uint8, int) {
	ctx.Poke(ctx.PeekWord(ctx.RegPC()+1), val)
	return 0, 0
}

// $ff, x
func ReadZeroPageIdxX(ctx CPUContext, _ uint8) (uint8, int) {
	return ctx.Peek(uint16(ctx.Peek(ctx.RegPC()+1) + ctx.RegX())), 0
}

// $ff, x
func WriteZeroPageIdxX(ctx CPUContext, val uint8) (uint8, int) {
	ctx.Poke(uint16(ctx.Peek(ctx.RegPC()+1)+ctx.RegX()), val)
	return 0, 0
}

// $ff, y
func ReadZeroPageIdxY(ctx CPUContext, _ uint8) (uint8, int) {
	return ctx.Peek(uint16(ctx.Peek(ctx.RegPC()+1) + ctx.RegY())), 0
}

// $ff, y
func WriteZeroPageIdxY(ctx CPUContext, val uint8) (uint8, int) {
	ctx.Poke(uint16(ctx.Peek(ctx.RegPC()+1)+ctx.RegY()), val)
	return 0, 0
}
