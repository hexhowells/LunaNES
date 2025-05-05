package emu

import (
	"fmt"
)

type AddrModeType int

const (
	AddrModeIMP AddrModeType = iota
	AddrModeIMM
	AddrModeZP0
	AddrModeZPX
	AddrModeZPY
	AddrModeIZX
	AddrModeIZY
	AddrModeABS
	AddrModeABX
	AddrModeABY
	AddrModeIND
	AddrModeREL
)


type CPU struct {
	bus *Bus
	A uint8  // Accumulator register
	X uint8  // X register
	Y uint8  // Y register
	Stkp uint8  // Stack pointer
	Pc uint16  // Program counter
	Status uint8  // Status register
	fetched uint8  // The working input value to the ALU
	addr_abs uint16
	addr_rel uint16  // Relative address - used for jump instructions
	opcode uint8
	cycles uint8  // How many cycles the instruction has left
	clock_count uint32
	lookup []INSTRUCTION
}


type INSTRUCTION struct {
	Name string
	Operate func(*CPU) uint8
	AddrMode func(*CPU) uint8
	ModeType AddrModeType
	Cycles uint8
}


const (
	C = 1 << 0  // Carry bit
	Z = 1 << 1  // Zero
	I = 1 << 2  // Disable interrupts
	D = 1 << 3  // Decimal mode
	B = 1 << 4  // Break
	U = 1 << 5  // Unused
	V = 1 << 6  // Overflow
	N = 1 << 7  // Negative
)


func NewCPU() *CPU {
	cpu := CPU{}
	cpu.A = 0x00
	cpu.X = 0x00
	cpu.Y = 0x00
	cpu.Stkp = 0x00
	cpu.Pc = 0x0000
	cpu.Status = 0x00

	cpu.fetched = 0x00
	cpu.addr_abs = 0x0000
	cpu.addr_rel = 0x00
	cpu.opcode = 0x00
	cpu.cycles = 0

	cpu.lookup = []INSTRUCTION {
		INSTRUCTION{ "BRK", (*CPU).BRK, (*CPU).IMM, AddrModeIMM, 7 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 3 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ZP0, AddrModeZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "PHP", (*CPU).PHP, (*CPU).IMP, AddrModeIMP, 3 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ABS, AddrModeABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },
		INSTRUCTION{ "BPL", (*CPU).BPL, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).IZY, AddrModeIZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ZPX, AddrModeZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "CLC", (*CPU).CLC, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ABX, AddrModeABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },
		INSTRUCTION{ "JSR", (*CPU).JSR, (*CPU).ABS, AddrModeABS, 6 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "BIT", (*CPU).BIT, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ZP0, AddrModeZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "PLP", (*CPU).PLP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "BIT", (*CPU).BIT, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ABS, AddrModeABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },
		INSTRUCTION{ "BMI", (*CPU).BMI, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).IZY, AddrModeIZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ZPX, AddrModeZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "SEC", (*CPU).SEC, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ABX, AddrModeABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },
		INSTRUCTION{ "RTI", (*CPU).RTI, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 3 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ZP0, AddrModeZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "PHA", (*CPU).PHA, (*CPU).IMP, AddrModeIMP, 3 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "JMP", (*CPU).JMP, (*CPU).ABS, AddrModeABS, 3 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ABS, AddrModeABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },
		INSTRUCTION{ "BVC", (*CPU).BVC, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).IZY, AddrModeIZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ZPX, AddrModeZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "CLI", (*CPU).CLI, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ABX, AddrModeABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },
		INSTRUCTION{ "RTS", (*CPU).RTS, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 3 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ZP0, AddrModeZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "PLA", (*CPU).PLA, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "JMP", (*CPU).JMP, (*CPU).IND, AddrModeIND, 5 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ABS, AddrModeABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },
		INSTRUCTION{ "BVS", (*CPU).BVS, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).IZY, AddrModeIZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ZPX, AddrModeZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "SEI", (*CPU).SEI, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ABX, AddrModeABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },
		INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "STY", (*CPU).STY, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "STX", (*CPU).STX, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 3 },INSTRUCTION{ "DEY", (*CPU).DEY, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "TXA", (*CPU).TXA, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "STY", (*CPU).STY, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "STX", (*CPU).STX, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 4 },
		INSTRUCTION{ "BCC", (*CPU).BCC, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).IZY, AddrModeIZY, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "STY", (*CPU).STY, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "STX", (*CPU).STX, (*CPU).ZPY, AddrModeZPY, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "TYA", (*CPU).TYA, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ABY, AddrModeABY, 5 },INSTRUCTION{ "TXS", (*CPU).TXS, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ABX, AddrModeABX, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },
		INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 3 },INSTRUCTION{ "TAY", (*CPU).TAY, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "TAX", (*CPU).TAX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 4 },
		INSTRUCTION{ "BCS", (*CPU).BCS, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).IZY, AddrModeIZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ZPY, AddrModeZPY, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "CLV", (*CPU).CLV, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "TSX", (*CPU).TSX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 4 },
		INSTRUCTION{ "CPY", (*CPU).CPY, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "CPY", (*CPU).CPY, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ZP0, AddrModeZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "INY", (*CPU).INY, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "DEX", (*CPU).DEX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "CPY", (*CPU).CPY, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ABS, AddrModeABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },
		INSTRUCTION{ "BNE", (*CPU).BNE, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).IZY, AddrModeIZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ZPX, AddrModeZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "CLD", (*CPU).CLD, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "NOP", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ABX, AddrModeABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },
		INSTRUCTION{ "CPX", (*CPU).CPX, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).IZX, AddrModeIZX, 6 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "CPX", (*CPU).CPX, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ZP0, AddrModeZP0, 3 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ZP0, AddrModeZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 5 },INSTRUCTION{ "INX", (*CPU).INX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).IMM, AddrModeIMM, 2 },INSTRUCTION{ "NOP", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).SBC, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "CPX", (*CPU).CPX, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ABS, AddrModeABS, 4 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ABS, AddrModeABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },
		INSTRUCTION{ "BEQ", (*CPU).BEQ, (*CPU).REL, AddrModeREL, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).IZY, AddrModeIZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ZPX, AddrModeZPX, 4 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ZPX, AddrModeZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 6 },INSTRUCTION{ "SED", (*CPU).SED, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ABY, AddrModeABY, 4 },INSTRUCTION{ "NOP", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, AddrModeIMP, 4 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ABX, AddrModeABX, 4 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ABX, AddrModeABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, AddrModeIMP, 7 },
	}

	return &cpu
}


