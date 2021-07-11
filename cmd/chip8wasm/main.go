//
// CHIP-8 emulator - wasm executable
// Ben C, June 2021
// Notes: JS side must set argv to pass in arguments
//

package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/benc-uk/chip8/pkg/emulator"
)

const generalErrorCode = 1

func main() {
	log.Println("WASM emulator starting")

	if len(os.Args) != 7 {
		checkError(errors.New("wrong number of arguments"))
	}

	progURL := os.Args[0]
	log.Printf("Fetching program file '%s' via HTTP\n", progURL)
	resp, err := http.Get(progURL)
	checkError(err)
	if resp.StatusCode != 200 {
		log.Printf("Failed to download: %s", resp.Status)
		os.Exit(resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	checkError(err)
	log.Println("Program file downloaded OK")

	debug, err := strconv.Atoi(os.Args[1])
	checkError(err)
	speed, err := strconv.Atoi(os.Args[2])
	checkError(err)
	pixelSize, err := strconv.Atoi(os.Args[3])
	checkError(err)
	fg, err := strconv.Atoi(os.Args[4])
	checkError(err)
	bg, err := strconv.Atoi(os.Args[5])
	checkError(err)
	palletteName := os.Args[6]

	pallette := emulator.PalletteSpectrum
	if strings.EqualFold(palletteName, "c64") {
		pallette = emulator.PalletteC64
	}
	if strings.EqualFold(palletteName, "vaporwave") {
		pallette = emulator.PalletteVaporWave
	}

	var colourMap *emulator.ColourMap
	mapFileURL := progURL + ".colours.yaml"
	mapResp, err := http.Get(mapFileURL)
	if err == nil && mapResp.StatusCode == 200 {
		mapBody, err := ioutil.ReadAll(mapResp.Body)
		defer mapResp.Body.Close()
		checkError(err)
		log.Printf("Enabling multi-colour mode, will map colours based on: %s\n", mapFileURL)
		colourMap, _ = emulator.LoadColourMap(mapBody, pallette)
	} else {
		log.Println("Basic 1-bit colour mode enabled")
		colourMap = emulator.SimpleColourMap(fg, bg, pallette)
	}

	if colourMap == nil {
		checkError(errors.New("failed to load colour map"))
	}

	emulator.Start(body, debug, speed, pixelSize, colourMap)
}

func checkError(err error) {
	if err == nil {
		return
	}

	log.Fatalln(err)
	os.Exit(generalErrorCode)
}
