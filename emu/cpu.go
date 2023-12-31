package emu

import (
	"fmt"
	"reflect"
)


type CPU struct {
	bus *Bus
	a uint8  // Accumulator register
	x uint8  // X register
	y uint8  // Y register
	stkp uint8  // Stack pointer
	pc uint16  // Program counter
	status uint8  // Status register
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
	cpu.a = 0x00
	cpu.x = 0x00
	cpu.y = 0x00
	cpu.stkp = 0x00
	cpu.pc = 0x0000
	cpu.status = 0x00

	cpu.fetched = 0x00
	cpu.addr_abs = 0x0000
	cpu.addr_rel = 0x00
	cpu.opcode = 0x00
	cpu.cycles = 0

	cpu.lookup = []INSTRUCTION {
		INSTRUCTION{ "BRK", (*CPU).BRK, (*CPU).IMM, 7 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).IZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 3 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ZP0, 3 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "PHP", (*CPU).PHP, (*CPU).IMP, 3 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).IMM, 2 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ABS, 4 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		INSTRUCTION{ "BPL", (*CPU).BPL, (*CPU).REL, 2 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).IZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ZPX, 4 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "CLC", (*CPU).CLC, (*CPU).IMP, 2 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "ORA", (*CPU).ORA, (*CPU).ABX, 4 },INSTRUCTION{ "ASL", (*CPU).ASL, (*CPU).ABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		INSTRUCTION{ "JSR", (*CPU).JSR, (*CPU).ABS, 6 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).IZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "BIT", (*CPU).BIT, (*CPU).ZP0, 3 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ZP0, 3 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "PLP", (*CPU).PLP, (*CPU).IMP, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).IMM, 2 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "BIT", (*CPU).BIT, (*CPU).ABS, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ABS, 4 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		INSTRUCTION{ "BMI", (*CPU).BMI, (*CPU).REL, 2 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).IZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ZPX, 4 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "SEC", (*CPU).SEC, (*CPU).IMP, 2 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "AND", (*CPU).AND, (*CPU).ABX, 4 },INSTRUCTION{ "ROL", (*CPU).ROL, (*CPU).ABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		INSTRUCTION{ "RTI", (*CPU).RTI, (*CPU).IMP, 6 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).IZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 3 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ZP0, 3 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "PHA", (*CPU).PHA, (*CPU).IMP, 3 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).IMM, 2 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "JMP", (*CPU).JMP, (*CPU).ABS, 3 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ABS, 4 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		INSTRUCTION{ "BVC", (*CPU).BVC, (*CPU).REL, 2 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).IZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ZPX, 4 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "CLI", (*CPU).CLI, (*CPU).IMP, 2 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "EOR", (*CPU).EOR, (*CPU).ABX, 4 },INSTRUCTION{ "LSR", (*CPU).LSR, (*CPU).ABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		INSTRUCTION{ "RTS", (*CPU).RTS, (*CPU).IMP, 6 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).IZX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 3 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ZP0, 3 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "PLA", (*CPU).PLA, (*CPU).IMP, 4 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).IMM, 2 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "JMP", (*CPU).JMP, (*CPU).IND, 5 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ABS, 4 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		INSTRUCTION{ "BVS", (*CPU).BVS, (*CPU).REL, 2 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).IZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ZPX, 4 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "SEI", (*CPU).SEI, (*CPU).IMP, 2 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ABY, 4 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "ADC", (*CPU).ADC, (*CPU).ABX, 4 },INSTRUCTION{ "ROR", (*CPU).ROR, (*CPU).ABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).IZX, 6 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "STY", (*CPU).STY, (*CPU).ZP0, 3 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ZP0, 3 },INSTRUCTION{ "STX", (*CPU).STX, (*CPU).ZP0, 3 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 3 },INSTRUCTION{ "DEY", (*CPU).DEY, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "TXA", (*CPU).TXA, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "STY", (*CPU).STY, (*CPU).ABS, 4 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ABS, 4 },INSTRUCTION{ "STX", (*CPU).STX, (*CPU).ABS, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 4 },
		INSTRUCTION{ "BCC", (*CPU).BCC, (*CPU).REL, 2 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).IZY, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "STY", (*CPU).STY, (*CPU).ZPX, 4 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ZPX, 4 },INSTRUCTION{ "STX", (*CPU).STX, (*CPU).ZPY, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 4 },INSTRUCTION{ "TYA", (*CPU).TYA, (*CPU).IMP, 2 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ABY, 5 },INSTRUCTION{ "TXS", (*CPU).TXS, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 5 },INSTRUCTION{ "STA", (*CPU).STA, (*CPU).ABX, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },
		INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).IMM, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).IZX, 6 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).IMM, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ZP0, 3 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ZP0, 3 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ZP0, 3 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 3 },INSTRUCTION{ "TAY", (*CPU).TAY, (*CPU).IMP, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).IMM, 2 },INSTRUCTION{ "TAX", (*CPU).TAX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ABS, 4 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ABS, 4 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ABS, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 4 },
		INSTRUCTION{ "BCS", (*CPU).BCS, (*CPU).REL, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).IZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ZPX, 4 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ZPX, 4 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ZPY, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 4 },INSTRUCTION{ "CLV", (*CPU).CLV, (*CPU).IMP, 2 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ABY, 4 },INSTRUCTION{ "TSX", (*CPU).TSX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 4 },INSTRUCTION{ "LDY", (*CPU).LDY, (*CPU).ABX, 4 },INSTRUCTION{ "LDA", (*CPU).LDA, (*CPU).ABX, 4 },INSTRUCTION{ "LDX", (*CPU).LDX, (*CPU).ABY, 4 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 4 },
		INSTRUCTION{ "CPY", (*CPU).CPY, (*CPU).IMM, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).IZX, 6 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "CPY", (*CPU).CPY, (*CPU).ZP0, 3 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ZP0, 3 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "INY", (*CPU).INY, (*CPU).IMP, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).IMM, 2 },INSTRUCTION{ "DEX", (*CPU).DEX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "CPY", (*CPU).CPY, (*CPU).ABS, 4 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ABS, 4 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		INSTRUCTION{ "BNE", (*CPU).BNE, (*CPU).REL, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).IZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ZPX, 4 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "CLD", (*CPU).CLD, (*CPU).IMP, 2 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ABY, 4 },INSTRUCTION{ "NOP", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "CMP", (*CPU).CMP, (*CPU).ABX, 4 },INSTRUCTION{ "DEC", (*CPU).DEC, (*CPU).ABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		INSTRUCTION{ "CPX", (*CPU).CPX, (*CPU).IMM, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).IZX, 6 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "CPX", (*CPU).CPX, (*CPU).ZP0, 3 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ZP0, 3 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ZP0, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 5 },INSTRUCTION{ "INX", (*CPU).INX, (*CPU).IMP, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).IMM, 2 },INSTRUCTION{ "NOP", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).SBC, (*CPU).IMP, 2 },INSTRUCTION{ "CPX", (*CPU).CPX, (*CPU).ABS, 4 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ABS, 4 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ABS, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		INSTRUCTION{ "BEQ", (*CPU).BEQ, (*CPU).REL, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).IZY, 5 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 8 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ZPX, 4 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ZPX, 6 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 6 },INSTRUCTION{ "SED", (*CPU).SED, (*CPU).IMP, 2 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ABY, 4 },INSTRUCTION{ "NOP", (*CPU).NOP, (*CPU).IMP, 2 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },INSTRUCTION{ "???", (*CPU).NOP, (*CPU).IMP, 4 },INSTRUCTION{ "SBC", (*CPU).SBC, (*CPU).ABX, 4 },INSTRUCTION{ "INC", (*CPU).INC, (*CPU).ABX, 7 },INSTRUCTION{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
	}

	return &cpu
}


