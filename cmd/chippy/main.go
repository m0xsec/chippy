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
	"chippy/pkg/debug"
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

	// Initialize our Debug Overlay
	debug.Init()

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

		// Update debug overlay
		overlay := debug.RenderOverlay(&chippy, renderer)

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

		// Copy debug overlay to renderer
		renderer.Copy(overlay, nil, &sdl.Rect{X: 0, Y: 0, W: 100, H: 100})

		// Render to screen <3
		renderer.Present()

		// Event handling
		// TODO: Add CHIP-8 keyboard support, per spec
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				emulating = false
			}
		}

		// Maintain CHIP-8 Clock Speed, this is 60hz by default
		// TODO: Make clock speed adjustable
		sdl.Delay(1000 / chippy.ClockSpeed())
	}
}