func (cpu *CPU) ConnectBus(n *Bus) {
	cpu.bus = n
}


func (cpu *CPU) Read(A uint16) uint8 {
	return cpu.bus.CpuRead(A, false)
}


func (cpu *CPU) Write(A uint16, d uint8) {
	cpu.bus.CpuWrite(A, d)
}


// Updates the clock cycle
func (cpu *CPU) Clock() bool {
	new_inst := cpu.cycles == 0

	if cpu.cycles == 0 {
		cpu.opcode = cpu.Read(cpu.Pc)
		cpu.SetFlag(U, true)
		cpu.Pc++
		cpu.cycles = cpu.lookup[cpu.opcode].Cycles

		additional_cycle1 := cpu.lookup[cpu.opcode].AddrMode(cpu)

		additional_cycle2 := cpu.lookup[cpu.opcode].Operate(cpu)

		cpu.cycles += (additional_cycle1 + additional_cycle2)
		cpu.SetFlag(U, true)
	}

	cpu.cycles--

	return new_inst
}


func (cpu *CPU) GetFlag(flag uint8) uint8 {
	if (cpu.Status & flag) > 0 {
		return 1
	}
	return 0
}


func (cpu *CPU) SetFlag(flag uint8, v bool) {
	if v {
		cpu.Status |= flag
	} else {
		cpu.Status &= ^flag
	}
}


// Reset interrupt
func (cpu *CPU) Reset() {
	cpu.A = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.Stkp = 0xFD
	cpu.Status = 0x00 | U

	cpu.addr_abs = 0xFFFC
	lo := uint16(cpu.Read(cpu.addr_abs + 0))
	hi := uint16(cpu.Read(cpu.addr_abs + 1))

	cpu.Pc = (hi << 8) | lo

	cpu.addr_rel = 0x0000
	cpu.addr_abs = 0x0000
	cpu.fetched = 0x00

	cpu.cycles = 8
}


