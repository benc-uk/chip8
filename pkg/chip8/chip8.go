//
// CHIP-8 - CPU and virtual machine, the core of the emulation is done here
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/benc-uk/chip8/pkg/console"
	"github.com/benc-uk/chip8/pkg/font"
)

// ProgBase is where programs should be loaded into memory
const ProgBase = 0x200

// FontBase address where fonts are loaded
const FontBase = 0x0050

// FontLargeBase is where the larger 10 byte fonts are stored
const FontLargeBase = 0x0050 + 0x40 // 0x40 bytes is the size of the low res font

// Normal CHIP-8 systems have 4KB of memory
const memSize = 0x1000 // 4096 bytes

// DisplayHeight is Super CHIP-8 display height, note this is backwards compatible
const DisplayHeight = 64

// DisplayWidth is Super CHIP-8 display width, note this is backwards compatible
const DisplayWidth = 128

// Used for the timer loop to pause 1/60 second
const sixtyHzMicroSecs = 16700

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
	delayTimer byte
	soundTimer byte
	display    [DisplayWidth][DisplayHeight]uint16
	stack      []uint16

	// Super CHIP-8 extensions
	HighRes bool

	// Supporting fields for emulation. not part of the system architecture
	debug          bool
	DisplayUpdated bool
	// Flag for quirks e.g instructions F
	modernMode bool
	// Keys that are currently pressed, values are 0x0 ~ 0xF
	Keys []uint8
}

func NewVM(modernMode bool) *VM {
	v := VM{}
	console.Info("CHIP-8 system created...")
	v.Reset()

	// Load fonts into lower memory
	for i, fontByte := range font.GetFont() {
		v.memory[FontBase+i] = fontByte
	}
	for i, fontByte := range font.GetLargeFont() {
		v.memory[FontLargeBase+i] = fontByte
	}

	// Default to modern / quirks mode
	v.modernMode = modernMode

	// Start the timer loops for the VM
	go v.TimerLoop()

	return &v
}

// Cycle is the heart of the CHIP-8 emulator, running a single processor cycle
func (v *VM) Cycle() error {
	v.debugLogf("============== PC: %02X ==================\n", v.pc)

	// First get the 16 bit opcode at the current PC
	opcodeRaw, err := v.fetch()
	if err != nil {
		return err
	}

	// Decode the raw opcode into an parsed Opcode
	opcode := decode(opcodeRaw)
	//opcode.dump()

	// Execute parses the opcode and excutes instructions
	err = v.execute(opcode)
	if err != nil {
		return err
	}

	// Debug VM system state, PC, index, registers, stack etc
	if v.debug {
		v.Dump()
	}

	return nil
}

func (v *VM) TimerLoop() {
	for {
		if v.delayTimer > 0 {
			v.delayTimer--
		}
		if v.soundTimer > 0 {
			v.soundTimer--
		}
		// Wait for 60hz
		time.Sleep(time.Duration(sixtyHzMicroSecs) * time.Microsecond)
	}
}

