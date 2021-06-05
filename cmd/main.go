package main

import (
	"fmt"
	"os"

	"github.com/benc-uk/chip8/pkg/chip8"
	"github.com/hajimehoshi/ebiten"
)

var version = "0.0.1"

const displayWidth = 800
const displayHeight = 600

func main() {
	_ = chip8.NewVM()

	if err := ebiten.Run(updateDisplay, displayWidth, displayHeight, 1, "CHIP-8 v"+version); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func updateDisplay(screen *ebiten.Image) error {
	if !needsUpdate {
		return nil
	}
}
