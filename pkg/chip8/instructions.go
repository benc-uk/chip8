//
// CHIP-8 - Implementation of opcodes / instructions here
// Ben C, June 2021
// Notes:

// See http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3.1

package chip8

import "github.com/benc-uk/chip8/pkg/console"

//
// Zero params
//

// CLS - clear screen - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#00E0
func (v *VM) CLS() {
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			v.display[x][y] = false
		}
	}
}

func (v *VM) RET() {
	console.Error("***** RET NOT IMPLEMENTED!")
}

//
// One param: 12-bits in nnn
//

// LDI - load nnn into i - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Annn
func (v *VM) LDI(addr uint16) {
	v.index = addr
}

// JP - jump to addr nnn - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#1nnn
func (v *VM) JP(addr uint16) {
	v.pc = addr
}

//
// One param: nibble in x
//

// LF - load font
func (v *VM) LF(reg uint8) {
	console.Error("***** LF NOT IMPLEMENTED!")
}

//
// Two params: x (nibble) indicating a V register, and nn byte
//

// LDVB - load byte nn into register Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#6xkk
func (v *VM) LDVB(reg uint8, byteData uint8) {
	v.registers[reg] = byteData
}

// ADDVB - add byte nn to value in register Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#7xkk
func (v *VM) ADDVB(reg uint8, byteData uint8) {
	v.registers[reg] = v.registers[reg] + byteData
}

//
// Three params: x & y (nibbles) indicating V registers, and n nibble
//

func (v *VM) DRW(reg1 uint8, reg2 uint8, height uint8) {
	x := v.registers[reg1] % DisplayWidth
	y := v.registers[reg2] % DisplayHeight
	v.registers[0xF] = 0

	// FIXME: Handle edge cases, literally... need to cope with edges of display

	var row byte
	for row = 0; row < height; row++ {
		spriteByte := v.memory[v.index+uint16(row)]
		var xline byte
		for xline = 0; xline < 8; xline++ {

			// Get bit from sprite - why this needs to be reversed I don't know!
			spriteBit := (spriteByte>>(7-xline))&1 == 1
			// Get bit from display
			displayBit := v.display[x+xline][y+row]
			// XOR logic and setting of VF
			if displayBit && spriteBit {
				v.display[x+xline][y+row] = false
				v.registers[0xF] = 1
			}
			if !displayBit && spriteBit {
				v.display[x+xline][y+row] = true
			}
		}
	}
}
