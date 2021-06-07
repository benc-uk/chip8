//
// CHIP-8 emulator - main executable
// Ben C, June 2021
// Notes:
//

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/benc-uk/chip8/pkg/console"
	"github.com/hajimehoshi/ebiten/v2"
)

var version = "0.0.1"

const banner = `
 ██████╗  ██████╗      ██████╗██╗  ██╗██╗██████╗        █████╗ 
██╔════╝ ██╔═══██╗    ██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
██║  ███╗██║   ██║    ██║     ███████║██║██████╔╝█████╗╚█████╔╝
██║   ██║██║   ██║    ██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
╚██████╔╝╚██████╔╝    ╚██████╗██║  ██║██║██║           ╚█████╔╝
 ╚═════╝  ╚═════╝      ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝`

// Wrapper for ebiten implements the ebiten.Game interface
type emulator struct {
	vm    *chip8.VM
	tick  int
	speed int
}

func main() {
	var debugFlag = flag.Bool("debug", false, "Enable debug")
	var speedFlag = flag.Int("speed", 10, "Processor cycles per second")
	flag.Parse()

	if len(flag.Args()) < 1 {
		console.Error("Please supply filename of program/ROM to load")
		os.Exit(1)
	}
	progFile := flag.Arg(0)
	console.Info(banner)
	console.Infof("Version v%s\n\n", version)

	// Create a new CHIP-8 virtual machine
	vm := chip8.NewVM(*debugFlag)

	_, err := vm.LoadProgramFile(progFile)
	if err != nil {
		console.Errorf("Error loading program: %s", err)
		os.Exit(1)
	}

	//vm.DumpMemory(chip8.FontBase-0x1, chip8.FontBase+0x16)

	// Wrap the VM in an emulator
	emu := &emulator{
		vm:    vm,
		speed: *speedFlag,
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
	// FIXME: I think the processor loop needs to be decoupled from the ebiten update loop
	// Otherwise our CPU is max 60Hz!
	if e.tick > e.speed {
		e.tick = 0
		return e.vm.Cycle()
	}
	e.tick++
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