// Interrupt request
func (cpu *CPU) IRQ() {
	if cpu.GetFlag(I) == 0 {
		cpu.Write(0x0100 + uint16(cpu.Stkp), uint8((cpu.Pc >> 8) & 0x00FF))
		cpu.Stkp--
		cpu.Write(0x0100 + uint16(cpu.Stkp), uint8(cpu.Pc & 0x00FF))
		cpu.Stkp--

		cpu.SetFlag(B, false)
		cpu.SetFlag(U, true)
		cpu.SetFlag(I, true)
		cpu.Write(0x0100 + uint16(cpu.Stkp), cpu.Status)
		cpu.Stkp--

		cpu.addr_abs = 0xFFFE
		lo := uint16(cpu.Read(cpu.addr_abs + 0))
		hi := uint16(cpu.Read(cpu.addr_abs + 1))
		cpu.Pc = (hi << 8) | lo

		cpu.cycles = 7
	}
}


// Non-Maskable interrupt request
func (cpu *CPU) NMI() {
	cpu.Write(0x0100 + uint16(cpu.Stkp), uint8((cpu.Pc >> 8) & 0x00FF))
	cpu.Stkp--
	cpu.Write(0x0100 + uint16(cpu.Stkp), uint8(cpu.Pc & 0x00FF))
	cpu.Stkp--

	cpu.SetFlag(B, false)
	cpu.SetFlag(U, true)
	cpu.SetFlag(I, true)
	cpu.Write(0x0100 + uint16(cpu.Stkp), cpu.Status)
	cpu.Stkp--

	cpu.addr_abs = 0xFFFA
	lo := uint16(cpu.Read(cpu.addr_abs + 0))
	hi := uint16(cpu.Read(cpu.addr_abs + 1))
	cpu.Pc = (hi << 8) | lo

	cpu.cycles = 8
}


func (cpu *CPU) Fetch() uint8 {
	if cpu.lookup[cpu.opcode].ModeType != AddrModeIMP {
		cpu.fetched = cpu.Read(cpu.addr_abs)
	}

	return cpu.fetched
}

// Addressing Modes

// Address Mode: Implied
func (cpu *CPU) IMP() uint8 {
	cpu.fetched = cpu.A

	return 0
}

// Address Mode: Immediate
func (cpu *CPU) IMM() uint8 {
	cpu.addr_abs = cpu.Pc
	cpu.Pc++

	return 0
}

// Address Mode: Zero Page
func (cpu *CPU) ZP0() uint8 {
	cpu.addr_abs = uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	cpu.addr_abs &= 0x00FF

	return 0
}

// Address Mode: Zero Page with X Offset
func (cpu *CPU) ZPX() uint8 {
	cpu.addr_abs = uint16(cpu.Read(cpu.Pc) + cpu.X)
	cpu.Pc++
	cpu.addr_abs &= 0x00FF

	return 0
}

// Address Mode: Zero Page with Y Offset
func (cpu *CPU) ZPY() uint8 {
	cpu.addr_abs = uint16(cpu.Read(cpu.Pc) + cpu.Y)
	cpu.Pc++
	cpu.addr_abs &= 0x00FF

	return 0
}

// Address Mode: Relative
func (cpu *CPU) REL() uint8 {
	cpu.addr_rel = uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	if (cpu.addr_rel & 0x80) != 0 {
		cpu.addr_rel |= 0xFF00
	}

	return 0
}

// Address Mode: Absolute
func (cpu *CPU) ABS() uint8 {
	lo := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	hi := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++

	cpu.addr_abs = (hi << 8) | lo

	return 0
}

// Address Mode: Absolute with X Offset
func (cpu *CPU) ABX() uint8 {
	lo := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	hi := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += uint16(cpu.X)

	if (cpu.addr_abs & 0xFF00) != uint16(hi << 8) {
		return 1
	}

	return 0
}

// Address Mode: Absolute with Y Offset
func (cpu *CPU) ABY() uint8 {
	lo := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	hi := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += uint16(cpu.Y)

	if (cpu.addr_abs & 0xFF00) != uint16(hi << 8) {
		return 1
	}

	return 0
}

