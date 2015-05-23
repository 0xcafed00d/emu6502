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
	execMaker ExecFuncMakerFunc
	tstates   int
	mode      AddressMode
}

var InstructionData = []InstructionInfo{
	{0xa9, LDA, 2, AddrMode_Immediate},
	{0xa5, LDA, 3, AddrMode_AbsoluteZeroPage},
	{0xb5, LDA, 4, AddrMode_ZeroPageIdxX},
	{0xa1, LDA, 6, AddrMode_PostIndexIndirect},
	{0xb1, LDA, 5, AddrMode_PreIndexIndirect},
	{0xad, LDA, 4, AddrMode_Absolute},
	{0xbd, LDA, 4, AddrMode_AbsoluteIndexedX},
	{0xb9, LDA, 4, AddrMode_AbsoluteIndexedY},
	{0xa2, LDX, 2, AddrMode_Immediate},
	{0xa6, LDX, 3, AddrMode_AbsoluteZeroPage},
	{0xb6, LDX, 4, AddrMode_ZeroPageIdxY},
	{0xae, LDX, 4, AddrMode_Absolute},
	{0xbe, LDX, 4, AddrMode_AbsoluteIndexedY},
	{0xa0, LDY, 2, AddrMode_Immediate},
	{0xa4, LDY, 3, AddrMode_AbsoluteZeroPage},
	{0xb4, LDY, 4, AddrMode_ZeroPageIdxX},
	{0xac, LDY, 4, AddrMode_Absolute},
	{0xbc, LDY, 4, AddrMode_AbsoluteIndexedX},
	{0x85, STA, 3, AddrMode_AbsoluteZeroPage},
	{0x95, STA, 4, AddrMode_ZeroPageIdxX},
	{0x81, STA, 6, AddrMode_PostIndexIndirect},
	{0x91, STA, 6, AddrMode_PreIndexIndirect},
	{0x8d, STA, 4, AddrMode_Absolute},
	{0x9d, STA, 5, AddrMode_AbsoluteIndexedX},
	{0x99, STA, 5, AddrMode_AbsoluteIndexedY},
	{0x86, STX, 3, AddrMode_AbsoluteZeroPage},
	{0x96, STX, 4, AddrMode_ZeroPageIdxY},
	{0x8e, STX, 4, AddrMode_Absolute},
	{0x84, STY, 3, AddrMode_AbsoluteZeroPage},
	{0x94, STY, 4, AddrMode_ZeroPageIdxX},
	{0x8c, STY, 4, AddrMode_Absolute},
	{0xe6, INC, 4, AddrMode_AbsoluteZeroPage},
	{0xf6, INC, 4, AddrMode_ZeroPageIdxX},
	{0xee, INC, 4, AddrMode_Absolute},
	{0xfe, INC, 4, AddrMode_AbsoluteIndexedX},
	{0xe8, INX, 4, AddrMode_Implicit},
	{0xc8, INY, 4, AddrMode_Implicit},
	{0xc6, DEC, 4, AddrMode_AbsoluteZeroPage},
	{0xd6, DEC, 4, AddrMode_ZeroPageIdxX},
	{0xce, DEC, 4, AddrMode_Absolute},
	{0xde, DEC, 4, AddrMode_AbsoluteIndexedX},
	{0xca, DEX, 4, AddrMode_Implicit},
	{0x88, DEY, 4, AddrMode_Implicit},
	{0xaa, TAX, 2, AddrMode_Implicit},
	{0xa8, TAY, 2, AddrMode_Implicit},
	{0xba, TSX, 2, AddrMode_Implicit},
	{0x8a, TXA, 2, AddrMode_Implicit},
	{0x9a, TXS, 2, AddrMode_Implicit},
	{0x98, TYA, 2, AddrMode_Implicit},
	{0x18, CLC, 2, AddrMode_Implicit},
	{0x38, SEC, 2, AddrMode_Implicit},
	{0xD8, CLD, 2, AddrMode_Implicit},
	{0xF8, SED, 2, AddrMode_Implicit},
	{0x58, CLI, 2, AddrMode_Implicit},
	{0x78, SEI, 2, AddrMode_Implicit},
	{0xB8, CLV, 2, AddrMode_Implicit},
	{0xea, NOP, 1, AddrMode_Implicit},
	{0x00, BRK, 7, AddrMode_Implicit},
	{0x48, PHA, 3, AddrMode_Implicit},
	{0x68, PLA, 4, AddrMode_Implicit},
	{0x08, PHP, 3, AddrMode_Implicit},
	{0x28, PLP, 4, AddrMode_Implicit},

	{0x10, BPL, 2, AddrMode_Relative},
	{0x30, BMI, 2, AddrMode_Relative},
	{0x50, BVC, 2, AddrMode_Relative},
	{0x70, BVS, 2, AddrMode_Relative},
	{0x90, BCC, 2, AddrMode_Relative},
	{0xB0, BCS, 2, AddrMode_Relative},
	{0xD0, BNE, 2, AddrMode_Relative},
	{0xF0, BEQ, 2, AddrMode_Relative},
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

func CLD(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_D, false)
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func SED(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_D, true)
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func CLI(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_I, false)
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func SEI(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_I, true)
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func CLV(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlag(Flag_V, false)
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

func NOP(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func BRK(info *InstructionInfo) InstructionExecFunc {

	return func(ctx CPUContext) int {
		Push16(ctx, ctx.RegPC()+1)
		ctx.SetFlag(Flag_B, true)
		Push8(ctx, ctx.Flags())
		ctx.SetFlag(Flag_I, true)
		ctx.SetRegPC(ctx.PeekWord(Vector_IRQ))
		return info.tstates
	}
}

func PHA(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		Push8(ctx, ctx.RegA())
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func PLA(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetRegA(setFlagsFromValue(ctx, Pop8(ctx)))
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func PHP(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		Push8(ctx, ctx.Flags())
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func PLP(info *InstructionInfo) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		ctx.SetFlags(Pop8(ctx))
		ctx.SetFlag(Flag_B, false)
		ctx.SetRegPC(ctx.RegPC() + length)
		return info.tstates
	}
}

func makeBranchExecFunc(info *InstructionInfo, testFunc func(CPUContext) bool) InstructionExecFunc {
	length := InstructionBytes(info.mode)

	return func(ctx CPUContext) int {
		if testFunc(ctx) {
			newPC, exclock := CalcPCRelativeAddr(ctx)
			ctx.SetRegPC(newPC)
			return info.tstates + exclock
		} else {
			ctx.SetRegPC(ctx.RegPC() + length)
			return info.tstates
		}
	}
}

func BPL(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return !ctx.Flag(Flag_N)
	})
}

func BMI(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return ctx.Flag(Flag_N)
	})
}

func BVC(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return !ctx.Flag(Flag_V)
	})
}

func BVS(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return ctx.Flag(Flag_V)
	})
}

func BCC(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return !ctx.Flag(Flag_C)
	})
}

func BCS(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return ctx.Flag(Flag_C)
	})
}

func BNE(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return !ctx.Flag(Flag_Z)
	})
}

func BEQ(info *InstructionInfo) InstructionExecFunc {
	return makeBranchExecFunc(info, func(ctx CPUContext) bool {
		return ctx.Flag(Flag_Z)
	})
}
