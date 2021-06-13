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

func main() {
	log.Println("WASM emulator starting")

	log.Println(os.Args)
	if len(os.Args) != 4 {
		checkError(errors.New("wrong number of arguments"))
	}

	resp, err := http.Get(os.Args[0])
	checkError(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	checkError(err)

	// pgmBase64 := os.Args[0]
	// pgmBytes, _ := base64.StdEncoding.DecodeString(pgmBase64)
	debug, err := strconv.ParseBool(os.Args[1])
	checkError(err)
	delay, err := strconv.Atoi(os.Args[2])
	checkError(err)
	pixelSize, err := strconv.Atoi(os.Args[3])
	checkError(err)

	emulator.Start(body, debug, delay, pixelSize)
}

func checkError(err error) {
	if err == nil {
		return
	}

	log.Fatalln(err)
	os.Exit(1)
}
