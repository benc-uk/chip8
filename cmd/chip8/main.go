//
// CHIP-8 emulator - main executable
// Ben C, June 2021
// Notes:
//

package main

import (
	"fmt"
	"os"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/benc-uk/chip8/pkg/console"
	"github.com/hajimehoshi/ebiten/v2"
)

var version = "0.0.1"

const debug = false
const cyclesPerSecond = 60

const banner = `
 ██████╗  ██████╗      ██████╗██╗  ██╗██╗██████╗        █████╗ 
██╔════╝ ██╔═══██╗    ██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
██║  ███╗██║   ██║    ██║     ███████║██║██████╔╝█████╗╚█████╔╝
██║   ██║██║   ██║    ██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
╚██████╔╝╚██████╔╝    ╚██████╗██║  ██║██║██║           ╚█████╔╝
 ╚═════╝  ╚═════╝      ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝ 
                                                               `

// Wrapper for ebiten implements the ebiten.Game interface
type emulator struct {
	vm   *chip8.VM
	tick int64
}

func main() {
	console.Info(banner)
	console.Infof("Version v%s\n\n", version)

	// Create a new CHIP-8 virtual machine
	vm := chip8.NewVM(debug)

	_, err := vm.LoadProgramFile("roms/ibm.ch8")
	if err != nil {
		console.Errorf("Error loading program: %s", err)
		os.Exit(1)
	}

	// Wrap the VM in an emulator
	emu := &emulator{
		vm: vm,
	}

	ebiten.SetWindowSize(chip8.DisplayWidth*chip8.PixelSize, chip8.DisplayHeight*chip8.PixelSize)
	//ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Go CHIP-8 v" + version)
	if err := ebiten.RunGame(emu); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Update is called every tick (1/60 [s] by default).
func (e *emulator) Update() error {
	// Call the CHIP-8 processor cycle but at given rate
	e.tick++
	if e.tick%(60/cyclesPerSecond) == 0 {
		return e.vm.Cycle()
	}
	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (e *emulator) Draw(screen *ebiten.Image) {
	e.vm.RenderDisplay(screen)
}

// Layout can control scaling
func (e *emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
