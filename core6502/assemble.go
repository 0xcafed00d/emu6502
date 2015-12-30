package core6502

import (
	"fmt"
	"github.com/simulatedsimian/assert"
)

type asmInfo map[AddressMode]uint8

var asmData map[string]asmInfo = make(map[string]asmInfo)

func init() {
	for n := 0; n < len(InstructionData); n++ {
		info := &InstructionData[n]
		name := assert.GetShortFuncName(info.execMaker)

		if _, ok := asmData[name]; !ok {
			asmData[name] = asmInfo{}
		}
		asmData[name][info.mode] = info.opcode
	}
}

func DetermineAddresMode(parts []string) AddressMode {
	return AddrMode_Absolute
}

func Assemble(ctx CPUContext, addr uint16, s string) (nextAddr uint16, err error) {
	parts := Split(s, " ,")
	addrMode := DetermineAddresMode(parts)

	info, ok := asmData[parts[0]]
	if !ok {
		err = fmt.Errorf("'%s' not valid instruction", parts[0])
		return
	}

	addrMode = addrMode
	info = info

	return
}
