//
// CHIP-8 emulator - wasm executable
// Ben C, June 2021
// Notes: JS side must set argv to pass in arguments
//

package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/benc-uk/chip8/pkg/emulator"
)

const generalErrorCode = 1

func main() {
	log.Println("WASM emulator starting")

	if len(os.Args) != 4 {
		checkError(errors.New("wrong number of arguments"))
	}

	log.Printf("Fetching program file '%s' via HTTP\n", os.Args[0])
	resp, err := http.Get(os.Args[0])
	checkError(err)
	if resp.StatusCode != 200 {
		log.Printf("Failed to download: %s", resp.Status)
		os.Exit(resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	checkError(err)
	log.Println("Program file downloaded OK")

	debug, err := strconv.ParseBool(os.Args[1])
	checkError(err)
	speed, err := strconv.Atoi(os.Args[2])
	checkError(err)
	pixelSize, err := strconv.Atoi(os.Args[3])
	checkError(err)

	emulator.Start(body, debug, speed, pixelSize)
}

func checkError(err error) {
	if err == nil {
		return
	}

	log.Fatalln(err)
	os.Exit(generalErrorCode)
}
