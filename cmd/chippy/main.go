package main

/*

         dP       oo
         88
.d8888b. 88d888b. dP 88d888b. 88d888b. dP    dP
88'  `"" 88'  `88 88 88'  `88 88'  `88 88    88
88.  ... 88    88 88 88.  .88 88.  .88 88.  .88
`88888P' dP    dP dP 88Y888P' 88Y888P' `8888P88
                     88       88            .88
                     dP       dP        d8888P

				CHIP-8 Emulator
					m0x <3
*/

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

		// Render CHIP-8 Screen
		// TODO: Might need to keep track of when to draw to prevent flickering,
		//       instead of rendering every frame / cycle.
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		buff := chippy.DisplayBuffer()
		for h := 0; h < len(buff); h++ {
			for w := 0; w < len(buff[h]); w++ {
				// CHIP-8 pixels are colored based on 1 or 0
				if buff[h][w] != 0 {
					renderer.SetDrawColor(255, 255, 255, 255)
				} else {
					renderer.SetDrawColor(0, 0, 0, 255)
				}

				// Render, keeping our display scaling in mind
				renderer.FillRect(&sdl.Rect{
					Y: int32(h) * chip8.DISPLAY_MODIFIER,
					X: int32(w) * chip8.DISPLAY_MODIFIER,
					W: chip8.DISPLAY_MODIFIER,
					H: chip8.DISPLAY_MODIFIER,
				})
			}
		}

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
