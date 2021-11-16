package chip8

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
	"fmt"
	"os"
)

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

// CHIP-8 structure that represents internal state and subsystems
type Chip8 struct {
	// CHIP-8 has 4K of memory
	memory [4096]uint8

	// CHIP-8 has a display that is 64x32
	display [DISPLAY_HEIGHT][DISPLAY_WIDTH]uint8

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

	// Current CHIP-8 Opcode
	// Keeps track of the current instruction opcode
	oc uint16
}

// Initializes the CHIP-8
func Init() Chip8 {
	fmt.Println("Initializing CHIP-8...")
	chippy := Chip8{
		// The first CHIP-8 interpreter, on the COMAC VIP, was located in RAM,
		// from address 000 to 1FF. It would expect a CHIP-8 program to be
		// loaded into memory after it, starting at address 200.
		// For chippy, we will use 000 to 1FF for our font set :)
		pc: 0x200,
	}

	// Load fontset into memory
	fmt.Println("Loading font set into memory...")
	for i := 0; i < len(fontset); i++ {
		chippy.memory[i] = fontset[i]
	}

	return chippy
}

// Loads a CHIP-8 ROM into memory from the given file
// Returns the size of the ROM, and an error if the ROM is invalid
func (c *Chip8) LoadROM(file string) (int64, error) {
	// Read ROM File
	rom, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return -1, err
	}
	defer rom.Close()

	// Make sure the ROM is the correct size, given that we
	// load into memory starting at 0x200
	stat, err := rom.Stat()
	if err != nil {
		return -1, err
	}
	if int64(len(c.memory)-0x200) < stat.Size() {
		return -1, fmt.Errorf("ROM is too large to fit in memory :(")
	}

	// Read the ROM into memory, starting at 0x200
	fmt.Println("Loading ROM into memory...")
	buffer := make([]byte, stat.Size())
	if _, err := rom.Read(buffer); err != nil {
		return -1, err
	}
	for i := 0; i < len(buffer); i++ {
		c.memory[i+0x200] = buffer[i]
	}

	return stat.Size(), nil
}

// Returns the current CHIP-8 Display Buffer
func (c *Chip8) DisplayBuffer() [DISPLAY_HEIGHT][DISPLAY_WIDTH]uint8 {
	return c.display
}

// Cycle the CHIP-8 CPU (Fetch, Decode, Execute)
func (c *Chip8) Cycle() {
	// Fetch Opcode (2 bytes), and merge into a single 16-bit value
	// Todo this we shift left by 8 bytes and use bitwise OR to merge
	// For example:
	// memory[pc] = 0xA2
	// memory[pc+1] = 0xF0
	// Resulting merge: 0xA2F0
	c.oc = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])

	// Decode & Execute Opcode
	switch c.oc & 0xF000 {
	// Instrucutions starting with 0x0
	// 0x00E0 - Clear the display
	// 0x00EE - Return from a subroutine
	case 0x0000:
		// Need to compare the last 4 bits of the opcode
		switch c.oc & 0x000F {
		case 0x0000: // 0x00E0 Clear the display
			// TODO: Implement 0x00E0

		case 0x000E: // 0x00EE Return from a subroutine
			// TODO: Implement 0x00EE

		default:
			fmt.Printf("[0x0000] Unknown opcode: 0x%X\n", c.oc)
		}

	case 0x1000: // 0x1NNN Jump to address NNN
		// TODO: Implement 0x1NNN

	case 0x6000: // 0x6XNN Set Register VX to NN
		// TODO: Implement 0x6XNN

	case 0x7000: // 0x7XNN Add NN to Register VX
		// TODO: Implement 0x7XNN

	case 0xA000: // 0xANNN Set I to NNN
		// TODO: Implement 0xANNN

	case 0xD000: // 0xDXYN Display (Drawing)
		// TODO: Implement 0xDXYN

	// TODO: Implement remaining CHIP-8 insturctions

	default:
		fmt.Printf("Unknown opcode: 0x%X\n", c.oc)
	}

	// TODO: Handle Delay Timer

	// TODO: Handle Sound Timer w/ Beeping
}
