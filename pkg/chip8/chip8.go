//
// CHIP-8 - CPU and virtual machine, the core of the emulation is done here
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"encoding/binary"
	"errors"

	"github.com/benc-uk/chip8/pkg/console"
	"github.com/benc-uk/chip8/pkg/font"
)

// Where fonts are loaded
const fontBase = 0x50

// ProgBase is where programs should be loaded into memory
const ProgBase = 0x200

// Normal CHIP-8 systems have 4KB of memory
const memSize = 0x1000 // 4096 bytes

// DisplayHeight standard CHIP-8 display height
const DisplayHeight = 32

// DisplayWidth standard CHIP-8 display width
const DisplayWidth = 64

// Opcode holds a decoded opcode see: docs/opcode.md
// IMPORTANT: the decoder will set all values BUT their use is opcode dependant
type Opcode struct {
	kind uint8  // nibble
	x    uint8  // nibble
	y    uint8  // nibble
	n    uint8  // nibble
	nn   uint8  // byte
	nnn  uint16 // 12 bits
}

// VM is a CHIP-8 vritual machine
type VM struct {
	// The implementation of a CHIP-8 system
	memory     [memSize]byte
	registers  [16]byte
	pc         uint16
	index      uint16
	timerDelay byte
	timerSound byte
	// TODO: Is using booleans dumb?
	display [DisplayWidth][DisplayHeight]bool

	// Supporting fields for emulation. not part of the system architecture
	running bool
}

func NewVM(debug bool) *VM {
	v := VM{}
	console.Info("CHIP-8 system created...")
	v.reset()
	enableDebug = debug

	// Load font into lower memory
	for i, fontByte := range font.GetFont() {
		v.memory[fontBase+i] = fontByte
	}

	return &v
}

// Cycle is the main emulator function, running a processor cycle
func (v *VM) Cycle() error {
	if !v.running {
		return nil
	}

	debug("______________________________________________________")

	// First get the 16 bit opcode at the current PC
	opcodeRaw, err := v.fetch()
	if err != nil {
		return err
	}

	// Decode the raw opcode into an parsed Opcode
	opcode := decode(opcodeRaw)
	opcode.dump()

	// Execute parses the opcode and excutes instructions
	v.execute(opcode)

	// Debug system state
	v.dump()

	return nil
}

func (v *VM) fetch() (uint16, error) {
	if v.pc >= memSize {
		return 0, errors.New("PC went outside of memory bounds")
	}

	op := binary.BigEndian.Uint16(v.memory[v.pc : v.pc+2])
	debugf("> FET >>> %04X (%02X)\n", v.memory[v.pc:v.pc+2], op)

	// VERY IMPORTANT! Move the PC to the next address in memory
	v.pc = v.pc + 2

	return op, nil
}

func decode(rawOpcode uint16) Opcode {
	return Opcode{
		kind: uint8(rawOpcode & 0xF000 >> 12),
		x:    uint8(rawOpcode & 0x0F00 >> 8),
		y:    uint8(rawOpcode & 0x00F0 >> 4),
		n:    uint8(rawOpcode & 0x000F >> 0),
		nn:   uint8(rawOpcode & 0x00FF),
		nnn:  rawOpcode & 0x0FFF,
	}
}

func (v *VM) execute(o Opcode) {
	switch o.kind {
	case 0x0:
		{
			if o.nn == 0xE0 {
				v.CLS()
			}
			if o.nn == 0xEE {
				v.CLS()
			}
		}
	case 0x1:
		v.JP(o.nnn)
	case 0x6:
		v.LDVB(o.x, o.nn)
	case 0x7:
		v.ADDVB(o.x, o.nn)
	case 0xA:
		v.LDI(o.nnn)
	case 0xD:
		v.DRW(o.x, o.y, o.n)
	}
}

func (v *VM) reset() {
	console.Info("System was reset")
	v.clearMemory()
	v.clearRegisters()
	v.index = 0
	v.pc = ProgBase
	v.running = true
}

func (v *VM) clearMemory() {
	for i := ProgBase; i < len(v.memory); i++ {
		v.memory[i] = 0
	}
}

func (v *VM) clearRegisters() {
	for i := range v.registers {
		v.registers[i] = 0
	}
}