// Address Mode: Indirect
func (cpu *CPU) IND() uint8 {
	ptr_lo := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++
	ptr_hi := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++

	ptr := uint16((ptr_hi << 8) | ptr_lo)

	if ptr_lo == 0x00FF {
		cpu.addr_abs = (uint16(cpu.Read(ptr & 0xFF00)) << 8) | uint16(cpu.Read(ptr + 0))
	} else {
		cpu.addr_abs = (uint16(cpu.Read(ptr + 1)) << 8) | uint16(cpu.Read(ptr + 0))
	}

	return 0
}

// Address Mode: Indirect X
func (cpu *CPU) IZX() uint8 {
	t := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++

	lo := uint16(cpu.Read(uint16(t + uint16(cpu.X)) & 0x00FF))
	hi := uint16(cpu.Read(uint16(t + uint16(cpu.X) + 1) & 0x00FF))

	cpu.addr_abs = (hi << 8) | lo

	return 0
}

// Address Mode: Indirect Y
func (cpu *CPU) IZY() uint8 {
	t := uint16(cpu.Read(cpu.Pc))
	cpu.Pc++

	lo := uint16(cpu.Read(t & 0x00FF))
	hi := uint16(cpu.Read((t + 1) & 0x00FF))

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += uint16(cpu.Y)

	if (cpu.addr_abs & 0xFF00) != uint16(hi << 8) {
		return 1
	}

	return 0
}


// Instruction Implementations

// Instruction: Add with Carry In
func (cpu *CPU) ADC() uint8 {
	cpu.Fetch()

	temp := uint16(cpu.A) + uint16(cpu.fetched) + uint16(cpu.GetFlag(C))

	cpu.SetFlag(C, temp > 255)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0)
	signed_overflow_flag := (^(uint16(cpu.A) ^ uint16(cpu.fetched)) & (uint16(cpu.A) ^ uint16(temp))) & 0x0080
	cpu.SetFlag(V, signed_overflow_flag != 0)
	cpu.SetFlag(N, (temp & 0x80) != 0)

	cpu.A = uint8(temp & 0x00FF)

	return 1
}

// Instruction: Subtraction with Borrow In
func (cpu *CPU) SBC() uint8 {
	cpu.Fetch()

	value := uint16(cpu.fetched) ^ 0x00FF
	temp := uint16(cpu.A) + value + uint16(cpu.GetFlag(C))

	cpu.SetFlag(C, (temp & 0xFF00) != 0)
	cpu.SetFlag(Z, ((temp & 0x00FF) == 0))
	signed_overflow_flag := (temp ^ uint16(cpu.A)) & (temp ^ value) & 0x0080
	cpu.SetFlag(V, signed_overflow_flag != 0)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	cpu.A = uint8(temp & 0x00FF)

	return 1
}

// Instruction: Bitwise logic AND
func (cpu *CPU) AND() uint8 {
	cpu.Fetch()
	cpu.A = cpu.A & cpu.fetched
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, (cpu.A & 0x80) != 0)

	return 1
}

// Instruction: Bitwise Shift Left
func (cpu *CPU) ASL() uint8 {
    var value uint8

    if cpu.lookup[cpu.opcode].ModeType == AddrModeIMP {
        value = cpu.A
    } else {
        value = cpu.Fetch()
    }

    temp := uint16(value) << 1
    cpu.SetFlag(C, (temp & 0xFF00) > 0)
    cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
    cpu.SetFlag(N, (temp & 0x80) != 0)

    if cpu.lookup[cpu.opcode].ModeType == AddrModeIMP {
        cpu.A = uint8(temp & 0x00FF)
    } else {
        cpu.Write(cpu.addr_abs, uint8(temp & 0x00FF))
    }

    return 0
}

