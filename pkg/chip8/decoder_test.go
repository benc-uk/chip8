//
// CHIP-8 - Look at me mum, I'm writing tests!
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type decodeTest struct {
	rawOpcode      uint16
	expectedOpcode Opcode
}

var testCases = []decodeTest{
	{
		// Test DRW 4, B, 7
		rawOpcode:      0xD4B7,
		expectedOpcode: Opcode{kind: 0xD, x: 0x4, y: 0xB, n: 0x7, nn: 0xB7, nnn: 0x4b7},
	},
	{
		// Test LD, I #21F
		rawOpcode:      0xA21F,
		expectedOpcode: Opcode{kind: 0xA, x: 0x2, y: 0x1, n: 0xF, nn: 0x1F, nnn: 0x21F},
	},
	{
		// Test LD, V5 #8F
		rawOpcode:      0x658F,
		expectedOpcode: Opcode{kind: 0x6, x: 0x5, y: 0x8, n: 0xF, nn: 0x8F, nnn: 0x58F},
	},
	{
		// Test LD, VE #70
		rawOpcode:      0x6E70,
		expectedOpcode: Opcode{kind: 0x6, x: 0xE, y: 0x7, n: 0x0, nn: 0x70, nnn: 0xE70},
	},
}

func TestDecoder(t *testing.T) {
	for _, tcase := range testCases {
		t.Run(fmt.Sprintf("Decode_%04X", tcase.rawOpcode), func(t *testing.T) {
			o := decode(tcase.rawOpcode)
			assert.Equal(t, o, tcase.expectedOpcode, "opcode was not decoded correctly")
		})
	}
}
