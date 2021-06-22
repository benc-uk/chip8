//
// CHIP-8 - Implementation of opcodes / instructions here
// Ben C, June 2021
// Notes:

// See http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3.1

package chip8

import (
	"math/rand"
)

//
// Zero params
//

// CLS - clear screen - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#00E0
func (v *VM) insCLS() {
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			v.display[x][y] = 0
		}
	}
	v.DisplayUpdated = true

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

// HIGH - Enable hires (Super CHIP-8)
func (v *VM) insHIGH() {
	v.HighRes = true
}

// LOW - Disable hires (Super CHIP-8)
func (v *VM) insLOW() {
	v.HighRes = false
}

// LOW - Scroll right (Super CHIP-8)
func (v *VM) insSCRR() {
	for y := 0; y < DisplayHeight; y++ {
		var x int
		for x = DisplayWidth - 1; x >= 4; x-- {
			v.display[x-4][y] = v.display[x][y]
		}
		// wipe the last 4 pixels
		v.display[x][y] = 0
		v.display[x-1][y] = 0
		v.display[x-2][y] = 0
		v.display[x-3][y] = 0
	}
}

// SCRL - Scroll left (Super CHIP-8)
func (v *VM) insSCRL() {
	for y := 0; y < DisplayHeight; y++ {
		var x int
		for x = 0; x < DisplayWidth-4; x++ {
			v.display[x][y] = v.display[x+4][y]
		}
		// wipe the last 4 pixels
		v.display[x][y] = 0
		v.display[x+1][y] = 0
		v.display[x+2][y] = 0
		v.display[x+3][y] = 0
	}
}

//
// One param: 12-bits in nnn
//

