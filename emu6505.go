package main

import (
	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/emu6502/core6502"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.HideCursor()

	var ctx core6502.BasicCPUContext
	core6502.HardResetCPU(&ctx, 0x400)
	var addr uint16

	//	termx, termy := termbox.Size()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyPgdn:
				addr += 32
			case termbox.KeyPgup:
				addr -= 32
			}

			printRegisters(1, 1, &ctx)
			printMemory(25, 1, addr, &ctx)

			termbox.Flush()

		case termbox.EventResize:
			//			termx, termy = ev.Width, ev.Height
			termbox.Flush()
		}
	}

}
