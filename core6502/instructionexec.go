package core6502

import (
	"fmt"
)

/*
	Executes the next instruction from the supplied CPUContext.
	Results of the execution are written back to the provided context
	Returns the number of clock cycles consumed by the instruction
	and an error. An error is returned on the condition of an invalid
	opcode.
*/
func Execute(ctx CPUContext) (int, error) {
	pc := ctx.RegPC()
	opcode := ctx.Peek(pc)

	if executors[opcode] == nil {
		return 0, fmt.Errorf("Invalid Opcode: $%02x @ $%04x", opcode, pc)
	}

	return executors[opcode](ctx), nil
}

type InstructionExecFunc func(ctx CPUContext) int
type ExecFuncMakerFunc func(InstructionInfo *InstructionInfo) InstructionExecFunc

type InstructionInfo struct {
	opcode    uint8
	name      string
	execMaker ExecFuncMakerFunc
	tstates   int
	mode      AddressMode
}

var InstructionData = []InstructionInfo{
	{0xa9, "LDA", LDA, 2, AddrMode_Immediate},
	{0xa5, "LDA", LDA, 3, AddrMode_AbsoluteZeroPage},
	{0xb5, "LDA", LDA, 4, AddrMode_ZeroPageIdxX},
	{0xa1, "LDA", LDA, 6, AddrMode_PostIndexIndirect},
	{0xb1, "LDA", LDA, 5, AddrMode_PreIndexIndirect},
	{0xad, "LDA", LDA, 4, AddrMode_Absolute},
	{0xbd, "LDA", LDA, 4, AddrMode_AbsoluteIndexedX},
	{0xb9, "LDA", LDA, 4, AddrMode_AbsoluteIndexedY},

	{0xa2, "LDX", LDX, 2, AddrMode_Immediate},
	{0xa6, "LDX", LDX, 3, AddrMode_AbsoluteZeroPage},
	{0xb6, "LDX", LDX, 4, AddrMode_ZeroPageIdxY},
	{0xae, "LDX", LDX, 4, AddrMode_Absolute},
	{0xbe, "LDX", LDX, 4, AddrMode_AbsoluteIndexedY},

	{0xa0, "LDY", LDY, 2, AddrMode_Immediate},
	{0xa4, "LDY", LDY, 3, AddrMode_AbsoluteZeroPage},
	{0xb4, "LDY", LDY, 4, AddrMode_ZeroPageIdxX},
	{0xac, "LDY", LDY, 4, AddrMode_Absolute},
	{0xbc, "LDY", LDY, 4, AddrMode_AbsoluteIndexedX},

	{0x85, "STA", STA, 3, AddrMode_AbsoluteZeroPage},
	{0x95, "STA", STA, 4, AddrMode_ZeroPageIdxX},
	{0x81, "STA", STA, 6, AddrMode_PostIndexIndirect},
	{0x91, "STA", STA, 6, AddrMode_PreIndexIndirect},
	{0x8d, "STA", STA, 4, AddrMode_Absolute},
	{0x9d, "STA", STA, 5, AddrMode_AbsoluteIndexedX},
	{0x99, "STA", STA, 5, AddrMode_AbsoluteIndexedY},

	{0x86, "STX", STX, 3, AddrMode_AbsoluteZeroPage},
	{0x96, "STX", STX, 4, AddrMode_ZeroPageIdxY},
	{0x8e, "STX", STX, 4, AddrMode_Absolute},

	{0x84, "STY", STY, 3, AddrMode_AbsoluteZeroPage},
	{0x94, "STY", STY, 4, AddrMode_ZeroPageIdxX},
	{0x8c, "STY", STY, 4, AddrMode_Absolute},

	// todo tstates
	{0xe6, "INC", INC, 4, AddrMode_AbsoluteZeroPage},
	{0xf6, "INC", INC, 4, AddrMode_ZeroPageIdxX},
	{0xee, "INC", INC, 4, AddrMode_Absolute},
	{0xfe, "INC", INC, 4, AddrMode_AbsoluteIndexedX},

	{0xe8, "INX", INX, 4, AddrMode_Implicit},
	{0xc8, "INY", INY, 4, AddrMode_Implicit},

	{0xc6, "DEC", DEC, 4, AddrMode_AbsoluteZeroPage},
	{0xd6, "DEC", DEC, 4, AddrMode_ZeroPageIdxX},
	{0xce, "DEC", DEC, 4, AddrMode_Absolute},
	{0xde, "DEC", DEC, 4, AddrMode_AbsoluteIndexedX},

	{0xca, "DEX", DEX, 4, AddrMode_Implicit},
	{0x88, "DEY", DEY, 4, AddrMode_Implicit},

	{0x18, "CLC", CLC, 2, AddrMode_Implicit},
	{0x38, "SEC", SEC, 2, AddrMode_Implicit},
	{0xaa, "TAX", TAX, 2, AddrMode_Implicit},
	{0xa8, "TAY", TAY, 2, AddrMode_Implicit},
	{0xba, "TSX", TSX, 2, AddrMode_Implicit},
	{0x8a, "TXA", TXA, 2, AddrMode_Implicit},
	{0x9a, "TXS", TXS, 2, AddrMode_Implicit},
	{0x98, "TYA", TYA, 2, AddrMode_Implicit},
}