// LD I, addr - load nnn into i - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Annn
func (v *VM) insLDI(addr uint16) {
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

// JP V0, addr - jump to addr nnn +V0 - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Bnnn
func (v *VM) insJPV0(addr uint16) {
	v.pc = addr + uint16(v.registers[0])
}

//
// One param: nibble in n
//

// SCRD n - Scroll down n pixels (Super CHIP-8)
func (v *VM) insSCRD(n byte) {
	var y uint8
	for y = DisplayHeight - 1; y >= n; y-- {
		for x := 0; x < DisplayWidth; x++ {

			v.display[x][y] = v.display[x][y-n]
		}
	}
	//wipe the last 4 pixels
	for y = 0; y < n; y++ {
		for x := 0; x < DisplayWidth; x++ {
			v.display[x][y] = 0
		}
	}
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

// ADD I, Vx - Add value in Vx to i - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx1E
func (v *VM) insADDIx(reg uint8) {
	v.index += uint16(v.registers[reg])
}

// SKP Vx - Skip if key with val Vx is held - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Ex9E
func (v *VM) insSKPx(reg uint8) {
	for _, k := range v.Keys {
		if k == v.registers[reg] {
			v.pc += 2
			return
		}
	}
}

// SKNP Vx - Skip if key with val Vx is not held - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#ExA1
func (v *VM) insSKNPx(reg uint8) {
	keyUp := true
	for _, k := range v.Keys {
		if k == v.registers[reg] {
			keyUp = false
			break
		}
	}
	if keyUp {
		v.pc += 2
	}
}

// LD Vx, DT - store the delay timer val into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx07
func (v *VM) insLDxDT(reg uint8) {
	v.registers[reg] = v.delayTimer
}

// LD DT, Vx - store Vx into the delay timer - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx15
func (v *VM) insLDDTx(reg uint8) {
	v.delayTimer = v.registers[reg]
}

// LD ST, Vx - store Vx into the sound timer - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx18
func (v *VM) insLDSTx(reg uint8) {
	v.soundTimer = v.registers[reg]
}

// LD B, Vx - store BCD version of Vx into mem i (3 bytes) - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx33
func (v *VM) insLDBx(reg uint8) {
	v.memory[v.index] = v.registers[reg] / 100
	v.memory[v.index+1] = v.registers[reg] % 100 / 10
	v.memory[v.index+2] = v.registers[reg] % 10
}

// LD [I], Vx - store reg V0 through Vx into mem i - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx55
func (v *VM) insLDIx(reg uint8) {
	for ix := uint16(0); ix <= uint16(reg); ix++ {
		addr := v.index + ix
		if addr >= memSize {
			break
		}
		v.memory[addr] = v.registers[ix]
	}
}

// LD Vx, [I] - load reg V0 through Vx from mem i - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx65
func (v *VM) insLDxI(reg uint8) {
	for ix := uint16(0); ix <= uint16(reg); ix++ {
		if v.index+ix >= memSize {
			break
		}
		v.registers[ix] = v.memory[v.index+ix]
	}
}

// LD Vx, K - Wait for any key in Vx to be pressed - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Fx0A
func (v *VM) insLDxK(reg uint8) {

	if len(v.Keys) > 0 {
		// Get last key pressed if there are multiple and exit the PC loop
		v.registers[reg] = v.Keys[0]
		return
	}

	// Madness, *decrement* the PC to keep the fetch loop waiting here
	v.pc -= 2
}

//
// Two params: x (nibble) indicating a V register, and nn byte
//

// LD Vx, byte - load byte nn into register Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#6xkk
func (v *VM) insLDvb(reg, byteData uint8) {
	v.registers[reg] = byteData
}

// ADD Vx, byte - add byte nn to value in register Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#7xkk
func (v *VM) insADDvb(reg, byteData uint8) {
	v.registers[reg] = v.registers[reg] + byteData
}

// SE Vx, byte - if byte nn == value in register Vx, advance pc - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3xkk
func (v *VM) insSEvb(reg, byteData uint8) {
	if v.registers[reg] == byteData {
		v.pc += 2
	}
}

// SNE Vx, byte - if byte nn != value in register Vx, advance pc - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#3xkk
func (v *VM) insSNEvb(reg, byteData uint8) {
	if v.registers[reg] != byteData {
		v.pc += 2
	}
}

// RND Vx, byte - random value AND'ed with byte nn store in Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Cxkk
func (v *VM) insRNDvb(reg, byteData uint8) {
	r := uint8(rand.Intn(256))
	v.registers[reg] = r & byteData
}

//
// Two params: x and y (nibbles) both indicate V regsiters, n not used
//

// SE Vx, Vy - skip and inc PC if Vx == Vy - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#5xy0
func (v *VM) insSExy(regx, regy uint8) {
	if v.registers[regx] == v.registers[regy] {
		v.pc += 2
	}
}

// LD Vx, Vy - place value of Vy into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy0
func (v *VM) insLDxy(regx, regy uint8) {
	v.registers[regx] = v.registers[regy]
}

// OR Vx, Vy - bitwise OR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy1
func (v *VM) insORxy(regx, regy uint8) {
	v.registers[regx] |= v.registers[regy]
}

// AND Vx, Vy - bitwise OR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy2
func (v *VM) insANDxy(regx, regy uint8) {
	v.registers[regx] &= v.registers[regy]
}

// XOR Vx, Vy - bitwise XOR Vx and Vy, store result into Vx - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy3
func (v *VM) insXORxy(regx, regy uint8) {
	v.registers[regx] ^= v.registers[regy]
}

// ADD Vx, Vy - Add Vx and Vy, store result into Vx. SETS VF - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy4
func (v *VM) insADDxy(regx, regy uint8) {
	regxPrev := v.registers[regx]
	v.registers[regx] = v.registers[regx] + v.registers[regy]
	if v.registers[regx] < regxPrev {
		v.SetFlag(1)
	} else {
		v.SetFlag(0)
	}
}

// SUB Vx, Vy - Sub Vy from Vx, store result into Vx. Sets VF - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy5
func (v *VM) insSUBxy(regx, regy uint8) {
	if v.registers[regx] > v.registers[regy] {
		v.SetFlag(1)
	} else {
		v.SetFlag(0)
	}
	v.registers[regx] -= v.registers[regy]
}

// SHR Vx - bit 0 of Vx into VF, shift Vx to divide by 2 - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy6
func (v *VM) insSHRxy(regx, regy uint8) {
	v.registers[0xF] = v.registers[regx] & 1
	v.registers[regx] >>= 1
	// v.registers[0xF] = v.registers[regy] & 0x1
	// v.registers[regx] = v.registers[regy] >> 1
}

// SUBN Vx, Vy - Sub Vx from Vy into Vx. Sets VF - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xy7
func (v *VM) insSUBNxy(regx, regy uint8) {
	// IMPORTANT: overflow when greater OR equal
	if v.registers[regy] >= v.registers[regx] {
		v.SetFlag(1)
	} else {
		v.SetFlag(0)
	}
	v.registers[regx] = v.registers[regy] - v.registers[regx]
}

// SHL Vx, Vy - Most sig bit of Vx into VF, shift Vx left to mult by 2 - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#8xyE
func (v *VM) insSHLxy(regx, regy uint8) {
	v.registers[0xF] = v.registers[regx] >> 7
	v.registers[regx] <<= 1
}

// SHL Vx, Vy - Skip PC if Vx != Vy - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#9xy0
func (v *VM) insSNExy(regx, regy uint8) {
	if v.registers[regx] != v.registers[regy] {
		v.pc += 2
	}
}

//
// Three params: x & y (nibbles) indicating V registers, and n nibble
//

// DRW Vx, Vy, nibble - Draw sprite located at i for n bytes at x, y - http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#Dxyn
func (v *VM) insDRW(reg1, reg2, height uint8) {
	x := v.registers[reg1] % DisplayWidth
	y := v.registers[reg2] % DisplayHeight
	v.SetFlag(0)

	if height == 0 && v.HighRes {
		draw16Sprite(x, y)
		v.DisplayUpdated = true
		return
	}

	var row uint8
	for row = 0; row < height; row++ {
		if y+row >= DisplayHeight {
			return
		}
		spriteByte := v.memory[v.index+uint16(row)]
		for xbit := uint8(0); xbit < 8; xbit++ {
			if x+xbit >= DisplayWidth {
				continue
			}
			// Get bit from sprite - we need to draw left to right, so we start at MSB
			//spriteBit := (spriteByte & (0x80 >> xbit))
			spriteBit := (spriteByte >> (7 - xbit)) & 1
			// Get bit from display
			displayBit := v.display[x+xbit][y+row]

			// XOR logic and setting of VF
			if spriteBit == 1 && displayBit == 1 {
				v.SetFlag(1)
				v.display[x+xbit][y+row] = 0
			}

			if spriteBit != displayBit {
				v.display[x+xbit][y+row] = 1
			}

		}
	}

	v.DisplayUpdated = true
}

func draw16Sprite(x, y byte) {
	var row byte
	for row = 0; row < 16; row++ {
		if y+row >= DisplayHeight {
			return
		}

	}
}
