package core6502

type AddressMode int

const (
	AddrMode_Invalid AddressMode = iota
	AddrMode_Implicit
	AddrMode_Immediate
	AddrMode_AbsoluteZeroPage
	AddrMode_Absolute
	AddrMode_ZeroPageIdxX
	AddrMode_ZeroPageIdxY
	AddrMode_PreIndexIndirect
	AddrMode_PostIndexIndirect
	AddrMode_AbsoluteIndexedX
	AddrMode_AbsoluteIndexedY
	AddrMode_Indirect
	AddrMode_Relative
)

type AddrModeReadFunc func(ctx CPUContext) (uint8, int)
type AddrModeWriteFunc func(ctx CPUContext, val uint8) int

func InstructionBytes(mode AddressMode) uint16 {
	switch mode {
	case AddrMode_Implicit:
		return 1
	case AddrMode_Immediate:
		return 2
	case AddrMode_AbsoluteZeroPage:
		return 2
	case AddrMode_Absolute:
		return 3
	case AddrMode_ZeroPageIdxX:
		return 2
	case AddrMode_ZeroPageIdxY:
		return 2
	case AddrMode_PreIndexIndirect:
		return 2
	case AddrMode_PostIndexIndirect:
		return 2
	case AddrMode_AbsoluteIndexedX:
		return 3
	case AddrMode_AbsoluteIndexedY:
		return 3
	case AddrMode_Indirect:
		return 3
	case AddrMode_Relative:
		return 2
	}

	return 0
}

func GetReadFunc(mode AddressMode) AddrModeReadFunc {
	switch mode {
	case AddrMode_Implicit:
		return nil
	case AddrMode_Immediate:
		return ReadImmediate
	case AddrMode_Absolute:
		return ReadAbsolute
	case AddrMode_AbsoluteZeroPage:
		return ReadAbsoluteZeroPage
	case AddrMode_ZeroPageIdxX:
		return ReadZeroPageIdxX
	case AddrMode_ZeroPageIdxY:
		return ReadZeroPageIdxY
	case AddrMode_PreIndexIndirect:
		return ReadPreIndexIndirect
	case AddrMode_PostIndexIndirect:
		return ReadPostIndexIndirect
	case AddrMode_AbsoluteIndexedX:
		return ReadAboluteIndexedX
	case AddrMode_AbsoluteIndexedY:
		return ReadAboluteIndexedY
	}

	panic("Invalid Address Mode")
}

func GetWriteFunc(mode AddressMode) AddrModeWriteFunc {
	switch mode {
	case AddrMode_Absolute:
		return WriteAbsolute
	case AddrMode_AbsoluteZeroPage:
		return WriteAbsoluteZeroPage
	case AddrMode_ZeroPageIdxX:
		return WriteZeroPageIdxX
	case AddrMode_ZeroPageIdxY:
		return WriteZeroPageIdxY
	case AddrMode_PreIndexIndirect:
		return WritePreIndexIndirect
	case AddrMode_PostIndexIndirect:
		return WritePostIndexIndirect
	case AddrMode_AbsoluteIndexedX:
		return WriteAboluteIndexedX
	case AddrMode_AbsoluteIndexedY:
		return WriteAboluteIndexedY
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

// ($ff, x)
func ReadPreIndexIndirect(ctx CPUContext) (uint8, int) {
	addr := ctx.PeekWord(uint16(ctx.Peek(ctx.RegPC()+1) + ctx.RegX()))
	return ctx.Peek(addr), 0
}

// ($ff, x)
func WritePreIndexIndirect(ctx CPUContext, val uint8) int {
	addr := ctx.PeekWord(uint16(ctx.Peek(ctx.RegPC()+1) + ctx.RegX()))
	ctx.Poke(addr, val)
	return 0
}

// ($ff), y
func ReadPostIndexIndirect(ctx CPUContext) (uint8, int) {
	addr := ctx.PeekWord(uint16(ctx.Peek(ctx.RegPC()+1))) + uint16(ctx.RegY())
	return ctx.Peek(addr), 0
}

// ($ff), y
func WritePostIndexIndirect(ctx CPUContext, val uint8) int {
	addr := ctx.PeekWord(uint16(ctx.Peek(ctx.RegPC()+1))) + uint16(ctx.RegY())
	ctx.Poke(addr, val)
	return 0
}

// $ffff, x
func ReadAboluteIndexedX(ctx CPUContext) (uint8, int) {
	addr := ctx.PeekWord(ctx.PeekWord(ctx.RegPC()+1)) + uint16(ctx.RegX())
	return ctx.Peek(addr), 0
}

// $ffff, x
func WriteAboluteIndexedX(ctx CPUContext, val uint8) int {
	addr := ctx.PeekWord(ctx.PeekWord(ctx.RegPC()+1)) + uint16(ctx.RegX())
	ctx.Poke(addr, val)
	return 0
}

// $ffff, x
func ReadAboluteIndexedY(ctx CPUContext) (uint8, int) {
	addr := ctx.PeekWord(ctx.PeekWord(ctx.RegPC()+1)) + uint16(ctx.RegY())
	return ctx.Peek(addr), 0
}

// $ffff, x
func WriteAboluteIndexedY(ctx CPUContext, val uint8) int {
	addr := ctx.PeekWord(ctx.PeekWord(ctx.RegPC()+1)) + uint16(ctx.RegY())
	ctx.Poke(addr, val)
	return 0
}

func signExtend(val uint8) uint16 {
	return uint16(int16(int8(val)))
}

func CalcPCRelativeAddr(ctx CPUContext) uint16 {
	offset := signExtend(ctx.Peek(ctx.RegPC() + 1))
	return ctx.RegPC() + offset
}