// Instruction: Branch if Carry Clear
func (cpu *CPU) BCC() uint8 {
	if cpu.GetFlag(C) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Carry Set
func (cpu *CPU) BCS() uint8 {
	if cpu.GetFlag(C) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Equal
func (cpu *CPU) BEQ() uint8 {
	if cpu.GetFlag(Z) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

func (cpu *CPU) BIT() uint8 {
	cpu.Fetch()
	temp := cpu.A & cpu.fetched

	cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
	cpu.SetFlag(N, (cpu.fetched & (1 << 7)) != 0)
	cpu.SetFlag(V, (cpu.fetched & (1 << 6)) != 0)

	return 0
}

// Instruction: Branch if Negative
func (cpu *CPU) BMI() uint8 {
	if cpu.GetFlag(N) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Not Equal
func (cpu *CPU) BNE() uint8 {
	if cpu.GetFlag(Z) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Positive
func (cpu *CPU) BPL() uint8 {
	if cpu.GetFlag(N) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Break
func (cpu *CPU) BRK() uint8 {
	cpu.Pc++

	cpu.SetFlag(I, true)
	cpu.Write(0x0100 + uint16(cpu.Stkp), uint8((cpu.Pc >> 8) & 0x00FF))
	cpu.Stkp--
	cpu.Write(0x0100 + uint16(cpu.Stkp), uint8(cpu.Pc & 0x00FF))
	cpu.Stkp--

	cpu.SetFlag(B, true)
	cpu.Write(0x0100 + uint16(cpu.Stkp), cpu.Status)
	cpu.Stkp--
	cpu.SetFlag(B, false)

	cpu.Pc = uint16(cpu.Read(0xFFFE)) | (uint16(cpu.Read(0xFFFF)) << 8)

	return 0
}

// Instruction: Branch if Overflow Clear
func (cpu *CPU) BVC() uint8 {
	if cpu.GetFlag(V) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Overflow Set
func (cpu *CPU) BVS() uint8 {
	if cpu.GetFlag(V) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.Pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.Pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.Pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Clear Carry Flag
func (cpu *CPU) CLC() uint8 {
	cpu.SetFlag(C, false)

	return 0
}

// Instruction: Clear Decimal Flag
func (cpu *CPU) CLD() uint8 {
	cpu.SetFlag(D, false)

	return 0
}

// Instruction: Disable Interrupts / Clear Interrupt Flag
func (cpu *CPU) CLI() uint8 {
	cpu.SetFlag(I, false)

	return 0
}

// Instruction: Clear Overflow Flag
func (cpu *CPU) CLV() uint8 {
	cpu.SetFlag(V, false)

	return 0
}

// Instruction: Compare Accumulator
func (cpu *CPU) CMP() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.A) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.A >= cpu.fetched)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	return 1
}

// Instruction: Compare X Register
func (cpu *CPU) CPX() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.X) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.X >= cpu.fetched)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	return 0
}

// Instruction: Comapre Y Register
func (cpu *CPU) CPY() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.Y) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.Y >= cpu.fetched)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	return 0
}

// Instruction: Decrement Value at Memory Location
func (cpu *CPU) DEC() uint8 {
	cpu.Fetch()
	temp := cpu.fetched - 1
	cpu.Write(cpu.addr_abs, temp & 0x00FF)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	return 0
}

// Instruction: Decrement X Register
func (cpu *CPU) DEX() uint8 {
	cpu.X--
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, (cpu.X & 0x80) != 0)

	return 0
}

// Instruction: Decrement Y Register
func (cpu *CPU) DEY() uint8 {
	cpu.Y--
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, (cpu.Y & 0x80) != 0)

	return 0
}

// Instruction: Bitwise Logic XOR
func (cpu *CPU) EOR() uint8 {
	cpu.Fetch()
	cpu.A = cpu.A ^ cpu.fetched
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, (cpu.A & 0x80) != 0)

	return 1
}

// Instruction: Increment Value at Memory Location
func (cpu *CPU) INC() uint8 {
	cpu.Fetch()
	temp := cpu.fetched + 1
	cpu.Write(cpu.addr_abs, temp & 0x00FF)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	return 0
}

// Instruction: Increment X Register
func (cpu *CPU) INX() uint8 {
	cpu.X++
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, (cpu.X & 0x80) != 0)

	return 0
}

// Instruction: Increment Y Register
func (cpu *CPU) INY() uint8 {
	cpu.Y++
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, (cpu.Y & 0x80) != 0)

	return 0
}

// Instruction: Jump to Location
func (cpu *CPU) JMP() uint8 {
	cpu.Pc = cpu.addr_abs

	return 0
}

