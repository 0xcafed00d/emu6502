package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/emu6502/core6502"
	"reflect"
	"strconv"
)

func testFunc(ctx core6502.CPUContext, cmd string, x uint8, y uint16, s string) {

}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	var ctx core6502.BasicCPUContext
	core6502.HardResetCPU(&ctx, 0x400)

	doQuit := false

	dl := DisplayList{}

	regDisp := RegisterDisplay{1, 1, &ctx}
	memDisp := MemoryDisplay{30, 1, 0, &ctx}
	stkDisp := StackDisplay{24, 1, &ctx}
	logDisp := ScrollingTextOutput{1, 20, 80, 10, nil}

	cmdInput := MakeTextInputField(10, 18, func(inp string) {
		if inp == "q" {
			doQuit = true
		}

		t := reflect.TypeOf(testFunc)
		for n := 0; n < t.NumIn(); n++ {
			logDisp.WriteLine(fmt.Sprint(t.In(n)))
		}

		logDisp.WriteLine(fmt.Sprint(t))

		i, err := strconv.ParseInt(inp, 16, 16)

		logDisp.WriteLine(fmt.Sprint(i, err))
	})

	dl.AddElement(cmdInput)
	dl.AddElement(&regDisp)
	dl.AddElement(&memDisp)
	dl.AddElement(&stkDisp)
	dl.AddElement(&logDisp)
	dl.AddElement(&StaticText{1, 18, "Command:"})
	dl.AddElement(&StaticText{1, 0, "Registers:"})
	dl.AddElement(&StaticText{30, 0, "Memory:"})
	dl.AddElement(&StaticText{24, 0, "TOS:"})

	cmdInput.GiveFocus()

	dl.Draw()
	termbox.Flush()

	for !doQuit {
		ev := termbox.PollEvent()

		if ev.Type == termbox.EventKey {
			dl.HandleInput(ev.Key, ev.Ch)
			dl.Draw()
			termbox.Flush()
		}

		if ev.Type == termbox.EventResize {
			termbox.Flush()
		}
	}
}
