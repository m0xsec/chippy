package debug

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
	"github.com/veandco/go-sdl2/ttf"
)

// Init initializes the CHIP-8 debug overlay
func Init() {
	// Initialize SDL2 TTF
	fmt.Println("Initializing SDL2 TTF...")
	err := ttf.Init()
	if err != nil {
		fmt.Println("Failed to initialize TTF: " + err.Error())
	}
}

// RenderOverlay renders the CHIP-8 debug overlay, returns an SDL2 texture
func RenderOverlay(chippy *chip8.Chip8, renderer *sdl.Renderer) *sdl.Texture {
	// Load font
	font, err := ttf.OpenFont("./fonts/VT323.ttf", 20)
	if err != nil {
		fmt.Println("Failed to load font: " + err.Error())
	}
	defer font.Close()

	// Create debug overlay surface
	width := 100
	height := 100
	overlay, err := sdl.CreateRGBSurface(0, int32(width), int32(height), 32, 0, 0, 0, 50)
	if err != nil {
		fmt.Println("Failed to create debug overlay: " + err.Error())
	}

	// Pull current CHIP-8 state
	opcode := fmt.Sprintf("OP [0x%X]", chippy.Opcode())
	pc := fmt.Sprintf("PC [0x%X]", chippy.PC())
	i := fmt.Sprintf("I [0x%X]", chippy.I())

	// Render text for each line
	opcodeSurface, err := font.RenderUTF8Blended(opcode, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Println("Failed to render surface: " + err.Error())
	}
	pcSurface, err := font.RenderUTF8Blended(pc, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Println("Failed to render surface: " + err.Error())
	}
	iSurface, err := font.RenderUTF8Blended(i, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Println("Failed to render surface: " + err.Error())
	}

	// Render text to overlay for each line
	opcodeSurface.Blit(nil, overlay, &sdl.Rect{X: 0, Y: 0, W: 100, H: 100})
	pcSurface.Blit(nil, overlay, &sdl.Rect{X: 0, Y: 20, W: 100, H: 100})
	iSurface.Blit(nil, overlay, &sdl.Rect{X: 0, Y: 40, W: 100, H: 100})

	// Create SDL2 texture from overlay
	texture, err := renderer.CreateTextureFromSurface(overlay)
	if err != nil {
		fmt.Println("Failed to create texture: " + err.Error())
	}

	return texture
}
