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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// var audioContext *audio.Context
// var wav []byte

func RenderDisplay(v *chip8.VM, screen *ebiten.Image, pixelSize int) {
	screen.Clear()
	for y := 0; y < chip8.DisplayHeight; y++ {
		for x := 0; x < chip8.DisplayWidth; x++ {
			if v.DisplayValueAt(x, y) == 1 {
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

func ReadKeyboard(v *chip8.VM) {
	v.Keys = []uint8{}
	for _, keycode := range inpututil.PressedKeys() {
		// This handles the "standard" mapping from PC keyboard to CHIP8 keypad num
		switch keycode {
		case ebiten.Key1:
			v.Keys = append(v.Keys, 0x1)
		case ebiten.Key2:
			v.Keys = append(v.Keys, 0x2)
		case ebiten.Key3:
			v.Keys = append(v.Keys, 0x3)
		case ebiten.Key4:
			v.Keys = append(v.Keys, 0xC)

		case ebiten.KeyQ:
			v.Keys = append(v.Keys, 0x4)
		case ebiten.KeyW:
			v.Keys = append(v.Keys, 0x5)
		case ebiten.KeyE:
			v.Keys = append(v.Keys, 0x6)
		case ebiten.KeyR:
			v.Keys = append(v.Keys, 0xD)

		case ebiten.KeyA:
			v.Keys = append(v.Keys, 0x7)
		case ebiten.KeyS:
			v.Keys = append(v.Keys, 0x8)
		case ebiten.KeyD:
			v.Keys = append(v.Keys, 0x9)
		case ebiten.KeyF:
			v.Keys = append(v.Keys, 0xE)

		case ebiten.KeyZ:
			v.Keys = append(v.Keys, 0xA)
		case ebiten.KeyX:
			v.Keys = append(v.Keys, 0x0)
		case ebiten.KeyC:
			v.Keys = append(v.Keys, 0xB)
		case ebiten.KeyV:
			v.Keys = append(v.Keys, 0xF)
		}
	}
}

func PlaySound(v *chip8.VM) {
	if v.GetSoundTimer() > 0 {
		// bleepPlay := audio.NewPlayerFromBytes(audioContext, wav)
		// bleepPlay.SetVolume(1)
		// bleepPlay.Play()
	}
}

// func init() {
// 	sampleRate := 32000

// 	audioContext = audio.NewContext(sampleRate)
// 	var err error
// 	wav, err = os.ReadFile("bleep.wav")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }
