package emulator

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"runtime"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/benc-uk/chip8/pkg/console"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Version is the emulator version
var Version = "0.0.1"
var pixelColour = color.RGBA{0x00, 0xff, 0x00, 0xff}

// Wrapper for ebiten implements the ebiten.Game interface
type chip8Emulator struct {
	vm        *chip8.VM
	errorChan chan error
	display   *ebiten.Image
	pixelSize int

	audioContext *audio.Context
	wav          []byte
	bleeper      *audio.Player
}

func Start(program []byte, debug bool, delay int, pixelSize int) {
	console.Infof("Starting CHIP-8 emulator version v%s\n\n", Version)

	if runtime.GOARCH == "js" || runtime.GOOS == "js" {
		ebiten.SetFullscreen(true)
	}

	// Create a new CHIP-8 virtual machine, and load program into it
	vm := chip8.NewVM(debug)
	vm.LoadProgram(program)

	// Wrap the VM in an chip8Emulator to allow us to use ebiten with it
	emu := &chip8Emulator{
		vm:        vm,
		errorChan: make(chan error),
		display:   ebiten.NewImage(chip8.DisplayWidth*pixelSize, chip8.DisplayHeight*pixelSize),
		pixelSize: pixelSize,
	}

	if runtime.GOOS == "windows" {
		emu.audioContext = audio.NewContext(44100)
		var err error
		emu.wav, err = os.ReadFile("bleep.wav")
		if err != nil {
			log.Fatalln(err)
		}
	}

	ebiten.SetWindowSize(chip8.DisplayWidth*pixelSize, chip8.DisplayHeight*pixelSize)
	ebiten.SetWindowTitle("Go CHIP-8 v" + Version)
	ebiten.SetMaxTPS(ebiten.UncappedTPS)
	ebiten.SetVsyncEnabled(false)

	// Run VM processor loop in a separate go-routine, with a channel used to raise errors
	go vm.Run(emu.errorChan, delay)

	if err := ebiten.RunGame(emu); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Update is called every tick (1/60 [s] by default).
func (e *chip8Emulator) Update() error {
	// Read the keyboard
	readKeyboard(e.vm)

	// Play sound
	e.playSound()

	// This is a *non-blocking* check for any errors on the channel
	select {
	case runtimeError := <-e.errorChan:
		// Try to see if we got a SystemError
		se, isSystemError := runtimeError.(chip8.SystemError)
		// Default code
		code := 50
		if isSystemError {
			code = se.Code()
		}
		log.Printf("Unrecoverable system error: %s", runtimeError.Error())
		os.Exit(code)
	default:
		// Noop
		_ = true
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

	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)
}

// Layout can control scaling
func (e *chip8Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
