package chip8

import "fmt"

// CHIP-8 Display Width 64px
const DISPLAY_WIDTH int32 = 64

// CHIP-8 Display Height 32px
const DISPLAY_HEIGHT int32 = 32

// CHIP-8 Display Scaling Factor
const DISPLAY_MODIFIER int32 = 10

// CHIP-8 Font Set
var fontset = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// CHIP-8
type Chip8 struct {
	// CHIP-8 has 4K of memory
	Memory [4096]uint8

	// CHIP-8 has a display that is 64x32
	Display [DISPLAY_HEIGHT][DISPLAY_WIDTH]uint8

	// CHIP-8 Program Counter
	// Points at the current instruction in memory
	PC uint16

	// CHIP-8 Index Register
	// Used to point at locations in memory
	I uint16

	// CHIP-8 Stack
	// Follows LIFO, used to call and return from subroutines
	Stack [16]uint16

	// CHIP-8 Delay Timer
	// Decremented at 60Hz until it reaches 0
	DT uint8

	// CHIP-8 Sound Timer
	// Decremented at 60Hz until it reaches 0
	// Makes a beeping sound as long as it is not 0
	ST uint8

	// CHIP-8 Registers
	// 16 general purpose variable registers
	// Each register is 8 bits
	// Called V0, V1, V2, V3, V4, V5, V6, V7, V8, V9, VA, VB, VC, VD, VE, VF
	// VF may be used as a flag register
	V [16]uint8
}

// Initializes the CHIP-8
func Init() Chip8 {
	fmt.Println("Initializing CHIP-8...")
	chippy := Chip8{
		// The first CHIP-8 interpreter, on the COMAC VIP, was located in RAM,
		// from address 000 to 1FF. It would expect a CHIP-8 program to be
		// loaded into memory after it, starting at address 200.
		// For chippy, we will use 000 to 1FF for our font set :)
		PC: 0x200,
	}

	// Load fontset into memory
	fmt.Println("Loading font set into memory...")
	for i := 0; i < len(fontset); i++ {
		chippy.Memory[i] = fontset[i]
	}

	return chippy
}