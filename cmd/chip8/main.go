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
	var debugFlag = flag.Bool("debug", false, "Enable debug, lots of output very slow")
	var speedFlag = flag.Int("speed", 12, "Speed of the emulator in cycles per tick")
	var scaleFlag = flag.Int("scale", 10, "Size of pixels, default results in a 640x320 window")
	var fgFlag = flag.String("fg-colour", "#22DD22", "Colour of foreground pixels")
	var bgFlag = flag.String("bg-colour", "#000000", "Colour of background")
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

	emulator.Start(pgmBytes, *debugFlag, *speedFlag, *scaleFlag, *fgFlag, *bgFlag)
}