func (cpu *CPU) ConnectBus(n *Bus) {
	cpu.bus = n
}


func (cpu *CPU) Read(a uint16) uint8 {
	return cpu.bus.Read(a, false)
}


func (cpu *CPU) Write(a uint16, d uint8) {
	cpu.bus.Write(a, d)
}


// Updates the clock cycle
func (cpu *CPU) Clock() bool {
	new_inst := cpu.cycles == 0

	if cpu.cycles == 0 {
		cpu.opcode = cpu.Read(cpu.pc)
		cpu.pc++

		cpu.cycles = cpu.lookup[cpu.opcode].Cycles

		additional_cycle1 := cpu.lookup[cpu.opcode].AddrMode(cpu)

		additional_cycle2 := cpu.lookup[cpu.opcode].Operate(cpu)

		cpu.cycles += (additional_cycle1 + additional_cycle2)
	}

	cpu.cycles--

	return new_inst
}


func (cpu *CPU) GetFlag(flag uint8) uint8 {
	if (cpu.status & flag) > 0 {
		return 1
	}
	return 0
}


func (cpu *CPU) SetFlag(flag uint8, v bool) {
	if v {
		cpu.status |= flag
	} else {
		cpu.status &= ^flag
	}
}


// Reset interrupt
func (cpu *CPU) Reset() {
	cpu.a = 0
	cpu.x = 0
	cpu.y = 0
	cpu.stkp = 0xFD
	cpu.status = 0x00 | U

	cpu.addr_abs = 0xFFFC
	lo := uint16(cpu.Read(cpu.addr_abs + 0))
	hi := uint16(cpu.Read(cpu.addr_abs + 1))

	cpu.pc = (hi << 8) | lo

	cpu.addr_rel = 0x0000
	cpu.addr_abs = 0x0000
	cpu.fetched = 0x00

	cpu.cycles = 8
}


