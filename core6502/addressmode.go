package core6502

type AddressMode int

const (
	AddrMode_Implicit         AddressMode = iota
	AddrMode_Immediate        AddressMode = iota
	AddrMode_AbsolutePageZero AddressMode = iota
	AddrMode_Absolute         AddressMode = iota
	AddrMode_ZeroPageIdxX     AddressMode = iota
	AddrMode_ZeroPageIdxY     AddressMode = iota
)

type AddrModeReadFunc func(ctx CPUContext) (uint8, int)
type AddrModeWriteFunc func(ctx CPUContext, val uint8) int

func GetReadFunc(mode AddressMode) AddrModeReadFunc {
	switch mode {
	case AddrMode_Implicit:
		return nil
	case AddrMode_Immediate:
		return ReadImmediate
	case AddrMode_Absolute:
		return ReadAbsolute
	case AddrMode_AbsolutePageZero:
		return ReadAbsoluteZeroPage
	case AddrMode_ZeroPageIdxX:
		return ReadZeroPageIdxX
	case AddrMode_ZeroPageIdxY:
		return ReadZeroPageIdxY
	}

	panic("Invalid Address Mode")
}

func GetWriteFunc(mode AddressMode) AddrModeWriteFunc {
	switch mode {
	case AddrMode_Absolute:
		return WriteAbsolute
	case AddrMode_AbsolutePageZero:
		return WriteAbsoluteZeroPage
	case AddrMode_ZeroPageIdxX:
		return WriteZeroPageIdxX
	case AddrMode_ZeroPageIdxY:
		return WriteZeroPageIdxY
	}

	panic("Invalid Address Mode")
}

// #$ff
func ReadImmediate(ctx CPUContext) (uint8, int) {
	return ctx.Peek(ctx.RegPC() + 1), 0
}

// $ff
func ReadAbsoluteZeroPage(ctx CPUContext) (uint8, int) {
	return ctx.Peek(uint16(ctx.Peek(ctx.RegPC() + 1))), 0
}

// $ff
func WriteAbsoluteZeroPage(ctx CPUContext, val uint8) int {
	ctx.Poke(uint16(ctx.Peek(ctx.RegPC()+1)), val)
	return 0
}

// $ffff
func ReadAbsolute(ctx CPUContext) (uint8, int) {
	return ctx.Peek(ctx.PeekWord(ctx.RegPC() + 1)), 0
}

// $ffff
func WriteAbsolute(ctx CPUContext, val uint8) int {
	ctx.Poke(ctx.PeekWord(ctx.RegPC()+1), val)
	return 0
}

// $ff, x
func ReadZeroPageIdxX(ctx CPUContext) (uint8, int) {
	return ctx.Peek(uint16(ctx.Peek(ctx.RegPC()+1) + ctx.RegX())), 0
}

// $ff, x
func WriteZeroPageIdxX(ctx CPUContext, val uint8) int {
	ctx.Poke(uint16(ctx.Peek(ctx.RegPC()+1)+ctx.RegX()), val)
	return 0
}

// $ff, y
func ReadZeroPageIdxY(ctx CPUContext) (uint8, int) {
	return ctx.Peek(uint16(ctx.Peek(ctx.RegPC()+1) + ctx.RegY())), 0
}

// $ff, y
func WriteZeroPageIdxY(ctx CPUContext, val uint8) int {
	ctx.Poke(uint16(ctx.Peek(ctx.RegPC()+1)+ctx.RegY()), val)
	return 0
}
