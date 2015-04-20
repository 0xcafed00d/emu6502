package main

import (
	//"errors"
	"fmt"
	"github.com/simulatedsimian/emu6502/core6502"
	"reflect"
	"strconv"
	"strings"
)

/*
	var test interface{}

	test = ctx
	logDisp.WriteLine(fmt.Sprint(reflect.TypeOf(test)))

	t := reflect.TypeOf(testFunc)
	for n := 0; n < t.NumIn(); n++ {
		logDisp.WriteLine(fmt.Sprint(t.In(n)))
	}

	testFuncVal := reflect.ValueOf(testFunc)

	args := []reflect.Value{reflect.ValueOf("str"), reflect.ValueOf(uint8(22)), reflect.ValueOf(uint16(7625))}

	testFuncVal.Call(args)

	logDisp.WriteLine(fmt.Sprint(t))

	//		i, err := strconv.ParseInt(inp, 16, 16)

	//		logDisp.WriteLine(fmt.Sprint(i, err))
*/

type commandInfo struct {
	name, help string
	handler    reflect.Value
}

var commands = []commandInfo{
	{"sm", "Set Memory:   sm <address> <value>", reflect.ValueOf(setMemory)},
	{"sb", "Set Block:    smb <address> <count> <value>", reflect.ValueOf(setMemoryBlock)},
	{"sr", "Set Register: sr <reg> <value>", reflect.ValueOf(setReg)},
}

func processArgs(handler reflect.Value, ctx core6502.CPUContext, parts []string) ([]reflect.Value, error) {

	args := []reflect.Value{}
	args = append(args, reflect.ValueOf(ctx))

	handler_t := handler.Type()
	//panic(fmt.Sprint(handler_t.In(1), reflect.TypeOf(uint8(0))))
	for n := 0; n < handler_t.NumIn(); n++ {
		switch handler_t.In(n) {

		case reflect.TypeOf(uint8(0)):
			i, err := strconv.ParseUint(parts[0], 16, 8)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(uint8(i)))
			parts = parts[1:]

		case reflect.TypeOf(uint16(0)):
			i, err := strconv.ParseUint(parts[0], 16, 16)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(uint16(i)))
			parts = parts[1:]

		case reflect.TypeOf(""):
			args = append(args, reflect.ValueOf(parts[0]))
			parts = parts[1:]
		}
	}

	// todo verify arg count

	return args, nil
}

func DispatchCommand(ctx core6502.CPUContext, cmd string) (bool, error) {
	if cmd == "q" {
		return true, nil
	}

	parts := strings.Split(cmd, " ")
	if len(parts) > 0 {
		for n := 0; n < len(commands); n++ {
			if commands[n].name == parts[0] {
				args, err := processArgs(commands[n].handler, ctx, parts[1:])
				if err != nil {
					return false, err
				}
				commands[n].handler.Call(args)
				return false, nil
			}
		}
		return false, fmt.Errorf("Unknown Command: %s", parts[0])
	}
	return false, nil
}

func setMemory(ctx core6502.CPUContext, addr uint16, val uint8) error {
	ctx.Poke(addr, val)
	return nil
}

func setMemoryBlock(ctx core6502.CPUContext, addr uint16, count uint16, val uint8) error {
	for count != 0 {
		ctx.Poke(addr, val)
		addr++
		count--
	}
	return nil
}

func setReg(ctx core6502.CPUContext, reg string, val uint8) error {
	switch reg {
	case "a":
		ctx.SetRegA(val)
	case "x":
		ctx.SetRegX(val)
	case "y":
		ctx.SetRegY(val)
	default:
		return fmt.Errorf("Invalid Register: %s", reg)
	}
	return nil
}
