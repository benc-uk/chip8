//
// CHIP-8 - Input and output to filesystem, keyboard & display
// Ben C, June 2021
// Notes:
//

package emulator

import (
	"runtime"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (e *chip8Emulator) renderDisplay() {
	e.display.Clear()
	//pixelSize := e.pixelImage.Bounds().Dx()

	for y := 0; y < chip8.DisplayHeight; y++ {
		for x := 0; x < chip8.DisplayWidth; x++ {
			if e.vm.DisplayValueAt(x, y) == 1 {
				e.display.Set(x, y, pixelColour)
				//opt := &ebiten.DrawImageOptions{}
				//opt.GeoM.Translate(float64(x), float64(y))
				//opt.CompositeMode = ebiten.CompositeModeCopy
				//e.display.DrawImage(e.pixelImage, opt)
			}
		}
	}
}

func readKeyboard(v *chip8.VM) {
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

func (e *chip8Emulator) playSound() {
	if runtime.GOOS == "windows" {
		if e.vm.GetSoundTimer() > 0 && e.bleeper == nil {
			e.bleeper = audio.NewPlayerFromBytes(e.audioContext, e.wav)
			e.bleeper.SetVolume(1)
			e.bleeper.Play()
		}
		if e.vm.GetSoundTimer() <= 0 && e.bleeper != nil {
			e.bleeper.Close()
			e.bleeper = nil
		}
	}
}
