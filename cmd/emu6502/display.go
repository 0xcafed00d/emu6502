package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/emu6502/core6502"
)

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

type RegisterDisplay struct {
	x, y int
	ctx  core6502.CPUContext
}

func (rd *RegisterDisplay) Draw() {
	printAtDef(rd.x, rd.y, fmt.Sprintf("PC: $%04x    SP: $%02x", rd.ctx.RegPC(), rd.ctx.RegSP()))
	printAtDef(rd.x, rd.y+1, fmt.Sprintf("A: $%02x X: $%02x Y: $%02x", rd.ctx.RegA(), rd.ctx.RegX(), rd.ctx.RegY()))
	printAtDef(rd.x, rd.y+2, fmt.Sprintf("FLAGS: N V B D I Z C"))
	printAtDef(rd.x, rd.y+3, fmt.Sprintf("       %x %x %x %x %x %x %x",
		btoi(rd.ctx.Flag(core6502.Flag_N)),
		btoi(rd.ctx.Flag(core6502.Flag_V)),
		btoi(rd.ctx.Flag(core6502.Flag_B)),
		btoi(rd.ctx.Flag(core6502.Flag_D)),
		btoi(rd.ctx.Flag(core6502.Flag_I)),
		btoi(rd.ctx.Flag(core6502.Flag_Z)),
		btoi(rd.ctx.Flag(core6502.Flag_C))))
}

func (rd *RegisterDisplay) GiveFocus() bool {
	termbox.SetCursor(rd.x, rd.y)
	return true
}

func (rd *RegisterDisplay) HandleInput(k termbox.Key, r rune) {
}

type MemoryDisplay struct {
	x, y int
	addr uint16
	ctx  core6502.CPUContext
}

func (md *MemoryDisplay) Draw() {
	addr := md.addr
	for l := 0; l < 16; l++ {
		printAtDef(md.x, md.y+l, fmt.Sprintf("$%04x:", addr))
		for n := 0; n < 16; n++ {
			printAtDef(md.x+6+n*3, md.y+l, fmt.Sprintf("%02x", md.ctx.Peek(addr)))

			c := rune(md.ctx.Peek(addr))
			if c < ' ' || c > 127 {
				c = '.'
			}
			termbox.SetCell(md.x+55+n, md.y+l, c, termbox.ColorDefault, termbox.ColorDefault)
			addr++
		}
	}
}

func (md *MemoryDisplay) GiveFocus() bool {
	termbox.SetCursor(md.x, md.y)
	return true
}

func (md *MemoryDisplay) HandleInput(k termbox.Key, r rune) {
	if k == termbox.KeyPgdn {
		md.addr += 64
	}

	if k == termbox.KeyPgup {
		md.addr -= 64
	}
}

type StackDisplay struct {
	x, y int
	ctx  core6502.CPUContext
}

func (sd *StackDisplay) Draw() {
	sp := sd.ctx.RegSP() + 1

	for l := 0; l < 16; l++ {
		printAtDef(sd.x, sd.y+l, fmt.Sprintf("$%02x", sd.ctx.Peek(uint16(sp)+0x100)))
		sp++
	}
}

func (sd *StackDisplay) GiveFocus() bool {
	return false
}

func (sd *StackDisplay) HandleInput(k termbox.Key, r rune) {
}

type DisasmDisplay struct {
	x, y  int
	lines int
	ctx   core6502.CPUContext
}

func (dd *DisasmDisplay) Draw() {
	pc := dd.ctx.RegPC()

	for l := 0; l < dd.lines; l++ {
		instr, len, _ := core6502.Disassemble(dd.ctx, pc)
		printAtDef(dd.x, dd.y+l, fmt.Sprintf("$%04x %s", pc, instr))
		pc += len
	}
}

func (dd *DisasmDisplay) GiveFocus() bool {
	return false
}

func (dd *DisasmDisplay) HandleInput(k termbox.Key, r rune) {
}