var executors [256]InstructionExecFunc

func init() {
	for n := 0; n < len(InstructionData); n++ {
		info := &InstructionData[n]
		executors[info.opcode] = info.execMaker(info)
	}
}

func setFlagsFromValue(ctx CPUContext, val uint8) uint8 {
	ctx.SetFlag(Flag_Z, val == 0)
	ctx.SetFlag(Flag_N, (val&0x80) == 0x80)
	return val
}

func INC(info *InstructionInfo) InstructionExecFunc {
	readFunc := GetReadFunc(info.mode)
	writeFunc := GetWriteFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		val, exclock := readFunc(ctx)
		writeFunc(ctx, setFlagsFromValue(ctx, val+1))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates + exclock
	}
}

func INX(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegX()+1))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func INY(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegY(setFlagsFromValue(ctx, ctx.RegY()+1))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func DEC(info *InstructionInfo) InstructionExecFunc {
	readFunc := GetReadFunc(info.mode)
	writeFunc := GetWriteFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		val, exclock := readFunc(ctx)
		writeFunc(ctx, setFlagsFromValue(ctx, val-1))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates + exclock
	}
}

func DEX(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegX()-1))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func DEY(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegY(setFlagsFromValue(ctx, ctx.RegY()-1))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func LDA(info *InstructionInfo) InstructionExecFunc {
	readFunc := GetReadFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		val, exclock := readFunc(ctx)
		ctx.SetRegA(setFlagsFromValue(ctx, val))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates + exclock
	}
}

func LDX(info *InstructionInfo) InstructionExecFunc {
	readFunc := GetReadFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		val, exclock := readFunc(ctx)
		ctx.SetRegY(setFlagsFromValue(ctx, val))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates + exclock
	}
}

func LDY(info *InstructionInfo) InstructionExecFunc {
	readFunc := GetReadFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		val, exclock := readFunc(ctx)
		ctx.SetRegY(setFlagsFromValue(ctx, val))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates + exclock
	}
}

func STA(info *InstructionInfo) InstructionExecFunc {
	writeFunc := GetWriteFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		writeFunc(ctx, ctx.RegA())
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func STX(info *InstructionInfo) InstructionExecFunc {
	writeFunc := GetWriteFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		writeFunc(ctx, ctx.RegX())
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func STY(info *InstructionInfo) InstructionExecFunc {
	writeFunc := GetWriteFunc(info.mode)
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		writeFunc(ctx, ctx.RegY())
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func CLC(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_C, false)
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func SEC(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_C, true)
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func TAX(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegA()))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func TAY(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegY(setFlagsFromValue(ctx, ctx.RegA()))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func TSX(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegSP()))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func TXA(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegA(setFlagsFromValue(ctx, ctx.RegX()))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func TXS(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegSP(ctx.RegX())
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func TYA(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegA(setFlagsFromValue(ctx, ctx.RegY()))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}
