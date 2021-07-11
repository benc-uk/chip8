//
// CHIP-8 - SUPER CHIP-8 instructions and opcodes
// Ben C, June 2021
// Notes:

package chip8

import "github.com/benc-uk/chip8/pkg/console"

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
	v.DisplayUpdated = true
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
	v.DisplayUpdated = true
}

// SCRD n - Scroll down n pixels (Super CHIP-8)
func (v *VM) insSCRD(n byte) {
	var y uint8
	for y = DisplayHeight - 1; y >= n; y-- {
		for x := 0; x < DisplayWidth; x++ {

			v.display[x][y] = v.display[x][y-n]
		}
	}
	// Wipe the remaining top n rows of pixels
	for y = 0; y < n; y++ {
		for x := 0; x < DisplayWidth; x++ {
			v.display[x][y] = 0
		}
	}
	v.DisplayUpdated = true
}

// LD SF, Vx - load addr of big/super font sprite for value of Vx into i (Super CHIP-8)
func (v *VM) insLDSf(reg uint8) {
	// NOTE: Each font sprite is 10 bytes "high"
	val := uint16(v.registers[reg]) * 10
	v.index = FontLargeBase + val
}

// EXIT - Exit the emulator, instead we just hang in a loop, it's less abrupt
func (v *VM) insEXIT() {
	v.pc -= 2
}

// Called from insDRW to draw a 16x16 sprite in super mode
// Note not really an instruction but nowhere better to put it
func (v *VM) draw16Sprite(x, y byte) {
	var row uint8
	v.SetFlag(0)

	for row = 0; row < 16; row++ {
		if y+row >= DisplayHeight {
			return
		}
		// Note we multiply by two as we process two bytes per row
		spriteByte1 := v.memory[v.index+uint16(row*2)]
		spriteByte2 := v.memory[v.index+uint16(row*2+1)]

		var bitIndex uint8
		for bitIndex = 0; bitIndex < 16; bitIndex++ {
			if x+bitIndex >= DisplayWidth {
				continue
			}

			var spriteBit uint8
			if bitIndex < 8 {
				spriteBit = (spriteByte1 >> (7 - bitIndex)) & 1
			} else {
				spriteBit = (spriteByte2 >> (15 - bitIndex)) & 1
			}

			// Get bit from display
			displayBit := v.display[x+bitIndex][y+row]

			// XOR logic and setting of VF
			if spriteBit == 1 && displayBit > 0 {
				v.SetFlag(1)
				v.display[x+bitIndex][y+row] = 0
				continue
			}

			if spriteBit == 1 && displayBit == 0 {
				// !NOTE! This is HIGHLY unorthodox, we store the sprite address here ONLY for colour remapping support
				// If you looking at this code and writing your own CHIP-8 emulator set this to 1 !
				v.display[x+bitIndex][y+row] = v.index

				// Only for debugging sprite values
				if !v.debugSpriteMap[v.index] && v.DebugLevel == DebugLevelSprite {
					// Only output sprite message the first time we see this sprite address
					v.debugSpriteMap[v.index] = true
					console.Successf("DRAWING SPRITE %04X at, %d,%d\n", v.index, x, y)
					v.debugSpriteMap[v.index] = true
				}
			}
		}
	}
}
