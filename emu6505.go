package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/emu6502/core6502"
)

type TextInputField struct {
	x, y      int
	inp       []rune
	cursorLoc int
}

func (t *TextInputField) HandleInput(k termbox.Key, r rune) {
	if r > ' ' {
		t.inp = append(t.inp, r)
	}

	if k == 32 {
		t.inp = append(t.inp, ' ')
	}

	if len(t.inp) > 0 && (k == termbox.KeyBackspace || k == termbox.KeyBackspace2) {
		t.inp = t.inp[:len(t.inp)-1]
	}

	printAtDef(t.x, t.y+1, fmt.Sprintf("%v, %v               ", k, r))

}

func (t *TextInputField) Draw() {
	printAtDef(t.x, t.y, string(t.inp)+" ")
}

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

	textInput := TextInputField{0, 10, nil, 0}

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
			textInput.HandleInput(ev.Key, ev.Ch)
			textInput.Draw()

			termbox.Flush()

		case termbox.EventResize:
			//			termx, termy = ev.Width, ev.Height
			termbox.Flush()
		}
	}
}
