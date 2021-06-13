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
	var debugFlag = flag.Bool("debug", false, "Enable debug")
	var slowFlag = flag.Int("slow", 1, "Pause the processor for this num microseconds each cycle")
	var scaleFlag = flag.Int("scale", 5, "Scale up the size of pixels")
	flag.Parse()

	if len(flag.Args()) < 1 {
		console.Error("Please supply filename of program/ROM to load")
		os.Exit(1)
	}
	progFile := flag.Arg(0)

	console.Info(banner)

	console.Infof("Loading program from disk %s\n", progFile)
	pgmBytes, err := ioutil.ReadFile(progFile)
	if err != nil {
		console.Errorf("Unable to load file %s", progFile)
	}

	emulator.Start(pgmBytes, *debugFlag, *slowFlag, *scaleFlag)
}
