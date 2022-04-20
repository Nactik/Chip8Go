package main

import (
	"fmt"
	"os"

	"github.com/Nactik/Chip8Go/core"
	"github.com/veandco/go-sdl2/sdl"
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

	fmt.Println("ROM Loading !")
	err := chip8.LoadRom(romPath)

	if err != nil {
		panic(err)
	}

	fmt.Println("ROM Loaded !")
	//Launch graphical interface

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, render, err := sdl.CreateWindowAndRenderer(64, 32, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer render.Destroy()

	fmt.Println("Emulate loop")
	// Emulate loops
	for {
		chip8.DrawFlag = false
		err = chip8.Cycle()
		if err != nil {
			panic(err)
		}

		//TODO : Display the screen
		if chip8.DrawFlag {
			display := chip8.GetDisplay()
			for y := 0; y < len(display); y++ {
				for x := 0; x < len(display[y]); x++ {
					if display[y][x] == 1 {
						render.SetDrawColor(255, 255, 255, 1)
					} else {
						render.SetDrawColor(0, 0, 0, 1)
					}
					//On render
					render.DrawPoint(int32(x), int32(y))
				}
			}
			render.Present()
		}
		//TODO : Check keyboard inputs
		event := sdl.PollEvent()
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			break
		}
	}
}
