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
	console.Successf("> SYS >>> pc:%04X i:%04X\n", v.pc, v.index)
	console.Successf("> SYS >>> ")
	for i := range v.registers {
		console.Successf("v%d:%02X ", i, v.registers[i])
	}
	console.Successf("\n")
	//console.Successf("> SYS >>> v0: %02X v1:%02X v0: %02X v1:%02X v0: %02X v1:%02X v0: %02X v1:%02X \n", v.registers[0], v.registers[1])
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
