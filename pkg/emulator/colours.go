package emulator

import (
	"fmt"
	"image/color"
)

type ColourMap map[int]int

var pallette []color.RGBA

func init() {
	pallette = make([]color.RGBA, 9)
	pallette[1] = color.RGBA{R: 255, G: 255, B: 255, A: 255} // white
	pallette[2] = color.RGBA{R: 255, G: 0, B: 0, A: 255}     // red
	pallette[3] = color.RGBA{R: 0, G: 255, B: 0, A: 255}     // green
	pallette[4] = color.RGBA{R: 0, G: 0, B: 255, A: 255}     // blue
	pallette[5] = color.RGBA{R: 255, G: 0, B: 255, A: 255}   // magenta
	pallette[6] = color.RGBA{R: 255, G: 255, B: 0, A: 255}   // yellow
	pallette[7] = color.RGBA{R: 0, G: 255, B: 255, A: 255}   // cyan
	pallette[8] = color.RGBA{R: 255, G: 120, B: 0, A: 255}   // orange
}

var joustMap = ColourMap{
	0xA6A: 2,
	0xA62: 2,
}

var carMap = ColourMap{
	0x32c: 1, // roadside
	0x330: 3, // grass
	0x338: 8, // cars
}

var invadersMap = ColourMap{
	0x03B7: 2,
	0x03CF: 1,
	0x03BD: 4,
	0x03C3: 5,
}

var brixMap = ColourMap{
	0x030C: 5,
	0x0312: 1,
	0x0310: 4,
	0x030E: 6,
}

var spacejamMap = ColourMap{
	0x0396: 1,
	0x03CC: 5,
	0x0397: 7,
}

func GetColourMap(hash [16]byte) ColourMap {
	hashString := fmt.Sprintf("%X", hash)
	switch hashString {
	case "214E7C967243CC8FD9E51CCEBE248113":
		return joustMap
	case "C497BB692EA4B32A4A7B11B1373EF92F":
		return carMap
	case "4FE20B951DBC801D7F682B88E672626C":
		return invadersMap
	case "D677C1B9DE941484D718799AEBAFEBF3":
		return brixMap
	case "21D9FF1620FC2D8AEC0D6DBCDA92C35E":
		return spacejamMap
	}
	return nil
}

func (m ColourMap) getColour(address int) *color.RGBA {
	if palletteNum, ok := m[address]; ok {
		return &pallette[palletteNum]
	}
	return nil
}
