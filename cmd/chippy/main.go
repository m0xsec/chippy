package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("henlo from chippy <3")

	// Initialize SDL2
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Create SDL2 window
	window, err := sdl.CreateWindow("chippy <3", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// TODO: Create SDL2 renderer

	// Emulator loop
	emulating := true
	for emulating {
		// TODO: CHIP-8 Fetch/decode/execution loop

		// TODO: CHIP-8 Drawing -- Utilize SDL2 renderer

		// Event handling
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				emulating = false
			}
		}

		// CHIP-8 Updates display at 60Hz, which is (1000/60)ms
		sdl.Delay(1000 / 60)
	}
}
