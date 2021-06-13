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

func (v *VM) LoadProgram(pgm []byte) {
	// Reset the machine before writing program data to memory
	v.reset()
	for i := range pgm {
		v.memory[ProgBase+i] = pgm[i]
	}

	console.Successf("Loaded %d bytes into memory OK\n", len(pgm))
}

func (v *VM) LoadProgramFile(fileName string) (int, error) {
	console.Infof("Loading program from disk %s\n", fileName)

	pgmBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	v.LoadProgram(pgmBytes)

	return len(pgmBytes), nil
}

func (v *VM) RenderDisplay(screen *ebiten.Image, pixelSize int) {
	screen.Clear()
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			if v.display[x][y] {
				drawPixel(x, y, screen, pixelSize)
			}
		}
	}
}

func drawPixel(x int, y int, screen *ebiten.Image, pixelSize int) {
	c := color.RGBA{0, 0xff, 0, 0xff}
	for yi := 0; yi < pixelSize; yi++ {
		for xi := 0; xi < pixelSize; xi++ {
			screen.Set((x*pixelSize)+xi, (y*pixelSize)+yi, c)
		}
	}
}