// Interrupt request
func (cpu *CPU) IRQ() {
	if cpu.GetFlag(I) == 0 {
		cpu.Write(0x0100 + uint16(cpu.stkp), uint8((cpu.pc >> 8) & 0x00FF))
		cpu.stkp--
		cpu.Write(0x0100 + uint16(cpu.stkp), uint8(cpu.pc & 0x00FF))
		cpu.stkp--

		cpu.SetFlag(B, false)
		cpu.SetFlag(U, true)
		cpu.SetFlag(I, true)
		cpu.Write(0x0100 + uint16(cpu.stkp), cpu.status)
		cpu.stkp--

		cpu.addr_abs = 0xFFFE
		lo := uint16(cpu.Read(cpu.addr_abs + 0))
		hi := uint16(cpu.Read(cpu.addr_abs + 1))
		cpu.pc = (hi << 8) | lo

		cpu.cycles = 7
	}
}


// Non-Maskable interrupt request
func (cpu *CPU) NMI() {
	cpu.Write(0x0100 + uint16(cpu.stkp), uint8((cpu.pc >> 8) & 0x00FF))
	cpu.stkp--
	cpu.Write(0x0100 + uint16(cpu.stkp), uint8(cpu.pc & 0x00FF))
	cpu.stkp--

	cpu.SetFlag(B, false)
	cpu.SetFlag(U, true)
	cpu.SetFlag(I, true)
	cpu.Write(0x0100 + uint16(cpu.stkp), cpu.status)
	cpu.stkp--

	cpu.addr_abs = 0xFFFA
	lo := uint16(cpu.Read(cpu.addr_abs + 0))
	hi := uint16(cpu.Read(cpu.addr_abs + 1))
	cpu.pc = (hi << 8) | lo

	cpu.cycles = 8
}


func (cpu *CPU) Fetch() uint8 {
	if !(reflect.ValueOf(cpu.lookup[cpu.opcode].AddrMode).Pointer() == reflect.ValueOf((*CPU).IMP).Pointer()) {
		cpu.fetched = cpu.Read(cpu.addr_abs)
	}

	return cpu.fetched
}

// Addressing Modes

// Address Mode: Implied
func (cpu *CPU) IMP() uint8 {
	cpu.fetched = cpu.a

	return 0
}

