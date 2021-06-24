//
// CHIP-8 - Emulator input and output to filesystem, keyboard & display
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
	e.display.Fill(e.bgColor)

	for y := 0; y < chip8.DisplayHeight; y++ {
		for x := 0; x < chip8.DisplayWidth; x++ {
			if e.vm.DisplayValueAt(x, y) == 1 {
				if !e.vm.HighRes {
					e.display.Set(x*2, y*2, e.fgColor)
					e.display.Set(x*2+1, y*2, e.fgColor)
					e.display.Set(x*2, y*2+1, e.fgColor)
					e.display.Set(x*2+1, y*2+1, e.fgColor)
				} else {
					e.display.Set(x, y, e.fgColor)
				}
			}
		}
	}
}

func (e *chip8Emulator) readKeyboard() {
	// Emulator specific keys
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		e.paused = !e.paused
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		e.vm.SetDebug(!e.vm.IsDebugging())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		e.SoftReset()
	}

	e.vm.Keys = []uint8{}
	for _, keycode := range inpututil.PressedKeys() {
		// This handles the "standard" mapping from PC keyboard to CHIP8 keypad num
		switch keycode {
		case ebiten.Key1:
			e.vm.Keys = append(e.vm.Keys, 0x1)
		case ebiten.Key2:
			e.vm.Keys = append(e.vm.Keys, 0x2)
		case ebiten.Key3:
			e.vm.Keys = append(e.vm.Keys, 0x3)
		case ebiten.Key4:
			e.vm.Keys = append(e.vm.Keys, 0xC)

		case ebiten.KeyQ:
			e.vm.Keys = append(e.vm.Keys, 0x4)
		case ebiten.KeyW:
			e.vm.Keys = append(e.vm.Keys, 0x5)
		case ebiten.KeyE:
			e.vm.Keys = append(e.vm.Keys, 0x6)
		case ebiten.KeyR:
			e.vm.Keys = append(e.vm.Keys, 0xD)

		case ebiten.KeyA:
			e.vm.Keys = append(e.vm.Keys, 0x7)
		case ebiten.KeyS:
			e.vm.Keys = append(e.vm.Keys, 0x8)
		case ebiten.KeyD:
			e.vm.Keys = append(e.vm.Keys, 0x9)
		case ebiten.KeyF:
			e.vm.Keys = append(e.vm.Keys, 0xE)

		case ebiten.KeyZ:
			e.vm.Keys = append(e.vm.Keys, 0xA)
		case ebiten.KeyX:
			e.vm.Keys = append(e.vm.Keys, 0x0)
		case ebiten.KeyC:
			e.vm.Keys = append(e.vm.Keys, 0xB)
		case ebiten.KeyV:
			e.vm.Keys = append(e.vm.Keys, 0xF)
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
