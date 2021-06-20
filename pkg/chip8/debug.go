//
// CHIP-8 - Debugging to the console
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"github.com/benc-uk/chip8/pkg/console"
)

var enableDebug = false

func (o Opcode) dump() {
	if !enableDebug {
		return
	}
	console.Infof("> OPC >>> kind: %X, x: %X, y: %X, n:%X, nn:%X, nnn:%X\n", o.kind, o.x, o.y, o.n, o.nn, o.nnn)
}

func (v VM) dump() {
	if !enableDebug {
		return
	}
	console.Successf("> SYS >>> PC:%04X  I:%04X  DT:%02X  ST:%02X \n", v.pc, v.index, v.delayTimer, v.soundTimer)

	// Dump registers
	console.Successf("> REG >>> ")
	for i := range v.registers {
		console.Successf("v%X:%02X  ", i, v.registers[i])
	}
	console.Successf("\n")

	// Dump stack
	if len(v.stack) > 0 {
		console.Successf("> SYS >>> ")
		for i := range v.stack {
			console.Successf("stack[%d]:%04X ", i, v.stack[i])
		}
		console.Successf("\n")
	}
}

func (v *VM) DumpMemory(start int, end int) {
	if !enableDebug {
		return
	}
	for i := start; i < end; i++ {
		console.Warningf("%04X: %02X\n", i, v.memory[i])
	}
}

func debug(s string) {
	if !enableDebug {
		return
	}
	console.Debug(s)
}

func debugf(f string, a ...interface{}) {
	if !enableDebug {
		return
	}
	console.Debugf(f, a...)
}