// Address Mode: Immediate
func (cpu *CPU) IMM() uint8 {
	cpu.addr_abs = cpu.pc
	cpu.pc++

	return 0
}

// Address Mode: Zero Page
func (cpu *CPU) ZP0() uint8 {
	cpu.addr_abs = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	cpu.addr_abs &= 0x00FF

	return 0
}

// Address Mode: Zero Page with X Offset
func (cpu *CPU) ZPX() uint8 {
	cpu.addr_abs = uint16(cpu.Read(cpu.pc) + cpu.x)
	cpu.pc++
	cpu.addr_abs &= 0x00FF

	return 0
}

// Address Mode: Zero Page with Y Offset
func (cpu *CPU) ZPY() uint8 {
	cpu.addr_abs = uint16(cpu.Read(cpu.pc) + cpu.y)
	cpu.pc++
	cpu.addr_abs &= 0x00FF

	return 0
}

// Address Mode: Relative
func (cpu *CPU) REL() uint8 {
	cpu.addr_rel = uint16(cpu.Read(cpu.pc))
	cpu.pc++
	if (cpu.addr_rel & 0x80) != 0 {
		cpu.addr_rel |= 0xFF00
	}

	return 0
}

// Address Mode: Absolute
func (cpu *CPU) ABS() uint8 {
	lo := uint16(cpu.Read(cpu.pc))
	cpu.pc++
	hi := uint16(cpu.Read(cpu.pc))
	cpu.pc++

	cpu.addr_abs = (hi << 8) | lo

	return 0
}

// Address Mode: Absolute with X Offset
func (cpu *CPU) ABX() uint8 {
	lo := uint16(cpu.Read(cpu.pc))
	cpu.pc++
	hi := uint16(cpu.Read(cpu.pc))
	cpu.pc++

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += uint16(cpu.x)

	if (cpu.addr_abs & 0xFF00) != uint16(hi << 8) {
		return 1
	}

	return 0
}

// Address Mode: Absolute with Y Offset
func (cpu *CPU) ABY() uint8 {
	lo := uint16(cpu.Read(cpu.pc))
	cpu.pc++
	hi := uint16(cpu.Read(cpu.pc))
	cpu.pc++

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += uint16(cpu.y)

	if (cpu.addr_abs & 0xFF00) != uint16(hi << 8) {
		return 1
	}

	return 0
}

// Address Mode: Indirect
func (cpu *CPU) IND() uint8 {
	ptr_lo := uint16(cpu.Read(cpu.pc))
	cpu.pc++
	ptr_hi := uint16(cpu.Read(cpu.pc))
	cpu.pc++

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
	t := uint16(cpu.Read(cpu.pc))
	cpu.pc++

	lo := uint16(cpu.Read(uint16(t + uint16(cpu.x)) & 0x00FF))
	hi := uint16(cpu.Read(uint16(t + uint16(cpu.x) + 1) & 0x00FF))

	cpu.addr_abs = (hi << 8) | lo

	return 0
}

// Address Mode: Indirect Y
func (cpu *CPU) IZY() uint8 {
	t := uint16(cpu.Read(cpu.pc))
	cpu.pc++

	lo := uint16(cpu.Read(t & 0x00FF))
	hi := uint16(cpu.Read((t + 1) & 0x00FF))

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += uint16(cpu.y)

	if (cpu.addr_abs & 0xFF00) != uint16(hi << 8) {
		return 1
	}

	return 0
}


// Instruction Implementations

// Instruction: Add with Carry In
func (cpu *CPU) ADC() uint8 {
	cpu.Fetch()

	temp := uint16(cpu.a) + uint16(cpu.fetched) + uint16(cpu.GetFlag(C))

	cpu.SetFlag(C, temp > 255)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0)
	signed_overflow_flag := (^(uint16(cpu.a) ^ uint16(cpu.fetched)) & (uint16(cpu.a) ^ uint16(temp))) & 0x0080
	cpu.SetFlag(V, signed_overflow_flag != 0)
	cpu.SetFlag(N, (temp & 0x80) != 0)

	cpu.a = uint8(temp & 0x00FF)

	return 1
}

