package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/emu6502/core6502"
)

func printAt(x, y int, s string, fg, bg termbox.Attribute) {
	for _, r := range s {
		termbox.SetCell(x, y, r, fg, bg)
		x++
	}
}

func printAtDef(x, y int, s string) {
	printAt(x, y, s, termbox.ColorDefault, termbox.ColorDefault)
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func printRegisters(x, y int, ctx core6502.CPUContext) {
	printAtDef(x, y, fmt.Sprintf("A: $%02x X: $%02x Y: $%02x", ctx.RegA(), ctx.RegX(), ctx.RegY()))
	printAtDef(x, y+1, fmt.Sprintf("PC: $%04x    SP: $%02x", ctx.RegPC(), ctx.RegSP()))
	printAtDef(x, y+2, fmt.Sprintf("FLAGS: N V B D I Z C"))
	printAtDef(x, y+3, fmt.Sprintf("       %x %x %x %x %x %x %x",
		btoi(ctx.Flag(core6502.Flag_N)),
		btoi(ctx.Flag(core6502.Flag_V)),
		btoi(ctx.Flag(core6502.Flag_B)),
		btoi(ctx.Flag(core6502.Flag_D)),
		btoi(ctx.Flag(core6502.Flag_I)),
		btoi(ctx.Flag(core6502.Flag_Z)),
		btoi(ctx.Flag(core6502.Flag_C))))
}

func printMemory(x, y int, addr uint16, ctx core6502.CPUContext) {

	for l := 0; l < 16; l++ {
		printAtDef(x, y+l, fmt.Sprintf("%04x:", addr))
		for n := 0; n < 16; n++ {
			printAtDef(x+6+n*4, y+l, fmt.Sprintf("%02x", ctx.Peek(addr)))
			addr++
		}
	}
}
