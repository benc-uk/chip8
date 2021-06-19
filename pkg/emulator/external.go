//
// CHIP-8 - Input and output to filesystem, keyboard & display
// Ben C, June 2021
// Notes:
//

package emulator

import (
	"runtime"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/benc-uk/chip8/pkg/sounds"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (e *chip8Emulator) renderDisplay() {
	e.display.Clear()

	for y := 0; y < chip8.DisplayHeight; y++ {
		for x := 0; x < chip8.DisplayWidth; x++ {
			if e.vm.DisplayValueAt(x, y) == 1 {
				e.display.Set(x, y, pixelColour)
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

// Play the bleep sound when it's active
// HACK: This just feels like a mess
func (e *chip8Emulator) playSound() {
	if runtime.GOOS != "linux" {
		// Play sound once, when timer is over zero and we aren't playing already
		if e.vm.GetSoundTimer() > 0 && e.bleeper == nil {
			e.bleeper = audio.NewPlayerFromBytes(e.audioContext, sounds.BleepWav)
			e.bleeper.SetVolume(1)
			e.bleeper.Play()
		}
		// When timer his zero, stop but only if we have an active bleeper
		if e.vm.GetSoundTimer() <= 0 && e.bleeper != nil {
			e.bleeper.Close()
			e.bleeper = nil
		}
	}
}
