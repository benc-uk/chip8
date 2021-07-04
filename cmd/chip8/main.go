//
// CHIP-8 emulator - main executable
// Ben C, June 2021
// Notes:
//

package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/benc-uk/chip8/pkg/console"
	"github.com/benc-uk/chip8/pkg/emulator"
)

const banner = `
 ██████╗  ██████╗      ██████╗██╗  ██╗██╗██████╗        █████╗ 
██╔════╝ ██╔═══██╗    ██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
██║  ███╗██║   ██║    ██║     ███████║██║██████╔╝█████╗╚█████╔╝
██║   ██║██║   ██║    ██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
╚██████╔╝╚██████╔╝    ╚██████╗██║  ██║██║██║           ╚█████╔╝
 ╚═════╝  ╚═════╝      ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝`

func main() {
	var debugLevelFlag = flag.Int("debug", 0, "Debug output level: 0 = off, 1 = sprites only, 2 = full")
	var speedFlag = flag.Int("speed", 12, "Speed of the emulator in cycles per tick")
	var scaleFlag = flag.Int("scale", 10, "Size of pixels, default results in a 640x320 window")
	var fgFlag = flag.Int("fg", 2, "Colour of foreground pixels, pallette index: 0-8")
	var bgFlag = flag.Int("bg", 0, "Colour of background, pallette index: 0-8")
	var palletteFlag = flag.String("pallette", "spectrum", "Colour pallette; spectrum, c64 or vaporwave")
	var monoFlag = flag.Bool("nocolour", false, "Force mono mode even if a colour map file is found")
	flag.Parse()

	if len(flag.Args()) < 1 {
		console.Error("💥 Please supply filename of program/ROM to load")
		os.Exit(1)
	}
	progFile := flag.Arg(0)

	console.Info(banner)

	console.Infof("💾 Loading program from disk %s\n", progFile)
	pgmBytes, err := ioutil.ReadFile(progFile)
	if err != nil {
		console.Errorf("💣 Unable to load file %s\n", progFile)
		os.Exit(1)
	}

	if *fgFlag < 0 || *fgFlag > 8 {
		console.Errorf("💣 Invalid foreground colour index %d\n", *fgFlag)
		os.Exit(1)
	}
	if *bgFlag < 0 || *bgFlag > 8 {
		console.Errorf("💣 Invalid background colour index %d\n", *bgFlag)
		os.Exit(1)
	}

	pallette := emulator.PalletteSpectrum
	if strings.EqualFold(*palletteFlag, "c64") {
		pallette = emulator.PalletteC64
	}
	if strings.EqualFold(*palletteFlag, "vaporwave") {
		pallette = emulator.PalletteVaporWave
	}

	var colourMap *emulator.ColourMap
	path := filepath.Dir(progFile)
	romName := filepath.Base(progFile)
	mapFilePath := filepath.Join(path, romName+".colours.yaml")
	yamlRaw, err := ioutil.ReadFile(mapFilePath)

	if err == nil && !*monoFlag {
		console.Infof("Enabling multi-colour mode, will map colours based on: %s\n", mapFilePath)
		colourMap, _ = emulator.LoadColourMap(yamlRaw, pallette)
	} else {
		console.Info("Basic 1-bit colour mode enabled")
		colourMap = emulator.SimpleColourMap(*fgFlag, *bgFlag, pallette)
	}

	if colourMap == nil {
		console.Errorf("💣 Colour map error!\n", progFile)
		os.Exit(1)
	}

	emulator.Start(pgmBytes, *debugLevelFlag, *speedFlag, *scaleFlag, colourMap)
}
