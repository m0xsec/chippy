package chip8

// CHIP-8 Display Width 64px
const DISPLAY_WIDTH int32 = 64

// CHIP-8 Display Height 32px
const DISPLAY_HEIGHT int32 = 32

// CHIP-8 Display Scaling Factor
const DISPLAY_MODIFIER int32 = 10

// CHIP-8
type Chip8 struct {
	// CHIP-8 has 4K of memory
	memory [4096]uint8

	// CHIP-8 has a display that is 64x32
	display [DISPLAY_WIDTH * DISPLAY_HEIGHT]uint8

	// CHIP-8 Program Counter
	// Points at the current instruction in memory
	pc uint16

	// CHIP-8 Index Register
	// Used to point at locations in memory
	i uint16

	// CHIP-8 Stack
	// Follows LIFO, used to call and return from subroutines
	stack [16]uint16

	// CHIP-8 Delay Timer
	// Decremented at 60Hz until it reaches 0
	dt uint8

	// CHIP-8 Sound Timer
	// Decremented at 60Hz until it reaches 0
	// Makes a beeping sound as long as it is not 0
	st uint8

	// CHIP-8 Registers
	// 16 general purpose variable registers
	// Each register is 8 bits
	// Called V0, V1, V2, V3, V4, V5, V6, V7, V8, V9, VA, VB, VC, VD, VE, VF
	// VF may be used as a flag register
	v [16]uint8
}
