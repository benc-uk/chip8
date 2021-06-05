package main

import (
	"github.com/benc-uk/chip8/pkg/chip8"
)

func main() {
	sys := chip8.NewSystem()
	sys.Run()
}
