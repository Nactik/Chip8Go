package core

import (
	"bufio"
	"fmt"
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

// TODO: Open Rom file and copy content to memory
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
