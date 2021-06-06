//
// CHIP-8 - Look at me mum, I'm writing tests!
// Ben C, June 2021
// Notes:
//

package chip8

import "testing"

// Test LD, I #21F
func TestDecode_0xA21F(t *testing.T) {
	o := decode(0xA21F)
	o.dump()
	if o.kind != 0xA {
		t.Error("kind is wrong")
	}
	if o.nnn != 0x21F {
		t.Error("nnn is wrong")
	}
}

// Test LD, V5 #8F
func TestDecode_0x658F(t *testing.T) {
	o := decode(0x658F)
	o.dump()
	if o.kind != 0x6 {
		t.Error("kind is wrong")
	}
	if o.x != 0x5 {
		t.Error("x is wrong")
	}
	if o.nn != 0x8F {
		t.Error("nn is wrong")
	}
}

// Test LD, V0 #05
func TestDecode_0x6E70(t *testing.T) {
	o := decode(0x6E70)
	o.dump()
	if o.kind != 0x6 {
		t.Error("kind is wrong")
	}
	if o.x != 0xE {
		t.Error("x is wrong")
	}
	if o.nn != 0x70 {
		t.Error("nn is wrong")
	}
}

// Test DRW 4, B, 7
func TestDecode_0xD4B7(t *testing.T) {
	o := decode(0xD4B7)
	o.dump()
	if o.kind != 0xD {
		t.Error("kind is wrong")
	}
	if o.x != 0x4 {
		t.Error("x is wrong")
	}
	if o.y != 0xB {
		t.Error("y is wrong")
	}
	if o.n != 0x7 {
		t.Error("n is wrong")
	}
}