// Instruction: Jump to Sub-Routine
func (cpu *CPU) JSR() uint8 {
	cpu.Pc--

	cpu.Write(0x0100 + uint16(cpu.Stkp), uint8((cpu.Pc >> 8) & 0x00FF))
	cpu.Stkp--
	cpu.Write(0x0100 + uint16(cpu.Stkp), uint8(cpu.Pc & 0x00FF))
	cpu.Stkp--

	cpu.Pc = cpu.addr_abs

	return 0
}

// Instruction: Load The Accumulator
func (cpu *CPU) LDA() uint8 {
	cpu.Fetch()
	cpu.A = cpu.fetched
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, (cpu.A & 0x80) != 0)

	return 1
}

// Instruction: Load the X Register
func (cpu *CPU) LDX() uint8 {
	cpu.Fetch()
	cpu.X = cpu.fetched
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, (cpu.X & 0x80) != 0)

	return 1
}

// Instruction: Load the Y Register
func (cpu *CPU) LDY() uint8 {
	cpu.Fetch()
	cpu.Y = cpu.fetched
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, (cpu.Y & 0x80) != 0)

	return 1
}

func (cpu *CPU) LSR() uint8 {
	cpu.Fetch()
	cpu.SetFlag(C, (cpu.fetched & 0x0001) != 0)
	temp := cpu.fetched >> 1
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	if cpu.lookup[cpu.opcode].ModeType == AddrModeIMP {
		cpu.A = temp & 0x00FF
	} else {
		cpu.Write(cpu.addr_abs, temp & 0x00FF)
	}

	return 0
}

func (cpu *CPU) NOP() uint8 {
	switch cpu.opcode {
	case 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
		return 1
	}
	return 0
}

// Instruction: Bitwise Logic OR
func (cpu *CPU) ORA() uint8 {
	cpu.Fetch()
	cpu.A = cpu.A | cpu.fetched
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, (cpu.A & 0x80) != 0)

	return 0
}

// Instruction: Push Accuumulator to Stack
func (cpu *CPU) PHA() uint8 {
	cpu.Write(0x0100 + uint16(cpu.Stkp), cpu.A)
	cpu.Stkp--

	return 0
}

// Instruction: Push Status Register to Stack
func (cpu *CPU) PHP() uint8 {
	cpu.Write(0x0100 + uint16(cpu.Stkp), cpu.Status | B | U)
	cpu.SetFlag(B, false)
	cpu.SetFlag(U, false)
	cpu.Stkp--

	return 0
}

// Instruction: Pop Accumulator off Stack
func (cpu *CPU) PLA() uint8 {
	cpu.Stkp++
	cpu.A = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, (cpu.A & 0x80) != 0)

	return 0
}

// Instruction: Pop Status Register off Stack
func (cpu *CPU) PLP() uint8 {
	cpu.Stkp++
	cpu.Status = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.SetFlag(U, true)

	return 0
}

func (cpu *CPU) ROL() uint8 {
	if cpu.lookup[cpu.opcode].ModeType == AddrModeIMP {
		carry := cpu.GetFlag(C)
		cpu.SetFlag(C, (cpu.A & 0x80) != 0)
		cpu.A = (cpu.A << 1) | carry
		cpu.SetFlag(Z, cpu.A == 0x00)
		cpu.SetFlag(N, (cpu.A & 0x80) != 0)
	} else {
		cpu.Fetch()
		temp := (cpu.fetched << 1) | cpu.GetFlag(C)
		cpu.SetFlag(C, (cpu.fetched & 0x80) != 0)
		cpu.SetFlag(Z, (temp & 0xFF) == 0x00)
		cpu.SetFlag(N, (temp & 0x80) != 0)
		cpu.Write(cpu.addr_abs, temp&0xFF)
	}
	return 0
}

func (cpu *CPU) ROR() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.GetFlag(C) << 7) | uint16(cpu.fetched >> 1)
	cpu.SetFlag(C, (cpu.fetched & 0x01) != 0)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	if cpu.lookup[cpu.opcode].ModeType == AddrModeIMP {
		cpu.A = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addr_abs, uint8(temp & 0x00FF))
	}

	return 0
}

