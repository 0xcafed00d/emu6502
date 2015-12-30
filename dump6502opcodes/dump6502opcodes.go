package main

import (
	"fmt"
	"github.com/simulatedsimian/emu6502/core6502"
)

func main() {
	var ctx core6502.BasicCPUContext
	core6502.HardResetCPU(&ctx, 0x400)

	for opcode := 0; opcode < 256; opcode++ {
		ctx.Poke(0x400, uint8(opcode))

		dis, _ := core6502.Disassemble(&ctx, 0x400)
		fmt.Printf("$%02x %s\n", opcode, dis)
	}
}
