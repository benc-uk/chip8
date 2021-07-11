//
// CHIP-8 - Emulator wraps the CHIP-8 VM and allows us to interact with it
// Ben C, June 2021
// Notes:
//

package emulator

import (
	"crypto/md5"
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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Version is the emulator version
var Version = "1.5.0"

//var pixelColour = color.RGBA{0x00, 0xff, 0x00, 0xff}

// Wrapper for ebiten implements the ebiten.Game interface
type chip8Emulator struct {
	vm        *chip8.VM
	display   *ebiten.Image
	pixelSize int
	speed     int
	paused    bool
	pgmData   []byte // Only stored so we can do a soft reset
	showSpeed int
	colourMap *ColourMap

	audioContext *audio.Context
	bleeper      *audio.Player
}

// Start is called by the WASM and console main.go to start everything
func Start(program []byte, debugLevel int, speed int, pixelSize int, colourMap *ColourMap) {
	console.Infof("Starting CHIP-8 emulator version v%s\n\n", Version)

	if runtime.GOARCH == "js" || runtime.GOOS == "js" {
		ebiten.SetFullscreen(true)
	}

	if speed < 1 {
		log.Fatalln("Speed must be greater than 0")
	}
	if pixelSize < 1 || pixelSize > 60 {
		log.Fatalln("Pixel size must be be between 1 and 60")
	}

	// Create a new CHIP-8 virtual machine, and load program into it
	vm := chip8.NewVM(true)
	vm.DebugLevel = debugLevel

	// Load supplied data as a program
	err := vm.LoadProgram(program)
	if err != nil {
		fmt.Printf("R Tape loading error: %s\n", err)
		os.Exit(1)
	}

	// Wrap the VM in an chip8Emulator to allow us to use ebiten with it
	emu := &chip8Emulator{
		vm:           vm,
		display:      ebiten.NewImage(chip8.DisplayWidth*pixelSize, chip8.DisplayHeight*pixelSize),
		pixelSize:    pixelSize,
		speed:        speed,
		audioContext: audio.NewContext(44100),
		pgmData:      program,
		showSpeed:    0,
		colourMap:    colourMap,
	}

	ebiten.SetWindowSize(chip8.DisplayWidth*pixelSize, chip8.DisplayHeight*pixelSize)
	ebiten.SetWindowTitle("Go CHIP-8 v" + Version)
	//ebiten.SetMaxTPS(ebiten.UncappedTPS)
	ebiten.SetVsyncEnabled(false)

	// Create a new audio player	and set it to play the program

	console.Successf("Program MD5: %X\n", md5.Sum(program))

	console.Successf("Starting CHIP-8 system, processor at address 0x%04X\n", 0x200)
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
		// Handle pausing and stepping through code, kinda funky but it works
		if e.paused && !inpututil.IsKeyJustPressed(ebiten.KeyF6) {
			c--
			return nil
		}

		// This advances the processor one tick/cycle
		runtimeError := e.vm.Cycle()
		e.checkErr(runtimeError)
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
	if e.vm.DebugEnabled() {
		debugMsg += "\nDEBUGGING"
	}
	if e.showSpeed > 0 {
		debugMsg += fmt.Sprintf("\nSPEED: %d", e.speed)
		e.showSpeed--
	}
	ebitenutil.DebugPrint(screen, debugMsg)
}

// Layout can control scaling
func (e *chip8Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (e *chip8Emulator) SoftReset() {
	e.vm.Reset()
	e.vm.LoadProgram(e.pgmData)
}

// Handle runtime errors which are always fatal
func (e *chip8Emulator) checkErr(err error) {
	if err != nil {
		se, isSystemError := err.(chip8.SystemError)
		code := 50 // Default code
		if isSystemError {
			code = se.Code()
		}

		log.Printf("Unrecoverable system error: %s", err.Error())
		e.vm.Dump()
		os.Exit(code)
	}
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	return
}
