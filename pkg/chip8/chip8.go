package chip8

import (
	"fmt"

	"github.com/benc-uk/chip8/pkg/font"
)

const fontBase = 0x50

type Opcode uint16

type System struct {
	memory    [4096]byte
	registers [16]byte
	pc        uint16
	index     uint16
}

func NewSystem() System {
	s := System{}
	s.ClearMemory()
	s.ClearRegisters()
	s.index = 0
	s.pc = 0x200

	// Load font into memory
	for i, fontByte := range font.GetFont() {
		s.memory[fontBase+i] = fontByte
	}

	return s
}

func (s System) Run() {
	for i := range s.memory {
		fmt.Println(i, s.memory[i])
		if i > 0x80 {
			break
		}
	}
}

func (s System) ClearMemory() {
	for i := range s.memory {
		s.memory[i] = 0
	}
}

func (s System) ClearRegisters() {
	for i := range s.registers {
		s.registers[i] = 0
	}
}

func (s System) FetchOpCode() Opcode {
	return 8
}
