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
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	fmt.Println("henlo from chippy <3")

	// Initialize SDL2
	fmt.Println("Initializing SDL2...")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Initialize SDL2 TTF
	fmt.Println("Initializing SDL2 TTF...")
	err := ttf.Init()
	if err != nil {
		fmt.Println("Failed to initialize TTF: " + err.Error())
	}

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
	size, err := chippy.LoadROM("./roms/space_invaders.ch8")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d bytes! <3\n", size)

	// Emulator loop
	emulating := true
	displayOverlay := true
	for emulating {
		// CHIP-8 CPU Cycle (Fetch/Decode/Execute)
		chippy.Cycle()

		// Update debug overlay
		var overlay *sdl.Texture
		if displayOverlay {
			overlay = debug.RenderOverlay(&chippy, renderer)
		}

		// Render CHIP-8 Screen
		// NOTE: Might need to keep track of when to draw to prevent flickering,
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
		if displayOverlay {
			renderer.Copy(overlay, nil, &sdl.Rect{X: 0, Y: 0, W: 100, H: 100})
		}

		// Render to screen <3
		renderer.Present()

		// Event handling
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("kthxbai<3")
				emulating = false

			case *sdl.KeyboardEvent:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					println("kthxbai<3")
					emulating = false

				case sdl.K_LALT:
					if t.State == sdl.PRESSED {
						displayOverlay = !displayOverlay
					}

				case chippy.KeyMap()[0x0]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x0)
					} else {
						chippy.KeyRelease(0x0)
					}

				case chippy.KeyMap()[0x1]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x1)
					} else {
						chippy.KeyRelease(0x1)
					}

				case chippy.KeyMap()[0x2]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x2)
					} else {
						chippy.KeyRelease(0x2)
					}

				case chippy.KeyMap()[0x3]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x3)
					} else {
						chippy.KeyRelease(0x3)
					}

				case chippy.KeyMap()[0x4]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x4)
					} else {
						chippy.KeyRelease(0x4)
					}

				case chippy.KeyMap()[0x5]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x5)
					} else {
						chippy.KeyRelease(0x5)
					}

				case chippy.KeyMap()[0x6]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x6)
					} else {
						chippy.KeyRelease(0x6)
					}

				case chippy.KeyMap()[0x7]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x7)
					} else {
						chippy.KeyRelease(0x7)
					}

				case chippy.KeyMap()[0x8]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x8)
					} else {
						chippy.KeyRelease(0x8)
					}

				case chippy.KeyMap()[0x9]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0x9)
					} else {
						chippy.KeyRelease(0x9)
					}

				case chippy.KeyMap()[0xA]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0xA)
					} else {
						chippy.KeyRelease(0xA)
					}

				case chippy.KeyMap()[0xB]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0xB)
					} else {
						chippy.KeyRelease(0xB)
					}

				case chippy.KeyMap()[0xC]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0xC)
					} else {
						chippy.KeyRelease(0xC)
					}

				case chippy.KeyMap()[0xD]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0xD)
					} else {
						chippy.KeyRelease(0xD)
					}

				case chippy.KeyMap()[0xE]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0xE)
					} else {
						chippy.KeyRelease(0xE)
					}

				case chippy.KeyMap()[0xF]:
					if t.State == sdl.PRESSED {
						chippy.KeyPress(0xF)
					} else {
						chippy.KeyRelease(0xF)
					}
				}
			}
		}

		// Maintain CHIP-8 Clock Speed, this is ~500hz by default
		// TODO: Make clock speed adjustable
		// TODO: Make sure display is updating at 60hz (60 FPS)?
		sdl.Delay(1000 / chippy.ClockSpeed())
	}
}