// Instruction: Subtraction with Borrow In
func (cpu *CPU) SBC() uint8 {
	cpu.Fetch()

	value := uint16(cpu.fetched) ^ 0x00FF
	temp := uint16(cpu.a) + value + uint16(cpu.GetFlag(C))

	cpu.SetFlag(C, (temp & 0xFF00) != 0)
	cpu.SetFlag(Z, ((temp & 0x00FF) == 0))
	signed_overflow_flag := (temp ^ uint16(cpu.a)) & (temp ^ value) & 0x0080
	cpu.SetFlag(V, signed_overflow_flag != 0)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	cpu.a = uint8(temp & 0x00FF)

	return 1
}

// Instruction: Bitwise logic AND
func (cpu *CPU) AND() uint8 {
	cpu.Fetch()
	cpu.a = cpu.a & cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, (cpu.a & 0x00) != 0)

	return 1
}

// Instruction: Bitwise Shift Left
func (cpu *CPU) ASL() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.fetched << 1)
	cpu.SetFlag(C, (temp & 0xFF00) > 0)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
	cpu.SetFlag(N, (temp & 0x80) != 0)

	if reflect.ValueOf(cpu.lookup[cpu.opcode].AddrMode).Pointer() == reflect.ValueOf((*CPU).IMP).Pointer() {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addr_abs, uint8(temp & 0x00FF))
	}

	return 0
}

// Instruction: Branch if Carry Clear
func (cpu *CPU) BCC() uint8 {
	if cpu.GetFlag(C) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Carry Set
func (cpu *CPU) BCS() uint8 {
	if cpu.GetFlag(C) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Equal
func (cpu *CPU) BEQ() uint8 {
	if cpu.GetFlag(Z) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
	}

	return 0
}

func (cpu *CPU) BIT() uint8 {
	cpu.Fetch()
	temp := cpu.a & cpu.fetched

	cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
	cpu.SetFlag(N, (cpu.fetched & (1 << 7)) != 0)
	cpu.SetFlag(V, (cpu.fetched & (1 << 6)) != 0)

	return 0
}

// Instruction: Branch if Negative
func (cpu *CPU) BMI() uint8 {
	if cpu.GetFlag(N) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Not Equal
func (cpu *CPU) BNE() uint8 {
	if cpu.GetFlag(Z) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Positive
func (cpu *CPU) BPL() uint8 {
	if cpu.GetFlag(N) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Break
func (cpu *CPU) BRK() uint8 {
	cpu.pc++

	cpu.SetFlag(I, true)
	cpu.Write(0x0100 + uint16(cpu.stkp), uint8((cpu.pc >> 8) & 0x00FF))
	cpu.stkp--
	cpu.Write(0x0100 + uint16(cpu.stkp), uint8(cpu.pc & 0x00FF))
	cpu.stkp--

	cpu.SetFlag(B, true)
	cpu.Write(0x0100 + uint16(cpu.stkp), cpu.status)
	cpu.stkp--
	cpu.SetFlag(B, false)

	cpu.pc = uint16(cpu.Read(0xFFFE)) | (uint16(cpu.Read(0xFFF)) << 8)

	return 0
}

// Instruction: Branch if Overflow Clear
func (cpu *CPU) BVC() uint8 {
	if cpu.GetFlag(V) == 0 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
	}

	return 0
}

// Instruction: Branch if Overflow Set
func (cpu *CPU) BVS() uint8 {
	if cpu.GetFlag(V) == 1 {
		cpu.cycles++
		cpu.addr_abs = cpu.pc + cpu.addr_rel

		if (cpu.addr_abs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}

		cpu.pc = cpu.addr_abs
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
	temp := uint16(cpu.a) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.a >= cpu.fetched)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	return 1
}

// Instruction: Compare X Register
func (cpu *CPU) CPX() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.x) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.x >= cpu.fetched)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	return 0
}

