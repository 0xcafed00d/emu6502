package core6502

import (
	//	"fmt"
	"github.com/simulatedsimian/assert"
)

type asmInfo map[AddressMode]uint8

var asmData map[string]asmInfo

func initx() {
	for n := 0; n < len(InstructionData); n++ {
		info := &InstructionData[n]
		name := assert.GetShortFuncName(info.execMaker)

		if _, ok := asmData[name]; !ok {
			asmData[name] = asmInfo{}
		}
		asmData[name][info.mode] = info.opcode
	}
}

func Assemble(ctx CPUContext, addr uint16, s string) (uint16, error) {
	return 0, nil
}
