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
	length    uint16
	tstates   int
	mode      AddressMode
}

var InstructionData = []InstructionInfo{
	{0xa9, "LDA", LDA, 2, 2, AddrMode_Immediate},
	{0xa5, "LDA", LDA, 2, 3, AddrMode_AbsoluteZeroPage},
	{0xb5, "LDA", LDA, 2, 4, AddrMode_ZeroPageIdxX},
	{0xad, "LDA", LDA, 3, 4, AddrMode_Absolute},
	{0x85, "STA", STA, 2, 3, AddrMode_AbsoluteZeroPage},
	{0x18, "CLC", CLC, 1, 2, AddrMode_Implicit},
	{0x38, "SEC", SEC, 1, 2, AddrMode_Implicit},
	{0xaa, "TAX", TAX, 1, 2, AddrMode_Implicit},
	{0xa8, "TAY", TAY, 1, 2, AddrMode_Implicit},
	{0xba, "TSX", TSX, 1, 2, AddrMode_Implicit},
	{0x8a, "TXA", TXA, 1, 2, AddrMode_Implicit},
	{0x9a, "TXS", TXS, 1, 2, AddrMode_Implicit},
	{0x98, "TYA", TYA, 1, 2, AddrMode_Implicit},
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

func LDA(info *InstructionInfo) InstructionExecFunc {
	readFunc := GetReadFunc(info.mode)

	return func(ctx CPUContext) int {
		val, exclock := readFunc(ctx)
		ctx.SetRegA(setFlagsFromValue(ctx, val))
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates + exclock
	}
}

func STA(info *InstructionInfo) InstructionExecFunc {
	writeFunc := GetWriteFunc(info.mode)

	return func(ctx CPUContext) int {
		exclock := writeFunc(ctx, ctx.RegA())
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates + exclock
	}
}

func CLC(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_C, false)
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}

func SEC(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_C, true)
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}

func TAX(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegA()))
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}

func TAY(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetRegY(setFlagsFromValue(ctx, ctx.RegA()))
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}

func TSX(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegSP()))
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}

func TXA(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetRegA(setFlagsFromValue(ctx, ctx.RegX()))
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}

func TXS(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetRegSP(ctx.RegX())
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}

func TYA(info *InstructionInfo) InstructionExecFunc {
	return func(ctx CPUContext) int {
		ctx.SetRegA(setFlagsFromValue(ctx, ctx.RegY()))
		ctx.SetRegPC(ctx.RegPC() + info.length)
		return info.tstates
	}
}
