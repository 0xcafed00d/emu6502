package main

import (
	"fmt"
	"github.com/simulatedsimian/emu6502/core6502"
)

func main() {
	var ctx core6502.BasicCPUContext
	core6502.HardResetCPU(&ctx, 0x400)
	ctx.PokeWord(0x401, 0xeeff)

	for opcode := 0; opcode < 256; opcode++ {
		ctx.Poke(0x400, uint8(opcode))

		dis, len := core6502.Disassemble(&ctx, 0x400)
		if len == 1 {
			fmt.Printf("$%02x:              %s\n", opcode, dis)
		} else if len == 2 {
			fmt.Printf("$%02x $%02x:          %s\n", opcode, ctx.Peek(0x401), dis)
		} else if len == 3 {
			fmt.Printf("$%02x $%02x $%02x:      %s\n", opcode, ctx.Peek(0x401), ctx.Peek(0x402), dis)
		}
	}
}
