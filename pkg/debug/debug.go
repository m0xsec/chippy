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
	overlay, err := sdl.CreateRGBSurface(0, int32(width), int32(height), 32, 0, 0, 0, 0)
	if err != nil {
		fmt.Println("Failed to create debug overlay: " + err.Error())
	}

	// Pull current CHIP-8 state
	opcodeStr := fmt.Sprintf("OP [0x%X]", chippy.Opcode())
	pcStr := fmt.Sprintf("PC [0x%X]", chippy.PC())
	iStr := fmt.Sprintf("I [0x%X]", chippy.I())

	// Render text for each line
	opcode, err := font.RenderUTF8Blended(opcodeStr, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Println("Failed to render surface: " + err.Error())
	}
	pc, err := font.RenderUTF8Blended(pcStr, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Println("Failed to render surface: " + err.Error())
	}
	i, err := font.RenderUTF8Blended(iStr, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Println("Failed to render surface: " + err.Error())
	}

	// Render text to overlay for each line
	opcode.Blit(nil, overlay, &sdl.Rect{X: 0, Y: 0, W: 100, H: 100})
	pc.Blit(nil, overlay, &sdl.Rect{X: 0, Y: 20, W: 100, H: 100})
	i.Blit(nil, overlay, &sdl.Rect{X: 0, Y: 40, W: 100, H: 100})

	// Create SDL2 texture from overlay
	texture, err := renderer.CreateTextureFromSurface(overlay)
	if err != nil {
		fmt.Println("Failed to create texture: " + err.Error())
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	texture.SetAlphaMod(200)

	return texture
}
