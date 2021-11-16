package main

/*



                   hhhhhhh               iiii
                   h:::::h              i::::i
                   h:::::h               iiii
                   h:::::h
    cccccccccccccccch::::h hhhhh       iiiiiiippppp   ppppppppp   ppppp   pppppppppyyyyyyy           yyyyyyy
  cc:::::::::::::::ch::::hh:::::hhh    i:::::ip::::ppp:::::::::p  p::::ppp:::::::::py:::::y         y:::::y
 c:::::::::::::::::ch::::::::::::::hh   i::::ip:::::::::::::::::p p:::::::::::::::::py:::::y       y:::::y
c:::::::cccccc:::::ch:::::::hhh::::::h  i::::ipp::::::ppppp::::::ppp::::::ppppp::::::py:::::y     y:::::y
c::::::c     ccccccch::::::h   h::::::h i::::i p:::::p     p:::::p p:::::p     p:::::p y:::::y   y:::::y
c:::::c             h:::::h     h:::::h i::::i p:::::p     p:::::p p:::::p     p:::::p  y:::::y y:::::y
c:::::c             h:::::h     h:::::h i::::i p:::::p     p:::::p p:::::p     p:::::p   y:::::y:::::y
c::::::c     ccccccch:::::h     h:::::h i::::i p:::::p    p::::::p p:::::p    p::::::p    y:::::::::y
c:::::::cccccc:::::ch:::::h     h:::::hi::::::ip:::::ppppp:::::::p p:::::ppppp:::::::p     y:::::::y
 c:::::::::::::::::ch:::::h     h:::::hi::::::ip::::::::::::::::p  p::::::::::::::::p       y:::::y
  cc:::::::::::::::ch:::::h     h:::::hi::::::ip::::::::::::::pp   p::::::::::::::pp       y:::::y
    cccccccccccccccchhhhhhh     hhhhhhhiiiiiiiip::::::pppppppp     p::::::pppppppp        y:::::y
                                               p:::::p             p:::::p               y:::::y
                                               p:::::p             p:::::p              y:::::y
                                              p:::::::p           p:::::::p            y:::::y
                                              p:::::::p           p:::::::p           y:::::y
                                              p:::::::p           p:::::::p          yyyyyyy
                                              ppppppppp           ppppppppp

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
