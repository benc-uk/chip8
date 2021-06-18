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
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}

//
// Zero params
//

func Test_InsCLS(t *testing.T) {
	v := vmForTest()
	v.insCLS()
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			p := v.display[x][y]
			assert.Equal(t, p, false, "pixel is not off")
		}
	}
}

func Test_InsRET(t *testing.T) {
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

func Test_InsLDi(t *testing.T) {
	v := vmForTest()
	var addr1 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	v.insLDI(addr1)
	assert.Equal(t, v.index, addr1, "index wasn't set correctly")
}

func Test_InsJP(t *testing.T) {
	v := vmForTest()
	var addr1 uint16 = uint16(ProgBase + rand.Intn(0x3000))
	v.insJP(addr1)
	assert.Equal(t, v.pc, addr1, "pc wasn't set correctly")
}

func Test_InsCALL(t *testing.T) {
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

func Test_InsLDf(t *testing.T) {
	v := vmForTest()
	char := uint8(0x2)
	v.registers[6] = char
	v.insLDf(6)
	assert.Equal(t, v.index, uint16(FontBase+char*5), "index wasn't set correctly")
}

func Test_InsADDIx(t *testing.T) {
	var r1 = uint8(1)
	v := vmForTest()
	prevI := v.index
	v.insADDIx(r1)
	assert.Equal(t, v.index, prevI+uint16(v.registers[r1]), "index wasn't set correctly")
}

//
// Two params: x (nibble) indicating a V register, and nn byte
//

func Test_InsLDvb(t *testing.T) {
	v := vmForTest()
	r := uint8(0x7)
	data := uint8(0xE)
	v.insLDvb(r, data)
	assert.Equal(t, v.registers[r], data, "register wasn't set correctly")
}

func Test_InsADDvb(t *testing.T) {
	v := vmForTest()
	var r uint8 = 0x3
	var data uint8 = 0x6

	rPrev := v.registers[r]
	v.insADDvb(r, data)
	assert.Equal(t, v.registers[r], rPrev+data, "register wasn't set correctly")
}

func Test_InsSEvb(t *testing.T) {
	v := vmForTest()
	pcPrev := v.pc
	var r uint8 = 0xD
	var data uint8 = 0x2

	v.registers[r] = data
	v.insSEvb(r, data)
	assert.Equal(t, v.pc, pcPrev+2, "pc wasn't set correctly")
}

func Test_InsSNEvb(t *testing.T) {
	v := vmForTest()
	pcPrev := v.pc
	var r uint8 = 0x3
	var data uint8 = 0x8

	v.registers[r] = data
	v.insSNEvb(r, data)
	assert.Equal(t, v.pc, pcPrev, "pc wasn't set correctly")
	v.registers[r] = 0xED
	v.insSNEvb(r, data)
	assert.Equal(t, v.pc, pcPrev+2, "pc is incorrect, should be +2")
}

//
// Two params: x and y (nibbles) both indicate V regsiters, n not used
//

func Test_InsSExy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x7
	var r2 uint8 = 0xA
	v.registers[r1] = 5
	v.registers[r2] = 5
	pcPrev := v.pc
	v.insSExy(r1, r2)
	assert.Equal(t, v.pc, pcPrev+2, "pc is incorrect, should be +2")

	v.registers[r1] = 8
	v.registers[r2] = 3
	pcPrev = v.pc
	v.insSExy(r1, r2)
	assert.Equal(t, v.pc, pcPrev, "pc is incorrect")
}

func Test_InsLDxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x1
	var r2 uint8 = 0x4

	v.insLDxy(r1, r2)
	assert.Equal(t, v.registers[r1], v.registers[r2], "register content is wrong")
}

func Test_InsORxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x7
	var r2 uint8 = 0xB
	r1Prev := v.registers[r1]

	v.insORxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev|v.registers[r2], "register content is wrong")
}

func Test_InsANDxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0x8
	var r2 uint8 = 0x2
	r1Prev := v.registers[r1]

	v.insANDxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev&v.registers[r2], "register content is wrong")
}

func Test_InsXORxy(t *testing.T) {
	v := vmForTest()
	var r1 uint8 = 0xD
	var r2 uint8 = 0x1
	r1Prev := v.registers[r1]

	v.insXORxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev^v.registers[r2], "register content is wrong")
}

func Test_InsADDxy(t *testing.T) {
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

func Test_InsSUBxy(t *testing.T) {
	var r1 = uint8(5)
	var r2 = uint8(6) // must be different
	v := vmForTest()
	r1Prev := v.registers[r1]
	flag := uint8(0)
	if v.registers[r1] > v.registers[r2] {
		flag = uint8(1)
	}
	v.insSUBxy(r1, r2)

	assert.Equal(t, v.registers[r1], r1Prev-v.registers[r2], "register content is wrong")
	assert.Equal(t, v.GetFlag(), flag, "flag register incorrect")
}

func Test_InsSHRxy(t *testing.T) {
	var r1 = randRegister()
	var r2 = randRegister()
	v := vmForTest()

	r1Prev := v.registers[r1]
	leastbit := v.registers[r1] & 1
	v.insSHRxy(r1, r2)
	assert.Equal(t, v.registers[r1], r1Prev>>1, "register content is wrong")
	assert.Equal(t, v.registers[0xF], leastbit, "flag content is wrong")
}

func Test_InsSUBNxy(t *testing.T) {
	var r1 = uint8(5)
	var r2 = uint8(6) // must be different
	v := vmForTest()
	v.registers[r1] = 30
	v.registers[r2] = 245
	r1Prev := v.registers[r1]
	v.insSUBNxy(r1, r2)
	flag := uint8(0)
	if v.registers[r2] > v.registers[r1] {
		flag = uint8(1)
	}
	assert.Equal(t, v.registers[r1], v.registers[r2]-r1Prev, "register content is wrong")
	assert.Equal(t, flag, v.GetFlag(), "flag register incorrect")
}

func Test_InsSHLxy(t *testing.T) {
	var r1 = randRegister()
	var r2 = randRegister()
	v := vmForTest()
	r1Prev := v.registers[r1]
	v.insSHLxy(r1, r2)
	newFlag := uint8(0)
	if r1Prev > 127 {
		newFlag = uint8(1)
	}
	assert.Equal(t, v.registers[r1], r1Prev*2, "register content is wrong")
	assert.Equal(t, v.GetFlag(), newFlag, "flag register incorrect")
}

func Test_InsSNExy(t *testing.T) {
	var r1 = uint8(1)
	var r2 = uint8(2) // must be different
	v := vmForTest()

	v.registers[r1] = uint8(44)
	v.registers[r2] = uint8(100)
	pcPrev := v.pc
	v.insSNExy(r1, r2)
	assert.Equal(t, v.pc, pcPrev+2, "pc value should be pc+2")

	v.registers[r1] = uint8(23)
	v.registers[r2] = uint8(23)
	pcPrev = v.pc
	v.insSNExy(r1, r2)
	assert.Equal(t, v.pc, pcPrev, "pc value is wrong")
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
	v.registers[0xf] = 4

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
