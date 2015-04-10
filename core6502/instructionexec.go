package core6502

func setFlagsFromValue(ctx CPUContext, val uint8) uint8 {
	ctx.SetFlag(Flag_Z, val == 0)
	ctx.SetFlag(Flag_N, (val&0x80) == 0x80)
	return val
}

type instrExec func(ctx CPUContext, pc uint16) int

var executors = []instrExec{
	exec_LDA_immediate,
	exec_LDA_zeropage,
}

func exec_LDA_immediate(ctx CPUContext, pc uint16) int {
	ctx.SetRegA(setFlagsFromValue(ctx, ctx.Peek(pc+1)))
	ctx.SetRegPC(pc + 2)
	return 2
}

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