func (v *VM) fetch() (uint16, error) {
	if v.pc >= memSize {
		err := SystemError{
			reason: "PC went outside of memory bounds",
			code:   errorCodeAddress,
		}
		return 0, err
	}

	op := binary.BigEndian.Uint16(v.memory[v.pc : v.pc+2])
	v.debugLogf("> FET >>> %04X\n", v.memory[v.pc:v.pc+2])

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

func (v *VM) execute(o Opcode) error {
	switch o.kind {
	case 0x0:
		{
			if o.y == 0xC {
				v.insSCRD(o.n) // Super CHIP-8
				return nil
			}
			switch o.nn {
			case 0xE0:
				v.insCLS()
			case 0xEE:
				v.insRET()
			case 0xEF:
				v.insLOW() // Super CHIP-8
			case 0xFB:
				v.insSCRR() // Super CHIP-8
			case 0xFC:
				v.insSCRL() // Super CHIP-8
			case 0xFD:
				v.insEXIT() // Super CHIP-8
			case 0xFF:
				v.insHIGH() // Super CHIP-8
			default:
				return SystemError{
					reason: fmt.Sprintf("Invalid opcode %+v", o),
					code:   errorBadOpcode,
				}
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
	case 0x5:
		v.insSExy(o.x, o.y)
	case 0x6:
		v.insLDvb(o.x, o.nn)
	case 0x7:
		v.insADDvb(o.x, o.nn)
	case 0x8:
		{
			switch o.n {
			case 0:
				v.insLDxy(o.x, o.y)
			case 1:
				v.insORxy(o.x, o.y)
			case 2:
				v.insANDxy(o.x, o.y)
			case 3:
				v.insXORxy(o.x, o.y)
			case 4:
				v.insADDxy(o.x, o.y)
			case 5:
				v.insSUBxy(o.x, o.y)
			case 6:
				v.insSHRxy(o.x, o.y)
			case 7:
				v.insSUBNxy(o.x, o.y)
			case 0xE:
				v.insSHLxy(o.x, o.y)
			default:
				return SystemError{
					reason: fmt.Sprintf("Invalid opcode %+v", o),
					code:   errorBadOpcode,
				}
			}
		}
	case 0x9:
		v.insSNExy(o.x, o.y)
	case 0xA:
		v.insLDI(o.nnn)
	case 0xB:
		v.insJPV0(o.nnn)
	case 0xC:
		v.insRNDvb(o.x, o.nn)
	case 0xD:
		v.insDRW(o.x, o.y, o.n)
	case 0xE:
		{
			switch o.nn {
			case 0x9E:
				v.insSKPx(o.x)
			case 0xA1:
				v.insSKNPx(o.x)
			}
		}
	case 0xF:
		{
			switch o.nn {
			case 0x07:
				v.insLDxDT(o.x)
			case 0x0A:
				v.insLDxK(o.x)
			case 0x15:
				v.insLDDTx(o.x)
			case 0x18:
				v.insLDSTx(o.x)
			case 0x1E:
				v.insADDIx(o.x)
			case 0x29:
				v.insLDf(o.x)
			case 0x30:
				v.insLDSf(o.x)
			case 0x33:
				v.insLDBx(o.x)
			case 0x55:
				v.insLDIx(o.x)
			case 0x65:
				v.insLDxI(o.x)
			case 0x75:
				console.Error("F075 unsupported")
				return nil
			case 0x85:
				console.Error("F085 unsupported")
				return nil
			default:
				return SystemError{
					reason: fmt.Sprintf("Invalid opcode %+v", o),
					code:   errorBadOpcode,
				}
			}
		}
	default:
		return SystemError{
			reason: fmt.Sprintf("Invalid opcode %+v", o),
			code:   errorBadOpcode,
		}
	}

	return nil
}

func (v *VM) Reset() {
	console.Info("System was reset")
	v.clearMemory()
	v.clearRegisters()
	v.insCLS()
	v.index = 0
	v.pc = ProgBase
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

func (v *VM) LoadProgram(pgm []byte) error {
	// Reset the machine before writing program data to memory
	v.Reset()

	if len(pgm)+ProgBase > memSize {
		return SystemError{
			reason: "Out of memory!",
			code:   errorOutOfMemory,
		}
	}
	for i := range pgm {
		v.memory[ProgBase+i] = pgm[i]
	}

	console.Successf("Loaded %d bytes into memory OK\n", len(pgm))
	return nil
}

func (v *VM) DisplayValueAt(x int, y int) uint16 {
	return v.display[x][y]
}

func (v *VM) GetFlag() uint8 {
	return v.registers[0xF]
}

func (v *VM) SetFlag(val uint8) {
	v.registers[0xF] = val
}

func (v *VM) GetSoundTimer() uint8 {
	return v.soundTimer
}

func (v *VM) SetDebug(d bool) {
	v.debug = d
}

func (v *VM) IsDebugging() bool {
	return v.debug
}
