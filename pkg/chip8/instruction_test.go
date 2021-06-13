//
// CHIP-8 - Look at me mum, I'm writing tests!
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test CALL and stack
func TestIns_CALL(t *testing.T) {
	var addr1 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	var addr2 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	v := vmForTest()
	prePC1 := v.pc
	preStackLen := len(v.stack)

	v.insCALL(addr1)
	assert.Equal(t, v.pc, addr1, "address is not in pc")
	assert.Equal(t, len(v.stack), preStackLen+1, "stack len is wrong")
	assert.Equal(t, v.stack[preStackLen], prePC1, "stack content is wrong")

	prePC2 := v.pc
	v.insCALL(addr2)
	assert.Equal(t, v.pc, addr2, "address is not in pc")
	assert.Equal(t, len(v.stack), preStackLen+2, "stack len is wrong")
	assert.Equal(t, v.stack[preStackLen+1], prePC2, "stack content is wrong")
}

func TestIns_LDxy(t *testing.T) {
	var r1 uint8 = 0x1
	var r2 uint8 = 0x4
	v := vmForTest()
	v.insLDxy(r1, r2)
	assert.Equal(t, v.registers[r1], v.registers[r2], "register content is wrong")
}

func vmForTest() *VM {
	v := NewVM(true)
	for i := range v.registers {
		r := rand.Int()
		v.registers[i] = uint8(r)
	}
	rand.Seed(time.Now().UTC().UnixNano())
	v.pc = uint16(ProgBase + rand.Intn(0x3000) + 1)
	v.index = uint16(ProgBase + rand.Intn(0x3000) + 1)
	v.stack = append(v.stack, uint16(rand.Intn(0x3000)))
	//v.dump()
	return v
}
