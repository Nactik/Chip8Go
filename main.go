package main

import (
	"fmt"
	"os"

	"github.com/Nactik/Chip8Go/core"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Missing argument : Rom path")
		os.Exit(-1)
	} else if len(os.Args) > 2 {
		fmt.Println("Too many arguments")
		os.Exit(-1)
	}

	romPath := os.Args[1]
	chip8 := core.Init()

	chip8.LoadRom(romPath)

}