func (cpu *CPU) RTI() uint8 {
	cpu.Stkp++
	cpu.Status = cpu.Read(0x0100 + uint16(cpu.Stkp))
	cpu.Status = uint8(uint16(cpu.Status) &^ uint16(B))
	cpu.Status = uint8(uint16(cpu.Status) &^ uint16(U))

	cpu.Stkp++
	cpu.Pc = uint16(cpu.Read(0x0100 + uint16(cpu.Stkp)))
	cpu.Stkp++
	cpu.Pc |= uint16(cpu.Read(0x0100 + uint16(cpu.Stkp))) << 8

	return 0
}

func (cpu *CPU) RTS() uint8 {
	cpu.Stkp++
	cpu.Pc = uint16(cpu.Read(0x0100 + uint16(cpu.Stkp)))
	cpu.Stkp++
	cpu.Pc |= uint16(cpu.Read(0x0100 + uint16(cpu.Stkp))) << 8

	cpu.Pc++

	return 0
}

// Instruction: Set Carry Flag
func (cpu *CPU) SEC() uint8 {
	cpu.SetFlag(C, true)

	return 0
}

// Instruction: Set Decimal Flag
func (cpu *CPU) SED() uint8 {
	cpu.SetFlag(D, true)

	return 0
}

// Instruction:  Set Interrupt Flag / Emable Interrupts
func (cpu *CPU) SEI() uint8 {
	cpu.SetFlag(I, true)

	return 0
}

// Instruction: Store Accumulator at Address
func (cpu *CPU) STA() uint8 {
	cpu.Write(cpu.addr_abs, cpu.A)

	return 0
}

// Instruction: Store X Register at Address
func (cpu *CPU) STX() uint8 {
	cpu.Write(cpu.addr_abs, cpu.X)

	return 0
}

// Instruction: Store Y Register at Address
func (cpu *CPU) STY() uint8 {
	cpu.Write(cpu.addr_abs, cpu.Y)

	return 0
}

// Instruction: Transfer Accumulator to X Register
func (cpu *CPU) TAX() uint8 {
	cpu.X = cpu.A
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, (cpu.X & 0x80) != 0)

	return 0
}

// Instruction: Transfer Accumulator to Y Register
func (cpu *CPU) TAY() uint8 {
	cpu.Y = cpu.A
	cpu.SetFlag(Z, cpu.Y == 0x00)
	cpu.SetFlag(N, (cpu.Y & 0x80) != 0)

	return 0
}

// Instruction: Transfer Stack Pointer to X Register
func (cpu *CPU) TSX() uint8 {
	cpu.X = cpu.Stkp
	cpu.SetFlag(Z, cpu.X == 0x00)
	cpu.SetFlag(N, (cpu.X & 0x80) != 0)

	return 0
}

// Instruction: Transfer X Register to Accumulator
func (cpu *CPU) TXA() uint8 {
	cpu.A = cpu.X
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, (cpu.A & 0x80) != 0)

	return 0
}

// Instruction: Transfer X Register to Stack Pointer
func (cpu *CPU) TXS() uint8 {
	cpu.Stkp = cpu.X

	return 0
}

// Instruction: Transfer Y Register to Accumulator
func (cpu *CPU) TYA() uint8 {
	cpu.A = cpu.Y
	cpu.SetFlag(Z, cpu.A == 0x00)
	cpu.SetFlag(N, (cpu.A & 0x80) != 0)

	return 0
}

// Illegal Opcodes
func (cpu *CPU) XXX() uint8 {
	return 0
}


//
// Visualisation functions
//
func (cpu *CPU) PrintCPU() {
	fmt.Println("\n------------------------------------------------------------")
	fmt.Printf("| %-12s | %-13s | %-25s |\n", "Field", "Value", "Description")
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "a", cpu.A, "Accumulator register")
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "x", cpu.X, "X register")
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "y", cpu.Y, "Y register")
	fmt.Printf("| %-12s | $%-12.4X | %-25s |\n", "stkp", cpu.Stkp, "Stack pointer")
	fmt.Printf("| %-12s | $%-12.4X | %-25s |\n", "pc", cpu.Pc, "Program counter")
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "fetched", cpu.fetched, "Working input to ALU")
	fmt.Printf("| %-12s | $%-12.4X | %-25s |\n", "addr_abs", cpu.addr_abs, "Absolute address")
	fmt.Printf("| %-12s | $%-12.4X | %-25s |\n", "addr_rel", cpu.addr_rel, "Relative address")
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "opcode", cpu.opcode, "CPU opcode")
	fmt.Printf("| %-12s | %-12d  | %-25s |\n", "cycles", cpu.cycles, "Cycles left")
	fmt.Printf("| %-12s | %-12d  | %-25s |\n", "clock_count", cpu.clock_count, "Number of clocks")
	fmt.Println("------------------------------------------------------------")
}

