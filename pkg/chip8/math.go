package chip8

func getBit(val uint8, bit int) uint8 {
	bitV := val >> (bit) & 1
	return bitV
}
