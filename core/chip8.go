package chip8

type Chip8 struct {
	/** Registers **/
	vx         [16]uint8 // General purpose registers referred to as Vx, where x is a hexadecimal digit (0 through F)
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