func (cpu *CPU) PrintStatusFlags() {
    // Define an array of flag names and their corresponding constants
    flags := []struct {
        name  string
        value uint8
    }{
        {"C (Carry)", C},
        {"Z (Zero)", Z},
        {"I (Disable interrupts)", I},
        {"D (Decimal mode)", D},
        {"B (Break)", B},
        {"U (Unused)", U},
        {"V (Overflow)", V},
        {"N (Negative)", N},
    }

    fmt.Println("\n--------------------------------------")
    fmt.Printf("| %-25s | %-6s |\n", "Flag", "Status")
    fmt.Println("--------------------------------------")

    // Iterate through flags and print their Status
    for _, flag := range flags {
        Status := 0
        if cpu.Status&flag.value != 0 {
            Status = 1
        }
        fmt.Printf("| %-25s | %-6d |\n", flag.name, Status)
    }

    fmt.Println("--------------------------------------")
}


func (cpu *CPU) PrintRAM(startPage int, pages int) {
	cpu.bus.PrintRAM(startPage, pages)
}


func (cpu *CPU) Disassemble(nStart uint16, nStop uint16) map[uint16]string {
	addr := uint32(nStart)
	value := uint8(0x00)
	lo := uint8(0x00)
	hi := uint8(0x00)
	mapLines := make(map[uint16]string)
	line_addr := uint16(0)

	for addr <= uint32(nStop) {
		line_addr = uint16(addr)

		sInst := "$" + fmt.Sprintf("%04x", addr) + ": "

		opcode := cpu.bus.CpuRead(uint16(addr), true)
		addr++
		sInst += cpu.lookup[opcode].Name + " "

		switch cpu.lookup[opcode].ModeType {
			case AddrModeIMP:
				sInst += " {IMP}"
			case AddrModeIMM:
				value = cpu.bus.CpuRead(uint16(addr), true)
				addr++
				sInst += "#$" + fmt.Sprintf("%02x", value) + " {IMM}"
			case AddrModeZP0:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = 0x00
		        sInst += "$" + fmt.Sprintf("%02x", lo) + " {ZP0}"
		    case AddrModeZPX:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = 0x00
		        sInst += "$" + fmt.Sprintf("%02x", lo) + ", X {ZPX}"
		    case AddrModeZPY:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = 0x00
		        sInst += "$" + fmt.Sprintf("%02x", lo) + ", Y {ZPY}"
		    case AddrModeIZX:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = 0x00
		        sInst += "($" + fmt.Sprintf("%02x", lo) + ", X) {IZX}"
		    case AddrModeIZY:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = 0x00
		        sInst += "($" + fmt.Sprintf("%02x", lo) + "), Y {IZY}"
		    case AddrModeABS:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        sInst += "$" + fmt.Sprintf("%04x", uint16(hi)<<8|uint16(lo)) + " {ABS}"
		    case AddrModeABX:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        sInst += "$" + fmt.Sprintf("%04x", uint16(hi)<<8|uint16(lo)) + ", X {ABX}"
		    case AddrModeABY:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        sInst += "$" + fmt.Sprintf("%04x", uint16(hi)<<8|uint16(lo)) + ", Y {ABY}"
		    case AddrModeIND:
		        lo = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        hi = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        sInst += "($" + fmt.Sprintf("%04x", uint16(hi)<<8|uint16(lo)) + ") {IND}"
		    case AddrModeREL:
		        value = cpu.bus.CpuRead(uint16(addr), true)
		        addr++
		        sInst += "$" + fmt.Sprintf("%02x", value) + " [$" + fmt.Sprintf("%04x", addr+uint32(value)) + "] {REL}"
		    default:
		    	sInst += ""
		}

		mapLines[line_addr] = sInst
	}

	return mapLines
}
