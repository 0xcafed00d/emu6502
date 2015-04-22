package main

import (
	"fmt"
	"github.com/simulatedsimian/emu6502/core6502"
	"reflect"
	"strconv"
	"strings"
)

type commandInfo struct {
	help    string
	handler reflect.Value
}

var commands = map[string]commandInfo{
	"sm": {"Set Memory:   sm <address> <value>", reflect.ValueOf(setMemory)},
	"sb": {"Set Block:    sb <address> <count> <value>", reflect.ValueOf(setMemoryBlock)},
	"sr": {"Set Register: sr <reg> <value>", reflect.ValueOf(setReg)},
	"ps": {"Push Stack:   ps <value>", reflect.ValueOf(push)},
}

func processArgs(cmd commandInfo, ctx core6502.CPUContext, parts []string) ([]reflect.Value, error) {

	args := []reflect.Value{reflect.ValueOf(ctx)}

	for n := 0; n < cmd.handler.Type().NumIn(); n++ {
		if len(parts) == 0 {
			return nil, fmt.Errorf("Not enough Args: %s", cmd.help)
		}

		switch cmd.handler.Type().In(n).Kind() {

		case reflect.Uint8:
			i, err := strconv.ParseUint(parts[0], 16, 8)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(uint8(i)))
			parts = parts[1:]

		case reflect.Uint16:
			i, err := strconv.ParseUint(parts[0], 16, 16)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(uint16(i)))
			parts = parts[1:]

		case reflect.String:
			args = append(args, reflect.ValueOf(parts[0]))
			parts = parts[1:]
		}
	}

	if len(parts) > 0 {
		return nil, fmt.Errorf("Too Many Args: %s", cmd.help)
	}

	return args, nil
}

func DispatchCommand(ctx core6502.CPUContext, cmd string) (bool, error) {
	if cmd == "q" {
		return true, nil
	}

	parts := strings.Split(cmd, " ")
	if len(parts) > 0 && parts[0] != "" {
		if cmd, ok := commands[parts[0]]; ok {
			if args, err := processArgs(cmd, ctx, parts[1:]); err == nil {
				cmd.handler.Call(args)
			} else {
				return false, err
			}
		} else {
			return false, fmt.Errorf("Unknown Command: %s", parts[0])
		}
	}
	return false, nil
}

func setMemory(ctx core6502.CPUContext, addr uint16, val uint8) error {
	ctx.Poke(addr, val)
	return nil
}

func push(ctx core6502.CPUContext, val uint8) error {
	sp := ctx.RegSP()
	ctx.Poke(0x100+uint16(sp), val)
	ctx.SetRegSP(sp - 1)
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