// Instruction: Comapre Y Register
func (cpu *CPU) CPY() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.y) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.y >= cpu.fetched)
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
	cpu.x--
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, (cpu.x & 0x80) != 0)

	return 0
}

// Instruction: Decrement Y Register
func (cpu *CPU) DEY() uint8 {
	cpu.y--
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, (cpu.y & 0x80) != 0)

	return 0
}

// Instruction: Bitwise Logic XOR
func (cpu *CPU) EOR() uint8 {
	cpu.Fetch()
	cpu.a = cpu.a ^ cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, (cpu.a & 0x80) != 0)

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
	cpu.x++
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, (cpu.x & 0x80) != 0)

	return 0
}

// Instruction: Increment Y Register
func (cpu *CPU) INY() uint8 {
	cpu.y++
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, (cpu.y & 0x80) != 0)

	return 0
}

// Instruction: Jump to Location
func (cpu *CPU) JMP() uint8 {
	cpu.pc = cpu.addr_abs

	return 0
}

// Instruction: Jump to Sub-Routine
func (cpu *CPU) JSR() uint8 {
	cpu.pc--

	cpu.Write(0x0100 + uint16(cpu.stkp), uint8((cpu.pc >> 8) & 0x00FF))
	cpu.stkp--
	cpu.Write(0x0100 + uint16(cpu.stkp), uint8(cpu.pc & 0x00FF))
	cpu.stkp--

	cpu.pc = cpu.addr_abs

	return 0
}

// Instruction: Load The Accumulator
func (cpu *CPU) LDA() uint8 {
	cpu.Fetch()
	cpu.a = cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, (cpu.a & 0x80) != 0)

	return 1
}

// Instruction: Load the X Register
func (cpu *CPU) LDX() uint8 {
	cpu.Fetch()
	cpu.x = cpu.fetched
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, (cpu.x & 0x80) != 0)

	return 1
}

// Instruction: Load the Y Register
func (cpu *CPU) LDY() uint8 {
	cpu.Fetch()
	cpu.y = cpu.fetched
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, (cpu.y & 0x80) != 0)

	return 1
}

func (cpu *CPU) LSR() uint8 {
	cpu.Fetch()
	cpu.SetFlag(C, (cpu.fetched & 0x0001) != 0)
	temp := cpu.fetched >> 1
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	if reflect.ValueOf(cpu.lookup[cpu.opcode].AddrMode).Pointer() == reflect.ValueOf((*CPU).IMP).Pointer() {
		cpu.a = temp & 0x00FF
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
	cpu.a = cpu.a | cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, (cpu.a & 0x80) != 0)

	return 0
}

// Instruction: Push Accuumulator to Stack
func (cpu *CPU) PHA() uint8 {
	cpu.Write(0x0100 + uint16(cpu.stkp), cpu.a)
	cpu.stkp--

	return 0
}

// Instruction: Push Status Register to Stack
func (cpu *CPU) PHP() uint8 {
	cpu.Write(0x0100 + uint16(cpu.stkp), cpu.status | B | U)
	cpu.SetFlag(B, false)
	cpu.SetFlag(U, false)
	cpu.stkp--

	return 0
}

// Instruction: Pop Accumulator off Stack
func (cpu *CPU) PLA() uint8 {
	cpu.stkp++
	cpu.a = cpu.Read(0x0100 + uint16(cpu.stkp))
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, (cpu.a & 0x80) != 0)

	return 0
}

// Instruction: Pop Status Register off Stack
func (cpu *CPU) PLP() uint8 {
	cpu.stkp++
	cpu.status = cpu.Read(0x0100 + uint16(cpu.stkp))
	cpu.SetFlag(U, true)

	return 0
}

func (cpu *CPU) ROL() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.fetched << 1) | uint16(cpu.GetFlag(C))
	cpu.SetFlag(C, (temp & 0xFF00) != 0)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	if reflect.ValueOf(cpu.lookup[cpu.opcode].AddrMode).Pointer() == reflect.ValueOf((*CPU).IMP).Pointer() {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addr_abs, uint8(temp & 0x00FF))
	}

	return 0
}

