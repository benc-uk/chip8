//
// CHIP-8 - Input and output to filesystem, keyboard & display
// Ben C, June 2021
// Notes:
//

package chip8

import (
	"image/color"
	"io/ioutil"

	"github.com/benc-uk/chip8/pkg/console"
	"github.com/hajimehoshi/ebiten/v2"
)

const PixelSize = 12

func (v *VM) LoadProgramFile(fileName string) (int, error) {
	console.Infof("Loading program from disk %s\n", fileName)

	pgmBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	// Reset the machine before writing program data to memory
	v.reset()
	for i := range pgmBytes {
		v.memory[ProgBase+i] = pgmBytes[i]
	}

	console.Successf("Loaded %d bytes into memory OK\n", len(pgmBytes))
	return len(pgmBytes), nil
}

func (v *VM) RenderDisplay(screen *ebiten.Image) {
	screen.Clear()
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			if v.display[x][y] {
				drawPixel(x, y, screen)
			}
		}
	}
}

func drawPixel(x int, y int, screen *ebiten.Image) {
	c := color.RGBA{0, 0xff, 0, 0xff}
	for yi := 0; yi < PixelSize; yi++ {
		for xi := 0; xi < PixelSize; xi++ {
			screen.Set((x*PixelSize)+xi, (y*PixelSize)+yi, c)
		}
	}
}
