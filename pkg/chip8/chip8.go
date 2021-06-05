package chip8

import (
	"encoding/binary"
	"fmt"

	"github.com/benc-uk/chip8/pkg/font"
)

const fontBase = 0x50

type Opcode struct {
	kind uint8  // nibble
	x    uint8  // nibble
	y    uint8  // nibble
	n    uint16 // n, nn, nnn
}

type VM struct {
	memory     [4096]byte
	registers  [16]byte
	pc         uint16
	index      uint16
	timerDelay byte
	timerSound byte
}

func NewVM() VM {
	v := VM{}
	v.clearMemory()
	v.clearRegisters()
	v.index = 0
	v.pc = 0x200

	// Load font into memory
	for i, fontByte := range font.GetFont() {
		v.memory[fontBase+i] = fontByte
	}

	return v
}

// The main emulator loop
func (v VM) Run() {
	for {
		opcodeRaw := v.fetch()
		opcode := decode(opcodeRaw)
		//v.execute(opcode)
		fmt.Println(opcode.kind)
	}
}

func (v VM) clearMemory() {
	for i := range v.memory {
		v.memory[i] = 0
	}
}

func (v VM) clearRegisters() {
	for i := range v.registers {
		v.registers[i] = 0
	}
}

func (v VM) fetch() uint16 {
	op := binary.BigEndian.Uint16(v.memory[v.pc : v.pc+1])
	v.pc = v.pc + 2
	return op
}

func decode(raw uint16) Opcode {
	k := uint8(raw & 0xF000)

	return Opcode{
		kind: k,
	}
}
