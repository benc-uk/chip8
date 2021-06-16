//
// CHIP-8 - Look at me mum, I'm writing tests!
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test setup
func TestMain(m *testing.M) {
	//log.SetOutput(ioutil.Discard)
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}

//
// Zero params
//

func TestIns_CLS(t *testing.T) {
	v := vmForTest()
	v.insCLS()
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			p := v.display[x][y]
			assert.Equal(t, p, false, "pixel is not off")
		}
	}
}

func TestIns_RET(t *testing.T) {
	v := vmForTest()
	var addr1 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	prePC1 := v.pc

	v.insCALL(addr1)
	assert.NotEqual(t, v.pc, prePC1, "pc hasn't changed")
	v.insRET()
	assert.Equal(t, v.pc, prePC1, "pc wasn't restored")
}

//
// One param: 12-bits in nnn
//

func TestIns_LDi(t *testing.T) {
	v := vmForTest()
	var addr1 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	v.insLDi(addr1)
	assert.Equal(t, v.index, addr1, "index wasn't set correctly")
}

func TestIns_JP(t *testing.T) {
	v := vmForTest()
	var addr1 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	v.insJP(addr1)
	assert.Equal(t, v.pc, addr1, "pc wasn't set correctly")
}

func TestIns_CALL(t *testing.T) {
	v := vmForTest()
	var addr1 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	var addr2 uint16 = uint16(ProgBase + rand.Intn(0x3000))
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

//
// One param: nibble in x
//

func TestIns_LDf(t *testing.T) {
	v := vmForTest()
	char := uint8(0x2)
	v.registers[6] = char
	v.insLDf(6)
	assert.Equal(t, v.index, uint16(FontBase+char*5), "index wasn't set correctly")
}

//
// Two params: x (nibble) indicating a V register, and nn byte
//

func TestIns_LDvb(t *testing.T) {
	v := vmForTest()
	r := uint8(0x7)
	data := uint8(0xE)
	v.insLDvb(r, data)
	assert.Equal(t, v.registers[r], data, "register wasn't set correctly")
}

func TestIns_ADDvb(t *testing.T) {
	v := vmForTest()
	var r uint8 = 0x3
	var data uint8 = 0x6

	rPrev := v.registers[r]
	v.insADDvb(r, data)
	assert.Equal(t, v.registers[r], rPrev+data, "register wasn't set correctly")
}

func TestIns_SEvb(t *testing.T) {
	v := vmForTest()
	pcPrev := v.pc
	var r uint8 = 0xD
	var data uint8 = 0x2

	v.registers[r] = data
	v.insSEvb(r, data)
	assert.Equal(t, v.pc, pcPrev+2, "pc wasn't set correctly")
}

func TestIns_SNEvb(t *testing.T) {
	v := vmForTest()
	pcPrev := v.pc
	var r uint8 = 0x3
	var data uint8 = 0x8

	v.registers[r] = data
	v.insSNEvb(r, data)
	assert.NotEqual(t, v.pc, pcPrev+2, "pc wasn't set correctly")
}

//
// Two params: x and y (nibbles) both indicate V regsiters, n not used
//

func TestIns_SExy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = randRegister()
	var r2 uint8 = randRegister()
	v.registers[r1] = 5
	v.registers[r2] = 5
	pcPrev := v.pc
	v.insSExy(r1, r2)
	assert.Equal(t, v.pc, pcPrev+2, "pc is incorrect")

	v.registers[r1] = 8
	v.registers[r2] = 3
	pcPrev = v.pc
	v.insSExy(r1, r2)
	assert.Equal(t, v.pc, pcPrev, "pc is incorrect")
}

func TestIns_LDxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x1
	var r2 uint8 = 0x4

	v.insLDxy(r1, r2)
	assert.Equal(t, v.registers[r1], v.registers[r2], "register content is wrong")
}

func TestIns_ORxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x7
	var r2 uint8 = 0xB
	r1Prev := v.registers[r1]

	v.insORxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev|v.registers[r2], "register content is wrong")
}

func TestIns_ANDxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x8
	var r2 uint8 = 0x2
	r1Prev := v.registers[r1]

	v.insANDxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev&v.registers[r2], "register content is wrong")
}

func TestIns_XORxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0xD
	var r2 uint8 = 0x1
	r1Prev := v.registers[r1]

	v.insXORxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev^v.registers[r2], "register content is wrong")
}

func TestIns_ADDxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x1
	var r2 uint8 = 0x2

	// Add without overflow
	v.registers[r1] = 101
	v.registers[r2] = 60
	r1Prev := v.registers[r1]
	v.insADDxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev+v.registers[r2], "register content is wrong")
	assert.Equal(t, v.GetFlag(), uint8(0), "flag register should be zero")

	// Add with overflow
	v.registers[r1] = 234
	v.registers[r2] = 80
	r1Prev = v.registers[r1]
	v.insADDxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev+v.registers[r2], "register content is wrong")
	assert.Equal(t, v.GetFlag(), uint8(1), "flag register should be one")
}

func TestIns_SUBxy(t *testing.T) {
	var r1 uint8 = 0xD
	var r2 uint8 = 0x1
	v := vmForTest()
	r1Prev := v.registers[r1]
	v.insSUBxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev-v.registers[r2], "register content is wrong")
}

func TestIns_SHRxy(t *testing.T) {
	var r1 = randRegister()
	var r2 = randRegister()
	v := vmForTest()

	r1Prev := v.registers[r1]
	leastbit := v.registers[r1] & 1
	v.insSHRxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev>>1, "register content is wrong")
	assert.Equal(t, v.registers[0xF], leastbit, "flag content is wrong")
}

//
// ===========================================
//
func vmForTest() *VM {
	v := NewVM(true)
	for i := range v.registers {
		r := rand.Int()
		v.registers[i] = uint8(r)
	}
	v.pc = uint16(ProgBase + rand.Intn(0x3000) + 1)
	v.index = uint16(ProgBase + rand.Intn(0x3000) + 1)
	v.stack = append(v.stack, uint16(rand.Intn(0x3000)))

	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			r := rand.Intn(100)
			v.display[x][y] = (r > 50)
		}
	}

	return v
}

// Random byte value
func randByte() uint8 {
	return uint8(rand.Intn(256))
}

// Random register between 0x0 and 0xE, note register F should not be used
func randRegister() uint8 {
	return uint8(rand.Intn(15))
}

// Generates a random address above the program line
func randAddress() uint16 {
	return uint16(rand.Intn(memSize-ProgBase) + ProgBase)
}
