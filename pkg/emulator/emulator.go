package emulator

import (
	"image/color"
	"log"
	"os"
	"runtime"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/benc-uk/chip8/pkg/console"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Version is the emulator version
var Version = "0.0.2"
var pixelColour = color.RGBA{0x00, 0xff, 0x00, 0xff}

// Wrapper for ebiten implements the ebiten.Game interface
type chip8Emulator struct {
	vm        *chip8.VM
	display   *ebiten.Image
	pixelSize int
	speed     int
	paused    bool
	pgmData   []byte

	audioContext *audio.Context
	bleeper      *audio.Player
}

func Start(program []byte, debug bool, speed int, pixelSize int) {
	console.Infof("Starting CHIP-8 emulator version v%s\n\n", Version)

	if runtime.GOARCH == "js" || runtime.GOOS == "js" {
		ebiten.SetFullscreen(true)
	}

	// Create a new CHIP-8 virtual machine, and load program into it
	vm := chip8.NewVM()
	vm.SetDebug(debug)
	vm.LoadProgram(program)

	// Wrap the VM in an chip8Emulator to allow us to use ebiten with it
	emu := &chip8Emulator{
		vm:           vm,
		display:      ebiten.NewImage(chip8.DisplayWidth*pixelSize, chip8.DisplayHeight*pixelSize),
		pixelSize:    pixelSize,
		speed:        speed,
		audioContext: audio.NewContext(44100),
		pgmData:      program,
	}

	ebiten.SetWindowSize(chip8.DisplayWidth*pixelSize, chip8.DisplayHeight*pixelSize)
	ebiten.SetWindowTitle("Go CHIP-8 v" + Version)
	//ebiten.SetMaxTPS(ebiten.UncappedTPS)
	ebiten.SetVsyncEnabled(false)

	if err := ebiten.RunGame(emu); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Update is called every tick (1/60 [s] by default).
func (e *chip8Emulator) Update() error {
	// Read the keyboard
	e.readKeyboard()

	// Play sound
	e.playSound()

	// Main emulator processor loop, we execute a number of CHIP-8 processor cycles
	// Depending on the speed
	for c := 0; c < e.speed; c++ {
		// Handle pausing and stepping through code
		if e.paused && !inpututil.IsKeyJustPressed(ebiten.KeyF6) {
			c--
			return nil
		}

		// This advances the processor one tick/cycle
		runtimeError := e.vm.Cycle()

		// Handle runtime errors which are always fatal
		if runtimeError != nil {
			se, isSystemError := runtimeError.(chip8.SystemError)
			code := 50 // Default code
			if isSystemError {
				code = se.Code()
			}

			log.Printf("Unrecoverable system error: %s", runtimeError.Error())
			os.Exit(code)
		}
	}

	return nil
}

// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (e *chip8Emulator) Draw(screen *ebiten.Image) {
	if e.vm.DisplayUpdated {
		e.renderDisplay()
		e.vm.DisplayUpdated = false
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(float64(e.pixelSize), float64(e.pixelSize))
	opt.Filter = ebiten.FilterNearest
	screen.DrawImage(e.display, opt)

	debugMsg := ""
	if e.paused {
		debugMsg = "PAUSED"
	}
	if e.vm.IsDebugging() {
		debugMsg += "\nDEBUGGING"
	}
	ebitenutil.DebugPrint(screen, debugMsg)
}

// Layout can control scaling
func (e *chip8Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (e *chip8Emulator) Reset() {
	e.vm.Reset()
	e.vm.LoadProgram(e.pgmData)
}
