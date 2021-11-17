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

	// CHIP-8 Stack Pointer
	sp uint16

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

	// Current CHIP-8 Clock Speed (Hz)
	// Default is 60Hz
	clockSpeed uint32

	//TODO: Need to add some way of handling keypad events
}

// Initializes the CHIP-8
func Init() Chip8 {
	fmt.Println("Initializing CHIP-8...")
	chippy := Chip8{
		// The first CHIP-8 interpreter, on the COMAC VIP, was located in RAM,
		// from address 000 to 1FF. It would expect a CHIP-8 program to be
		// loaded into memory after it, starting at address 200.
		// For chippy, we will use 000 to 1FF for our font set :)
		pc:         0x200,
		oc:         0x0,
		i:          0x0,
		sp:         0x0,
		dt:         0x0,
		st:         0x0,
		clockSpeed: 60,
	}

	// Load fontset into memory
	fmt.Println("Loading font set into memory...")
	for i := 0; i < len(fontset); i++ {
		chippy.memory[i] = fontset[i]
	}

	return chippy
}

// Returns the current CHIP-8 Clock Speed
func (chippy *Chip8) ClockSpeed() uint32 {
	return chippy.clockSpeed
}

// Returns the current CHIP-8 Display Buffer
func (c *Chip8) DisplayBuffer() [DISPLAY_HEIGHT][DISPLAY_WIDTH]uint8 {
	return c.display
}

// Returns the current CHIP-8 Opcode
func (c *Chip8) Opcode() uint16 {
	return c.oc
}

// Returns the current CHIP-8 Program Counter
func (c *Chip8) PC() uint16 {
	return c.pc
}

