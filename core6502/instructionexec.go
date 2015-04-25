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

	if executors[opcode].exec == nil {
		return 0, fmt.Errorf("Invalid Opcode: $%02x @ $%04x", opcode, pc)
	}

	return executors[opcode].exec(ctx, &executors[opcode]), nil
}

type InstructionExecFunc func(ctx CPUContext, execInfo *ExecInfo) int

type ExecInfo struct {
	exec    InstructionExecFunc
	length  uint16
	tstates int
	mode    AddrModeFunc
}

var executors = [256]ExecInfo{}

func init() {
	executors[0xA9] = ExecInfo{LDA, 2, 2, ReadImmediate}
	executors[0xA5] = ExecInfo{LDA, 2, 3, ReadAbsoluteZeroPage}
	executors[0xB5] = ExecInfo{LDA, 2, 4, ReadZeroPageIdxX}
	executors[0xAD] = ExecInfo{LDA, 3, 4, ReadAbsolute}

	executors[0x85] = ExecInfo{STA, 2, 3, WriteAbsoluteZeroPage}

	executors[0x18] = ExecInfo{CLC, 1, 2, nil}
	executors[0x38] = ExecInfo{SEC, 1, 2, nil}

	executors[0xaa] = ExecInfo{TAX, 1, 2, nil}
	executors[0xa8] = ExecInfo{TAY, 1, 2, nil}
	executors[0xba] = ExecInfo{TSX, 1, 2, nil}
	executors[0x8a] = ExecInfo{TXA, 1, 2, nil}
	executors[0x9a] = ExecInfo{TXS, 1, 2, nil}
	executors[0x98] = ExecInfo{TYA, 1, 2, nil}

}

func setFlagsFromValue(ctx CPUContext, val uint8) uint8 {
	ctx.SetFlag(Flag_Z, val == 0)
	ctx.SetFlag(Flag_N, (val&0x80) == 0x80)
	return val
}

func LDA(ctx CPUContext, info *ExecInfo) int {
	val, exclock := info.mode(ctx, 0)
	ctx.SetRegA(setFlagsFromValue(ctx, val))
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates + exclock
}

func STA(ctx CPUContext, info *ExecInfo) int {
	_, exclock := info.mode(ctx, ctx.RegA())
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates + exclock
}

func CLC(ctx CPUContext, info *ExecInfo) int {
	ctx.SetFlag(Flag_C, false)
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

func SEC(ctx CPUContext, info *ExecInfo) int {
	ctx.SetFlag(Flag_C, true)
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

func TAX(ctx CPUContext, info *ExecInfo) int {
	ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegA()))
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

func TAY(ctx CPUContext, info *ExecInfo) int {
	ctx.SetRegY(setFlagsFromValue(ctx, ctx.RegA()))
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

func TSX(ctx CPUContext, info *ExecInfo) int {
	ctx.SetRegX(setFlagsFromValue(ctx, ctx.RegSP()))
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

func TXA(ctx CPUContext, info *ExecInfo) int {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.RegX()))
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

func TXS(ctx CPUContext, info *ExecInfo) int {
	ctx.SetRegSP(ctx.RegX())
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

func TYA(ctx CPUContext, info *ExecInfo) int {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.RegY()))
	ctx.SetRegPC(ctx.RegPC() + info.length)
	return info.tstates
}

/*
func exec_LDA_zeropage(ctx CPUContext, pc uint16) int {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(uint16(ctx.Peek(pc+1)))))
	ctx.SetRegPC(pc + 2)
	return 3
}

func exec_LDA_zeropageX(ctx CPUContext, pc uint16) int {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(uint16(ctx.Peek(pc+1)+ctx.RegX()))))
	ctx.SetRegPC(pc + 2)
	return 4
}

func exec_LDA_absolute(ctx CPUContext, pc uint16) int {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(ctx.PeekWord(pc+1))))
	ctx.SetRegPC(pc + 3)
	return 4
}

func exec_LDA_absoluteX(ctx CPUContext, pc uint16) int {
	addr := ctx.PeekWord(pc + 1)
	addrx := addr + uint16(ctx.RegX())
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(addrx)))
	ctx.SetRegPC(pc + 3)
	if (addr & 0xff00) == (addrx & 0xf00) {
		return 4
	} else {
		return 5
	}
}

func exec_LDA_absoluteY(ctx CPUContext, pc uint16) int {
	addr := ctx.PeekWord(pc + 1)
	addry := addr + uint16(ctx.RegY())
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(addry)))
	ctx.SetRegPC(pc + 3)
	if (addr & 0xff00) == (addry & 0xf00) {
		return 4
	} else {
		return 5
	}
}
*/