func (cpu *CPU) ROR() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.GetFlag(C) << 7) | uint16(cpu.fetched >> 1)
	cpu.SetFlag(C, (cpu.fetched & 0x01) != 0)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
	cpu.SetFlag(N, (temp & 0x0080) != 0)

	if reflect.ValueOf(cpu.lookup[cpu.opcode].AddrMode).Pointer() == reflect.ValueOf((*CPU).IMP).Pointer() {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.Write(cpu.addr_abs, uint8(temp & 0x00FF))
	}

	return 0
}

func (cpu *CPU) RTI() uint8 {
	cpu.stkp++
	cpu.status = cpu.Read(0x0100 + uint16(cpu.stkp))
	cpu.status = uint8(uint16(cpu.status) &^ uint16(B))
	cpu.status = uint8(uint16(cpu.status) &^ uint16(U))

	cpu.stkp++
	cpu.pc = uint16(cpu.Read(0x0100 + uint16(cpu.stkp)))
	cpu.stkp++
	cpu.pc |= uint16(cpu.Read(0x0100 + uint16(cpu.stkp))) << 8

	return 0
}

func (cpu *CPU) RTS() uint8 {
	cpu.stkp++
	cpu.pc = uint16(cpu.Read(0x0100 + uint16(cpu.stkp)))
	cpu.stkp++
	cpu.pc |= uint16(cpu.Read(0x0100 + uint16(cpu.stkp))) << 8

	cpu.pc++

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
	cpu.Write(cpu.addr_abs, cpu.a)

	return 0
}

// Instruction: Store X Register at Address
func (cpu *CPU) STX() uint8 {
	cpu.Write(cpu.addr_abs, cpu.x)

	return 0
}

// Instruction: Store Y Register at Address
func (cpu *CPU) STY() uint8 {
	cpu.Write(cpu.addr_abs, cpu.y)

	return 0
}

// Instruction: Transfer Accumulator to X Register
func (cpu *CPU) TAX() uint8 {
	cpu.x = cpu.a
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, (cpu.x & 0x80) != 0)

	return 0
}

// Instruction: Transfer Accumulator to Y Register
func (cpu *CPU) TAY() uint8 {
	cpu.y = cpu.a
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, (cpu.y & 0x80) != 0)

	return 0
}

// Instruction: Transfer Stack Pointer to X Register
func (cpu *CPU) TSX() uint8 {
	cpu.x = cpu.stkp
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, (cpu.x & 0x80) != 0)

	return 0
}

// Instruction: Transfer X Register to Accumulator
func (cpu *CPU) TXA() uint8 {
	cpu.a = cpu.x
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, (cpu.a & 0x80) != 0)

	return 0
}

// Instruction: Transfer X Register to Stack Pointer
func (cpu *CPU) TXS() uint8 {
	cpu.stkp = cpu.x

	return 0
}

// Instruction: Transfer Y Register to Accumulator
func (cpu *CPU) TYA() uint8 {
	cpu.a = cpu.y
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, (cpu.a & 0x80) != 0)

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
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "a", cpu.a, "Accumulator register")
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "x", cpu.x, "X register")
	fmt.Printf("| %-12s | $%-12.2X | %-25s |\n", "y", cpu.y, "Y register")
	fmt.Printf("| %-12s | $%-12.4X | %-25s |\n", "stkp", cpu.stkp, "Stack pointer")
	fmt.Printf("| %-12s | $%-12.4X | %-25s |\n", "pc", cpu.pc, "Program counter")
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

    // Iterate through flags and print their status
    for _, flag := range flags {
        status := 0
        if cpu.status&flag.value != 0 {
            status = 1
        }
        fmt.Printf("| %-25s | %-6d |\n", flag.name, status)
    }

    fmt.Println("--------------------------------------")
}


func (cpu *CPU) PrintRAM(startPage int, pages int) {
	cpu.bus.PrintRAM(startPage, pages)
}