// Returns the current CHIP-8 Index Register
func (c *Chip8) I() uint16 {
	return c.i
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

// Cycle the CHIP-8 CPU (Fetch, Decode, Execute)
func (c *Chip8) Cycle() {
	// Fetch Opcode (2 bytes), and merge into a single 16-bit value
	// Todo this we shift left by 8 bytes and use bitwise OR to merge
	// For example:
	// memory[pc] = 0xA2
	// memory[pc+1] = 0xF0
	// Resulting merge: 0xA2F0
	c.oc = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	fmt.Printf("[0x%X]\n", c.oc)

	// Decode & Execute Opcode
	// Ex: 0xA2F0 & 0xF000 -> 0xA000
	switch c.oc & 0xF000 {
	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x0
	// 0x00E0 - Clear the display
	// 0x00EE - Return from a subroutine
	// NOTE: Did not implement 0x0NNN - Used for running machine language outside of CHIP-8
	case 0x0000:
		// Need to compare the last 4 bits of the opcode
		// Ex: 0x000E & 0x000F-> 0x000E
		switch c.oc & 0x000F {
		case 0x0000: // 0x00E0 Clear the display
			for h := 0; h < len(c.display); h++ {
				for w := 0; w < len(c.display[h]); w++ {
					c.display[h][w] = 0
				}
			}
			c.pc += 2

		case 0x000E: // 0x00EE Return from a subroutine
			// Decrease the stack pointer
			// Set the PC to the stored return address
			// Increment PC
			c.sp -= 1
			c.pc = c.stack[c.sp]
			c.pc += 2

		default:
			fmt.Printf("[0x0000] Unknown opcode: 0x%X\n", c.oc)
		}

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x1
	// 0x1NNN - Jump to address NNN
	case 0x1000: // 0x1NNN Jump to address NNN
		c.pc = c.oc & 0x0FFF

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x2
	// 0x2NNN - Call subroutine at NNN
	case 0x2000: // 0x2NNN Call subroutine at NNN
		// Store the current PC on the stack,
		// and increase stack pointer since we put something on the stack
		// Set the PC to the address NNN
		c.stack[c.sp] = c.pc
		c.sp += 1
		c.pc = c.oc & 0x0FFF

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x3
	// 0x3XNN - Skip next instruction if VX equals NN
	case 0x3000: // 0x3XNN Skip next instruction if VX equals NN
		if c.v[(c.oc&0x0F00)>>8] == uint8(c.oc&0x00FF) {
			c.pc += 4
		} else {
			c.pc += 2
		}

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x4
	// 0x4XNN - Skip next instruction if VX doesn't equal NN
	case 0x4000: //Skip next instruction if VX doesn't equal NN
		if c.v[(c.oc&0x0F00)>>8] != uint8(c.oc&0x0FF) {
			c.pc += 4
		} else {
			c.pc += 2
		}

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x5
	// 0x5XY0 - Skip next instruction if VX equals VY
	case 0x5000:
		if c.v[(c.oc&0x0F00)>>8] == c.v[(c.oc&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x6
	// 0x6XNN - Set Register VX to NN
	case 0x6000: // 0x6XNN Set Register VX to NN
		// Need to shift right 8 bytes to get X
		c.v[(c.oc&0x0F00)>>8] = uint8(c.oc & 0x00FF)
		c.pc += 2

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x7
	// 0x7XNN - Add NN to Register VX
	case 0x7000: // 0x7XNN Add NN to Register VX
		c.v[(c.oc&0x0F00)>>8] += uint8(c.oc & 0x00FF)
		c.pc += 2

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x8
	// 0x8XY0 - Set Register VX to Register VY
	// 0x8XY1 - Set Register VX to Register VX OR Register VY
	// 0x8XY2 - Set Register VX to Register VX AND Register VY
	// 0x8XY3 - Set Register VX to Register VX XOR Register VY
	// 0x8XY4 - Add Register VY to Register VX, set VF to 1 if carry, 0 if not
	// 0x8XY5 - Subtract Register VY from Register VX, set VF to 0 if borrow, 1 if not
	// 0x8XY6 - Set VX to VY. Store the least significant bit of Register VX in VF, and then shift Register VX right by 1
	// 0x8XY7 - Set Register VX to Register VY minus Register VX, set VF to 0 if borrow, 1 if not
	// 0x8XYE - Set VX to VY. Store the most significant bit of Register VX in VF, and then shift Register VX left by 1
	case 0x8000:
		switch c.oc & 0x000F {
		case 0x0000: // 0x8XY0 - Set Register VX to Register VY
			c.v[(c.oc&0x0F00)>>8] = c.v[(c.oc&0x00F0)>>4]
			c.pc += 2

		case 0x0001: // 0x8XY1 - Set Register VX to Register VX OR Register VY
			c.v[(c.oc&0x0F00)>>8] = (c.v[(c.oc&0x0F00)>>8] | c.v[(c.oc&0x00F0)>>4])
			c.pc += 2

		case 0x0002: // 0x8XY2 - Set Register VX to Register VX AND Register VY
			c.v[(c.oc&0x0F00)>>8] = (c.v[(c.oc&0x0F00)>>8] & c.v[(c.oc&0x00F0)>>4])
			c.pc += 2

		case 0x0003: // 0x8XY3 - Set Register VX to Register VX XOR Register VY
			c.v[(c.oc&0x0F00)>>8] = (c.v[(c.oc&0x0F00)>>8] ^ c.v[(c.oc&0x00F0)>>4])
			c.pc += 2

		case 0x0004: // 0x8XY4 - Add Register VY to Register VX, set VF to 1 if carry, 0 if not
			// Do we need to carry?
			if (0xFF - c.v[(c.oc&0x0F00)>>8]) < c.v[(c.oc&0x00F0)>>4] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[(c.oc&0x0F00)>>8] += c.v[(c.oc&0x00F0)>>4]
			c.pc += 2

		case 0x0005: // 0x8XY5 - Subtract Register VY from Register VX, set VF to 0 if borrow, 1 if not
			// If the first operand is larger than the second operand, set borrow flag
			if c.v[(c.oc&0x0F00)>>8] > c.v[(c.oc&0x00F0)>>4] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[(c.oc&0x0F00)>>8] -= c.v[(c.oc&0x00F0)>>4]
			c.pc += 2

		case 0x0006: // 0x8XY6 - Set VX to VY. Store the least significant bit of Register VX in VF, and then shift Register VX right by 1
			// Set VX to VY
			c.v[(c.oc&0x0F00)>>8] = c.v[(c.oc&0x00F0)>>4]

			// Store least signigicant bit of VX in VF
			c.v[0xF] = c.v[(c.oc&0x0F00)>>8] & 0x1

			// Shift VX right by 1
			c.v[(c.oc&0x0F00)>>8] >>= 1

			c.pc += 2

		case 0x0007: // 0x8XY7 - Set Register VX to Register VY minus Register VX, set VF to 0 if borrow, 1 if not
			// Do we need to borrow for VY-VX?
			if c.v[(c.oc&0x0F00)>>8] > c.v[(c.oc&0x00F0)>>4] {
				c.v[0xF] = 0 // Set to 0 if borrow
			} else {
				c.v[0xF] = 1 // Set to 1 if no borrow is needed
			}
			c.v[(c.oc&0x0F00)>>8] = c.v[(c.oc&0x00F0)>>4] - c.v[(c.oc&0x0F00)>>8]
			c.pc += 2

		case 0x000E: // 0x8XYE - Set VX to VY. Store the most significant bit of Register VX in VF, and then shift Register VX left by 1
			// Set VX to VY
			c.v[(c.oc&0x0F00)>>8] = c.v[(c.oc&0x00F0)>>4]

			// Store least signigicant bit of VX in VF
			c.v[0xF] = c.v[(c.oc&0x0F00)>>8] & 0x1

			// Shift VX left by 1
			c.v[(c.oc&0x0F00)>>8] <<= 1

			c.pc += 2

		default:
			fmt.Printf("[0x8000] Unknown opcode: 0x%X\n", c.oc)
		}

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0x9
	// 0x9XY0 - Skip next instruction if VX doesn't equal VY
	case 0x9000:
		if c.v[(c.oc&0x0F00)>>8] != c.v[(c.oc&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0xA
	// 0xANNN - Set I to NNN
	case 0xA000: // 0xANNN Set I to NNN
		c.i = c.oc & 0x0FFF
		c.pc += 2

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0xB
	// 0xBNNN - Jump to address NNN + V0
	// TODO: Implement 0xB instructions
	// ...

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0xC
	// 0xCXNN - Set VX to a random number and NN
	// TODO: Implement 0xC instructions
	// ...

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0xD
	// 0xDXYN - Draw a sprite at position VX, VY with N bytes of sprite data starting at the address
	//			stored in I. Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
	case 0xD000: // 0xDXYN Display (Drawing)
		// Fetch (X,Y) from VX and VY
		x := c.v[(c.oc&0x0F00)>>8] % uint8(DISPLAY_WIDTH)
		y := c.v[(c.oc&0x00F0)>>4] % uint8(DISPLAY_HEIGHT)

		// Fetch N from opcode (N is our height)
		n := c.oc & 0x000F

		// Reset VF Register
		c.v[0xF] = 0

		// Loop through the height (rows)
		for i := 0; i < int(n); i++ {
			// Fetch Nth byte of sprite data, at I register + i
			byte := c.memory[c.i+uint16(i)]

			// Each sprite row, there are 8 bits for each pixel
			for j := 0; j < 8; j++ {
				// Is the current pixel set?
				if byte&(0x80>>uint8(j)) != 0 {
					// Is that display pixel at X,Y also set?
					if c.display[y+uint8(i)][x+uint8(j)] == 1 {
						// Turn off the pixel, and set VF to 1
						c.display[y+uint8(i)][x+uint8(j)] = 0
						c.v[0xF] = 1
					}

					// Turn on display pixel
					c.display[y+uint8(i)][x+uint8(j)] ^= 1
				}
			}
		}

		fmt.Printf("[0xDXYN] X: %d, Y: %d, N: %d\n", x, y, n)
		c.pc += 2

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0xE
	// 0xEX9E - Skip next instruction if key stored in VX is pressed
	// 0xEXA1 - Skip next instruction if key stored in VX isn't pressed
	// TODO: Implement 0xE instructions
	// ...

	/////////////////////////////////////////////////////////////////////////////////////////
	// Instrucutions starting with 0xF
	// 0xFX07 - Set VX to the value of the delay timer
	// 0xFX0A - Wait for a key press, store the value of the key in VX
	// 0xFX15 - Set the delay timer to VX
	// 0xFX18 - Set the sound timer to VX
	// 0xFX1E - Add VX to I
	// 0xFX29 - Set I to the location of the sprite for the character in VX
	// 0xFX33 - Store the binary-coded decimal representation of VX in memory locations
	//			I, I+1, and I+2
	// 0xFX55 - Store the values of registers V0 to VX inclusive in memory starting at address I
	//			I is set to I + X + 1 after operation
	// 0xFX65 - Fill registers V0 to VX inclusive with the values stored in memory starting at address I
	//			I is set to I + X + 1 after operation
	// TODO: Implement 0xF instructions
	// ...

	default:
		fmt.Printf("Unknown opcode: 0x%X\n", c.oc)
	}

	// NOTE: Timers decrement at 60Hz, independant of the clock speed used in our cycle loop
	// TODO: Handle Delay Timer
	// TODO: Handle Sound Timer w/ Beeping
}
