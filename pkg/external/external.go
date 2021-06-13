//
// CHIP-8 - Input and output to filesystem, keyboard & display
// Ben C, June 2021
// Notes:
//

package external

import (
	"image/color"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/hajimehoshi/ebiten/v2"
)

func RenderDisplay(v *chip8.VM, screen *ebiten.Image, pixelSize int) {
	screen.Clear()
	for y := 0; y < chip8.DisplayHeight; y++ {
		for x := 0; x < chip8.DisplayWidth; x++ {
			if v.DisplayValueAt(x, y) {
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
