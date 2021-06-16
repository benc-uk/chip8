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
func (v *VM) insCLS() {
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			v.display[x][y] = false
		}
	}
}

// RET - Return - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#00EE
func (v *VM) insRET() {
	if len(v.stack) == 0 {
		return
	}

	// pop from stack and set pc
	i := len(v.stack) - 1
	stackAddr := v.stack[i]
	v.stack = v.stack[:i]
	v.pc = stackAddr
}

//
// One param: 12-bits in nnn
//

// LD I, addr - load nnn into i - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Annn
func (v *VM) insLDi(addr uint16) {
	v.index = addr
}

// JP addr - jump to addr nnn - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#1nnn
func (v *VM) insJP(addr uint16) {
	v.pc = addr
}

// CALL addr - put PC on stack & jump to addr nnn - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#2nnn
func (v *VM) insCALL(addr uint16) {
	v.stack = append(v.stack, v.pc)
	v.pc = addr
}

//
// One param: nibble in x
//

// LD F, Vx - load addr of font sprite for value in Vx into i - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx29
func (v *VM) insLDf(reg uint8) {
	// NOTE: Each font sprite is 5 bytes "high"
	val := uint16(v.registers[reg]) * 5
	v.index = FontBase + val
}

//
// Two params: x (nibble) indicating a V register, and nn byte
//

// LD Vx, byte - load byte nn into register Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#6xkk
func (v *VM) insLDvb(reg uint8, byteData uint8) {
	v.registers[reg] = byteData
}

// ADD Vx, byte - add byte nn to value in register Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#7xkk
func (v *VM) insADDvb(reg uint8, byteData uint8) {
	v.registers[reg] = v.registers[reg] + byteData
}

// SE Vx, byte - if byte nn == value in register Vx, advance pc - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3xkk
func (v *VM) insSEvb(reg uint8, byteData uint8) {
	if v.registers[reg] == byteData {
		v.pc += 2
	}
}

// SNE Vx, byte - if byte nn != value in register Vx, advance pc - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3xkk
func (v *VM) insSNEvb(reg uint8, byteData uint8) {
	if v.registers[reg] != byteData {
		v.pc += 2
	}
}

//
// Two params: x and y (nibbles) both indicate V regsiters, n not used
//

// SE Vx, Vy - skip and inc PC if Vx == Vy - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#5xy0
func (v *VM) insSExy(regx uint8, regy uint8) {
	if v.registers[regx] == v.registers[regy] {
		v.pc += 2
	}
}

// LD Vx, Vy - place value of Vy into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy0
func (v *VM) insLDxy(regx uint8, regy uint8) {
	v.registers[regx] = v.registers[regy]
}

// OR Vx, Vy - bitwise OR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy1
func (v *VM) insORxy(regx uint8, regy uint8) {
	v.registers[regx] |= v.registers[regy]
}

// AND Vx, Vy - bitwise OR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy2
func (v *VM) insANDxy(regx uint8, regy uint8) {
	v.registers[regx] &= v.registers[regy]
}

// XOR Vx, Vy - bitwise XOR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy3
func (v *VM) insXORxy(regx uint8, regy uint8) {
	v.registers[regx] ^= v.registers[regy]
}

// ADD Vx, Vy - Add Vx and Vy, store result into Vx. SETS VF - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy4
func (v *VM) insADDxy(regx uint8, regy uint8) {
	regxPrev := v.registers[regx]
	v.registers[regx] = v.registers[regx] + v.registers[regy]
	if v.registers[regx] < regxPrev {
		v.registers[0xF] = 1
	} else {
		v.registers[0xF] = 0
	}
}

// SUB Vx, Vy - Sub Vy from Vx, store result into Vx. SETS VF - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy5
func (v *VM) insSUBxy(regx uint8, regy uint8) {
	v.registers[regx] -= v.registers[regy]
}

// SHR Vx - bit 0 of Vx into VF, shift Vx to divide by 2 - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy6
func (v *VM) insSHRxy(regx uint8, regy uint8) {
	v.registers[0xF] = v.registers[regx] & 1
	v.registers[regx] >>= 1
	// v.registers[0xF] = v.registers[regy] & 0x1
	// v.registers[regx] = v.registers[regy] >> 1
}

// SUBN Vx, Vy - bitwise XOR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy7
func (v *VM) insSUBNxy(regx uint8, regy uint8) {
	console.Error("NOT IMPLEMENTED")
}

// SHL Vx, Vy - bitwise XOR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xyE
func (v *VM) insSHLxy(regx uint8, regy uint8) {
	console.Error("NOT IMPLEMENTED")
}

//
// Three params: x & y (nibbles) indicating V registers, and n nibble
//

// DRW Vx, Vy, nibble - Draw sprite located at i for n bytes at x, y - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Dxyn
func (v *VM) insDRW(reg1 uint8, reg2 uint8, height uint8) {
	x := v.registers[reg1] % DisplayWidth
	y := v.registers[reg2] % DisplayHeight
	v.registers[0xF] = 0

	// FIXME: Handle edge cases, literally... need to cope with edges of display

	var row byte
	for row = 0; row < height; row++ {
		spriteByte := v.memory[v.index+uint16(row)]
		var xbit byte
		for xbit = 0; xbit < 8; xbit++ {
			// Get bit from sprite - we need to draw left to right, so we start at MSB
			spriteBit := (spriteByte>>(7-xbit))&1 == 1
			// Get bit from display
			displayBit := v.display[x+xbit][y+row]
			// XOR logic and setting of VF
			if displayBit && spriteBit {
				v.display[x+xbit][y+row] = false
				v.registers[0xF] = 1
			}
			if !displayBit && spriteBit {
				v.display[x+xbit][y+row] = true
			}
		}
	}
}
