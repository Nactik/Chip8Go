package core

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

var sprites = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //O
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type Chip8 struct {
	/** Registers **/
	vx         [16]uint8 // General purpose registers referred to as Vx, where x is a hexadecimal digit (0 through F)
	i          uint16    // Store memory addresses
	delayTimer uint8     // Automatically decremented at a rate of 60Hz
	soundTimer uint8     // Automatically decremented at a rate of 60Hz - Buzz when not non-zero

	/** Pseudo-registers **/
	pc uint16 //Program counter - Store the current executing address
	sp uint8  //Stack pointer - Point to the topmost level of the stack

	stack [16]uint16 // store the address that the interpreter shoud return to when finished with a subroutine

	/** Other **/
	display [32][64]uint8 //64x32 monochrome display (should be an array of bits but the minimum adressable memory is a byte)
	memory  [4096]uint8   //Memory of 4096 bytes, from Ox000 to OxFFF (4095)
}

func Init() Chip8 {
	chip8 := Chip8{
		pc: 0x200, //pc starts at 0x200 for most of programs
	}

	for i := 0; i < len(sprites); i++ {
		chip8.memory[i] = sprites[i]
	}

	return chip8
}

func (c *Chip8) LoadRom(romPath string) error {
	rom, err := os.Open(romPath)

	if err != nil {
		return err
	}

	defer rom.Close()

	stats, statsErr := rom.Stat()
	if statsErr != nil {
		return statsErr
	}

	romSize := stats.Size()
	memoryCap := cap(c.memory[c.pc:])

	// Check if can be put in memory
	if romSize > int64(memoryCap) {
		fmt.Println("Rom is too big")
		return fmt.Errorf("Rom is too big")
	}

	// Write rom content into memory
	bufr := bufio.NewReader(rom)
	_, readErr := bufr.Read(c.memory[c.pc:])

	if readErr != nil {
		return readErr
	} else {
		return nil
	}
}

// TODO
func (c *Chip8) readInstruction(opcode uint16) error {
	switch opcode & 0xF000 {
	case 0x0000:

	case 0x1000:
		// JMP Instruction 0x1nnn - Jump to adress nnn
		nextPc := opcode & 0x0FFF
		c.pc = nextPc
	case 0x2000:
		// Call subroutine at nnn.
		// The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
		c.sp++
		c.stack[c.sp] = c.pc
		c.pc = opcode & 0x0FFF
	case 0x3000:
		// Skip next instruction if Vx = kk.
		// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
		// Skip next instruction if Vx = kk.
		// The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
		kk := uint8(opcode & 0x00FF)
		x := (opcode & 0x0F00) >> 8
		if c.vx[x] == kk {
			c.pc += 2
		}
	case 0x4000:
		// Skip next instruction if Vx != kk.
		// The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
		kk := uint8(opcode & 0x00FF)
		x := (opcode & 0x0F00) >> 8
		if c.vx[x] != kk {
			c.pc += 2
		}
	case 0x5000:
		// 0x5xy0
		// Skip next instruction if Vx = Vy.
		// The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4
		if c.vx[x] != c.vx[y] {
			c.pc += 2
		}
	case 0x6000:
		// Ox6xkk
		// Set Vx = kk.
		// The interpreter puts the value kk into register Vx.
		x := (opcode & 0x0F00) >> 8
		kk := uint8(opcode & 0x00FF)
		c.vx[x] = kk
	case 0x7000:
		// 0x7xkk
		// Set Vx = Vx + kk.
		// Adds the value kk to the value of register Vx, then stores the result in Vx.
		x := (opcode & 0x0F00) >> 8
		kk := uint8(opcode & 0x00FF)
		c.vx[x] = c.vx[x] + kk
	case 0x8000:
	case 0x9000:
		// 0x9xy0
		// Skip next instruction if Vx != Vy.
		// The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
		x := (opcode & 0x0F00) >> 8
		y := (opcode & 0x00F0) >> 4

		if c.vx[x] != c.vx[y] {
			c.pc += 2
		}
	case 0xA000:
		// OxAnnn
		// Set I = nnn.
		// The value of register I is set to nnn.
		nnn := opcode & 0x0FFF
		c.i = nnn
	case 0xB000:
		// 0xBnnn
		// Jump to location nnn + V0.
		// The program counter is set to nnn plus the value of V0.
		nnn := opcode & 0x0FFF
		v0 := c.vx[0]
		c.pc = nnn + uint16(v0)
	case 0xC000:
		// Cxkk - RND Vx, byte
		// Set Vx = random byte AND kk.
		// The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx
		x := (opcode & 0x0F00) >> 8
		kk := uint8(opcode & 0x00FF)
		randomValue := uint8(rand.Intn(256)) & kk
		c.vx[x] = randomValue
	case 0xD000:
	case 0xE000:
	case 0xF000:
	default:
		return fmt.Errorf("Unknown Opcode")
	}
}
