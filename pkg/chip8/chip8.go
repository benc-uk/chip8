//
// CHIP-8 - CPU and virtual machine, the core of the emulation is done here
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/benc-uk/chip8/pkg/console"
	"github.com/benc-uk/chip8/pkg/font"
)

// Where fonts are loaded
const FontBase = 0x0050

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
	stack   []uint16

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
		v.memory[FontBase+i] = fontByte
	}

	return &v
}

// Run the VM processor with a channel for reporting errors
func (v *VM) Run(errors chan error, delay int) {
	for {
		err := v.Cycle()
		if err != nil {
			errors <- err

			// Halt the processor
			return
		}
		time.Sleep(time.Duration(delay) * time.Microsecond)
	}
}

// Cycle is the heart of the CHIP-8 emulator, running a single processor cycle
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

	// Debug VM system state, PC, index, registers, stack etc
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
				v.insCLS()
			}
			if o.nn == 0xEE {
				v.insRET()
			}
		}
	case 0x1:
		v.insJP(o.nnn)
	case 0x2:
		v.insCALL(o.nnn)
	case 0x3:
		v.insSEvb(o.x, o.nn)
	case 0x4:
		v.insSNEvb(o.x, o.nn)
	case 0x6:
		v.insLDvb(o.x, o.nn)
	case 0x7:
		v.insADDvb(o.x, o.nn)
	case 0xA:
		v.insLDi(o.nnn)
	case 0xD:
		v.insDRW(o.x, o.y, o.n)
	case 0xF:
		{
			switch o.nn {
			case 0x29:
				v.insLDf(o.x)
			}
		}
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
