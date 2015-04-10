package core6502

func setFlagsFromValue(ctx CPUContext, val uint8) uint {
	ctx.SetFlag(Flag_Z, val == 0)
	ctx.SetFlag(Flag_N, (val&0x80) == 0x80)
	return val
}

type instrExec func(ctx CPUContext, opcode uint8, addr uint16)

var executors = []instrExec{
	exec_LDA_imediate,
	exec_LDA_zeropage,
}

func exec_LDA_imediate(ctx CPUContext, opcode uint8, addr uint16) {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(addr+1)))
	ctx.SetRegPC(addr + 2)
}

func exec_LDA_zeropage(ctx CPUContext, opcode uint8, addr uint16) {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(ctx.Peek(addr+1))))
	ctx.SetRegPC(addr + 2)
}
