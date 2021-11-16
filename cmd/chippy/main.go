package main

import (
	"chippy/pkg/chip8"
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
		chip8.DISPLAY_WIDTH*chip8.DISPLAY_MODIFIER, chip8.DISPLAY_HEIGHT*chip8.DISPLAY_MODIFIER,
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Create SDL2 renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// Initilaize CHIP-8 and load ROM :3
	chippy := chip8.Init()
	size, err := chippy.LoadROM("./roms/ibm_logo.ch8")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d bytes! <3\n", size)

	// Emulator loop
	emulating := true
	for emulating {
		// CHIP-8 CPU Cycle (Fetch/Decode/Execute)
		chippy.Cycle()

		// TODO: CHIP-8 Drawing -- Utilize SDL2 renderer to draw CHIP-8 screen
		// NOTE: We should only draw if required.

		// Testing of rendering -- this should be removed later
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.Present()

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
