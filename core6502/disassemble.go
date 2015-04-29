package core6502

import (
	"fmt"
)

type disasmInfo struct {
	name   string
	length uint16
	mode   AddressMode
}

var disasmData [256]disasmInfo

func init() {
	for n := 0; n < len(InstructionData); n++ {
		info := &InstructionData[n]
		disasmData[info.opcode] = disasmInfo{info.name, info.length, info.mode}
	}
}

func addressModeToStr(mode AddressMode, ctx CPUContext, addr uint16) string {
	switch mode {
	case AddrMode_Immediate:
		return fmt.Sprintf("#$%02x", ctx.Peek(addr))
	case AddrMode_Implicit:
		return ""
	case AddrMode_Absolute:
		return fmt.Sprintf("$%04x", ctx.PeekWord(addr))
	case AddrMode_AbsoluteZeroPage:
		return fmt.Sprintf("$%02x", ctx.Peek(addr))
	case AddrMode_ZeroPageIdxX:
		return fmt.Sprintf("$%02x, X", ctx.Peek(addr))
	case AddrMode_ZeroPageIdxY:
		return fmt.Sprintf("$%02x, X", ctx.Peek(addr))
	case AddrMode_PreIndexIndirect:
		return fmt.Sprintf("($%02x, X)", ctx.Peek(addr))
	case AddrMode_PostIndexIndirect:
		return fmt.Sprintf("($%02x), Y", ctx.Peek(addr))
	case AddrMode_AbsoluteIndexedX:
		return fmt.Sprintf("$%04x, X", ctx.PeekWord(addr))
	case AddrMode_AbsoluteIndexedY:
		return fmt.Sprintf("$%04x, Y", ctx.PeekWord(addr))
	case AddrMode_Indirect:
		return fmt.Sprintf("($%04x)", ctx.PeekWord(addr))
	case AddrMode_Relative:
		return fmt.Sprintf("%d", int(ctx.Peek(addr)))
	}
	return "Invalid"
}

func Disassemble(ctx CPUContext, addr uint16) (string, uint16) {
	info := &disasmData[ctx.Peek(addr)]

	if info.mode == AddrMode_Invalid {
		return fmt.Sprintf("db  $%02x", ctx.Peek(addr)), 1
	}
	return info.name + " " + addressModeToStr(info.mode, ctx, addr+1), info.length
}
