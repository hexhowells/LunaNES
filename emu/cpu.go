type CPU struct {
	bus Bus
	a unit8  // Accumulator register
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
	cpu := cpu{}
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

	lookup := []INSTRUCTION
	{
		{ "BRK", (*CPU).BRK, (*CPU).IMM, 7 },{ "ORA", (*CPU).ORA, (*CPU).IZX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 3 },{ "ORA", (*CPU).ORA, (*CPU).ZP0, 3 },{ "ASL", (*CPU).ASL, (*CPU).ZP0, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "PHP", (*CPU).PHP, (*CPU).IMP, 3 },{ "ORA", (*CPU).ORA, (*CPU).IMM, 2 },{ "ASL", (*CPU).ASL, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "ORA", (*CPU).ORA, (*CPU).ABS, 4 },{ "ASL", (*CPU).ASL, (*CPU).ABS, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		{ "BPL", (*CPU).BPL, (*CPU).REL, 2 },{ "ORA", (*CPU).ORA, (*CPU).IZY, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "ORA", (*CPU).ORA, (*CPU).ZPX, 4 },{ "ASL", (*CPU).ASL, (*CPU).ZPX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "CLC", (*CPU).CLC, (*CPU).IMP, 2 },{ "ORA", (*CPU).ORA, (*CPU).ABY, 4 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "ORA", (*CPU).ORA, (*CPU).ABX, 4 },{ "ASL", (*CPU).ASL, (*CPU).ABX, 7 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		{ "JSR", (*CPU).JSR, (*CPU).ABS, 6 },{ "AND", (*CPU).AND, (*CPU).IZX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "BIT", (*CPU).BIT, (*CPU).ZP0, 3 },{ "AND", (*CPU).AND, (*CPU).ZP0, 3 },{ "ROL", (*CPU).ROL, (*CPU).ZP0, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "PLP", (*CPU).PLP, (*CPU).IMP, 4 },{ "AND", (*CPU).AND, (*CPU).IMM, 2 },{ "ROL", (*CPU).ROL, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "BIT", (*CPU).BIT, (*CPU).ABS, 4 },{ "AND", (*CPU).AND, (*CPU).ABS, 4 },{ "ROL", (*CPU).ROL, (*CPU).ABS, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		{ "BMI", (*CPU).BMI, (*CPU).REL, 2 },{ "AND", (*CPU).AND, (*CPU).IZY, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "AND", (*CPU).AND, (*CPU).ZPX, 4 },{ "ROL", (*CPU).ROL, (*CPU).ZPX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "SEC", (*CPU).SEC, (*CPU).IMP, 2 },{ "AND", (*CPU).AND, (*CPU).ABY, 4 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "AND", (*CPU).AND, (*CPU).ABX, 4 },{ "ROL", (*CPU).ROL, (*CPU).ABX, 7 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		{ "RTI", (*CPU).RTI, (*CPU).IMP, 6 },{ "EOR", (*CPU).EOR, (*CPU).IZX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 3 },{ "EOR", (*CPU).EOR, (*CPU).ZP0, 3 },{ "LSR", (*CPU).LSR, (*CPU).ZP0, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "PHA", (*CPU).PHA, (*CPU).IMP, 3 },{ "EOR", (*CPU).EOR, (*CPU).IMM, 2 },{ "LSR", (*CPU).LSR, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "JMP", (*CPU).JMP, (*CPU).ABS, 3 },{ "EOR", (*CPU).EOR, (*CPU).ABS, 4 },{ "LSR", (*CPU).LSR, (*CPU).ABS, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		{ "BVC", (*CPU).BVC, (*CPU).REL, 2 },{ "EOR", (*CPU).EOR, (*CPU).IZY, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "EOR", (*CPU).EOR, (*CPU).ZPX, 4 },{ "LSR", (*CPU).LSR, (*CPU).ZPX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "CLI", (*CPU).CLI, (*CPU).IMP, 2 },{ "EOR", (*CPU).EOR, (*CPU).ABY, 4 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "EOR", (*CPU).EOR, (*CPU).ABX, 4 },{ "LSR", (*CPU).LSR, (*CPU).ABX, 7 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		{ "RTS", (*CPU).RTS, (*CPU).IMP, 6 },{ "ADC", (*CPU).ADC, (*CPU).IZX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 3 },{ "ADC", (*CPU).ADC, (*CPU).ZP0, 3 },{ "ROR", (*CPU).ROR, (*CPU).ZP0, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "PLA", (*CPU).PLA, (*CPU).IMP, 4 },{ "ADC", (*CPU).ADC, (*CPU).IMM, 2 },{ "ROR", (*CPU).ROR, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "JMP", (*CPU).JMP, (*CPU).IND, 5 },{ "ADC", (*CPU).ADC, (*CPU).ABS, 4 },{ "ROR", (*CPU).ROR, (*CPU).ABS, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		{ "BVS", (*CPU).BVS, (*CPU).REL, 2 },{ "ADC", (*CPU).ADC, (*CPU).IZY, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "ADC", (*CPU).ADC, (*CPU).ZPX, 4 },{ "ROR", (*CPU).ROR, (*CPU).ZPX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "SEI", (*CPU).SEI, (*CPU).IMP, 2 },{ "ADC", (*CPU).ADC, (*CPU).ABY, 4 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "ADC", (*CPU).ADC, (*CPU).ABX, 4 },{ "ROR", (*CPU).ROR, (*CPU).ABX, 7 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "STA", (*CPU).STA, (*CPU).IZX, 6 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "STY", (*CPU).STY, (*CPU).ZP0, 3 },{ "STA", (*CPU).STA, (*CPU).ZP0, 3 },{ "STX", (*CPU).STX, (*CPU).ZP0, 3 },{ "???", (*CPU).XXX, (*CPU).IMP, 3 },{ "DEY", (*CPU).DEY, (*CPU).IMP, 2 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "TXA", (*CPU).TXA, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "STY", (*CPU).STY, (*CPU).ABS, 4 },{ "STA", (*CPU).STA, (*CPU).ABS, 4 },{ "STX", (*CPU).STX, (*CPU).ABS, 4 },{ "???", (*CPU).XXX, (*CPU).IMP, 4 },
		{ "BCC", (*CPU).BCC, (*CPU).REL, 2 },{ "STA", (*CPU).STA, (*CPU).IZY, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "STY", (*CPU).STY, (*CPU).ZPX, 4 },{ "STA", (*CPU).STA, (*CPU).ZPX, 4 },{ "STX", (*CPU).STX, (*CPU).ZPY, 4 },{ "???", (*CPU).XXX, (*CPU).IMP, 4 },{ "TYA", (*CPU).TYA, (*CPU).IMP, 2 },{ "STA", (*CPU).STA, (*CPU).ABY, 5 },{ "TXS", (*CPU).TXS, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "???", (*CPU).NOP, (*CPU).IMP, 5 },{ "STA", (*CPU).STA, (*CPU).ABX, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },
		{ "LDY", (*CPU).LDY, (*CPU).IMM, 2 },{ "LDA", (*CPU).LDA, (*CPU).IZX, 6 },{ "LDX", (*CPU).LDX, (*CPU).IMM, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "LDY", (*CPU).LDY, (*CPU).ZP0, 3 },{ "LDA", (*CPU).LDA, (*CPU).ZP0, 3 },{ "LDX", (*CPU).LDX, (*CPU).ZP0, 3 },{ "???", (*CPU).XXX, (*CPU).IMP, 3 },{ "TAY", (*CPU).TAY, (*CPU).IMP, 2 },{ "LDA", (*CPU).LDA, (*CPU).IMM, 2 },{ "TAX", (*CPU).TAX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "LDY", (*CPU).LDY, (*CPU).ABS, 4 },{ "LDA", (*CPU).LDA, (*CPU).ABS, 4 },{ "LDX", (*CPU).LDX, (*CPU).ABS, 4 },{ "???", (*CPU).XXX, (*CPU).IMP, 4 },
		{ "BCS", (*CPU).BCS, (*CPU).REL, 2 },{ "LDA", (*CPU).LDA, (*CPU).IZY, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "LDY", (*CPU).LDY, (*CPU).ZPX, 4 },{ "LDA", (*CPU).LDA, (*CPU).ZPX, 4 },{ "LDX", (*CPU).LDX, (*CPU).ZPY, 4 },{ "???", (*CPU).XXX, (*CPU).IMP, 4 },{ "CLV", (*CPU).CLV, (*CPU).IMP, 2 },{ "LDA", (*CPU).LDA, (*CPU).ABY, 4 },{ "TSX", (*CPU).TSX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 4 },{ "LDY", (*CPU).LDY, (*CPU).ABX, 4 },{ "LDA", (*CPU).LDA, (*CPU).ABX, 4 },{ "LDX", (*CPU).LDX, (*CPU).ABY, 4 },{ "???", (*CPU).XXX, (*CPU).IMP, 4 },
		{ "CPY", (*CPU).CPY, (*CPU).IMM, 2 },{ "CMP", (*CPU).CMP, (*CPU).IZX, 6 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "CPY", (*CPU).CPY, (*CPU).ZP0, 3 },{ "CMP", (*CPU).CMP, (*CPU).ZP0, 3 },{ "DEC", (*CPU).DEC, (*CPU).ZP0, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "INY", (*CPU).INY, (*CPU).IMP, 2 },{ "CMP", (*CPU).CMP, (*CPU).IMM, 2 },{ "DEX", (*CPU).DEX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "CPY", (*CPU).CPY, (*CPU).ABS, 4 },{ "CMP", (*CPU).CMP, (*CPU).ABS, 4 },{ "DEC", (*CPU).DEC, (*CPU).ABS, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		{ "BNE", (*CPU).BNE, (*CPU).REL, 2 },{ "CMP", (*CPU).CMP, (*CPU).IZY, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "CMP", (*CPU).CMP, (*CPU).ZPX, 4 },{ "DEC", (*CPU).DEC, (*CPU).ZPX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "CLD", (*CPU).CLD, (*CPU).IMP, 2 },{ "CMP", (*CPU).CMP, (*CPU).ABY, 4 },{ "NOP", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "CMP", (*CPU).CMP, (*CPU).ABX, 4 },{ "DEC", (*CPU).DEC, (*CPU).ABX, 7 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
		{ "CPX", (*CPU).CPX, (*CPU).IMM, 2 },{ "SBC", (*CPU).SBC, (*CPU).IZX, 6 },{ "???", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "CPX", (*CPU).CPX, (*CPU).ZP0, 3 },{ "SBC", (*CPU).SBC, (*CPU).ZP0, 3 },{ "INC", (*CPU).INC, (*CPU).ZP0, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 5 },{ "INX", (*CPU).INX, (*CPU).IMP, 2 },{ "SBC", (*CPU).SBC, (*CPU).IMM, 2 },{ "NOP", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).SBC, (*CPU).IMP, 2 },{ "CPX", (*CPU).CPX, (*CPU).ABS, 4 },{ "SBC", (*CPU).SBC, (*CPU).ABS, 4 },{ "INC", (*CPU).INC, (*CPU).ABS, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },
		{ "BEQ", (*CPU).BEQ, (*CPU).REL, 2 },{ "SBC", (*CPU).SBC, (*CPU).IZY, 5 },{ "???", (*CPU).XXX, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 8 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "SBC", (*CPU).SBC, (*CPU).ZPX, 4 },{ "INC", (*CPU).INC, (*CPU).ZPX, 6 },{ "???", (*CPU).XXX, (*CPU).IMP, 6 },{ "SED", (*CPU).SED, (*CPU).IMP, 2 },{ "SBC", (*CPU).SBC, (*CPU).ABY, 4 },{ "NOP", (*CPU).NOP, (*CPU).IMP, 2 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },{ "???", (*CPU).NOP, (*CPU).IMP, 4 },{ "SBC", (*CPU).SBC, (*CPU).ABX, 4 },{ "INC", (*CPU).INC, (*CPU).ABX, 7 },{ "???", (*CPU).XXX, (*CPU).IMP, 7 },
	}
}


func (cpu *CPU) ConnectBus(n *Bus) {
	cpu.bus = n
}


func (cpu *CPU) Read(a uint16) (uint16, error) {
	return cpu.bus.Read(a, false)
}


func (cpu *CPU) Write(a uint16, d uint16) error {
	return cpu.bus.Write(a, d)
}


// Updates the clock cycle
func (cpu *CPU) Clock() {
	if cpu.cycles == 0 {
		cpu.opcode = cpu.Read(cpu.pc)
		cpu.pc++

		cpu.cycles = lookup[opcode].Cycles

		additional_cycle1 := lookup[opcode].AddrMode(cpu)

		additional_cycle2 := lookup[opcode].Operate(cpu)

		cpu.cycles += (additional_cycle1 + additional_cycle2)
	}

	cpu.cycles--
}


func (cpu *CPU) GetFlag(flag int) uint8 {
	if (cpu.status & flag) > 0 {
		return 1
	}
	return 0
}


func (cpu *CPU) SetFlag(flag int, v bool) {
	if v {
		cpu.status |= flag
	}
	else {
		cpu.status &= ^flag
	}
}


// Reset interrupt
func (cpu *CPU) Reset() {
	;
}


// Interrupt request
func (cpu *CPU) IRQ() {
	;
}


// Non-Maskable interrupt request
func (cpu *CPU) NMI() {
	;
}


func (cpu *CPU) Fetch() {
	if !(cpu.lookup[cpu.opcode].AddrMode == (*CPU).IMP) {
		cpu.fetched = cpu.Read(cpu.addr_abs)
	}

	return fetched
}

// Addressing Modes

// Address Mode: Implied
func (cpu *CPU) IMP() uint8 {
	cpu.fetched = cpu.a

	return 0
}

// Address Mode: Immediate
func (cpu *CPU) IMM() uint8 {
	cpu.addr_abs = cpu.pc++

	return 0
}

// Address Mode: Zero Page
func (cpu *CPU) ZP0() uint8 {
	cpu.addr_abs = cpu.Read(cpu.pc)
	cpu.pc++
	addr_abs &= 0x00FF

	return 0
}

// Address Mode: Zero Page with X Offset
func (cpu *CPU) ZPX() uint8 {
	cpu.addr_abs = (cpu.Read(cpu.pc) + cpu.x)
	cpu.pc++
	addr_abs &= 0x00FF

	return 0
}

// Address Mode: Zero Page with Y Offset
func (cpu *CPU) ZPY() uint8 {
	cpu.addr_abs = (cpu.Read(cpu.pc) + cpu.y)
	cpu.pc++
	addr_abs &= 0x00FF

	return 0
}

// Address Mode: Relative
func (cpu *CPU) REL() uint8 {
	cpu.addr_rel = cpu.Read(cpu.pc)
	cpu.pc++
	if cpu.addr_rel & 0x80 {
		cpu.addr_rel |= 0xFF00
	}

	return 0
}

// Address Mode: Absolute
func (cpu *CPU) ABS() uint8 {
	lo := cpu.Read(cpu.pc)
	cpu.pc++
	hi := cpu.Read(cpu.pc)
	cpu.pc++

	cpu.addr_abs = (hi << 8) | lo

	return 0
}

// Address Mode: Absolute with X Offset
func (cpu *CPU) ABX() uint8 {
	lo := cpu.Read(cpu.pc)
	cpu.pc++
	hi := cpu.Read(cpu.pc)
	cpu.pc++

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += cpu.x

	if (cpu.addr_abs & 0xFF00) != (hi << 8) {
		return 1
	}

	return 0
}

// Address Mode: Absolute with Y Offset
func (cpu *CPU) ABY() uint8 {
	lo := cpu.Read(cpu.pc)
	cpu.pc++
	hi := cpu.Read(cpu.pc)
	cpu.pc++

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += cpu.y

	if (cpu.addr_abs & 0xFF00) != (hi << 8) {
		return 1
	}

	return 0
}

// Address Mode: Indirect
func (cpu *CPU) IND() uint8 {
	ptr_lo := cpu.Read(cpu.pc)
	cpu.pc++
	ptr_hi := cpu.Read(cpu.pc)
	cpu.pc++

	ptr := (ptr_hi << 8) | ptr_lo

	if ptr_lo == 0x00FF {
		cpu.addr_abs = (cpu.Read(ptr & 0xFF00) << 8) | cpu.Read(ptr + 0)
	}
	else {
		cpu.addr_abs = (cpu.Read(ptr + 1) << 8) | cpu.Read(ptr + 0)
	}

	return 0
}

// Address Mode: Indirect X
func (cpu *CPU) IZX() uint8 {
	t := cpu.Read(cpu.pc)
	cpu.pc++

	lo := cpu.Read(uint16(cpu.t + uint16(cpu.x)) & 0x00FF)
	hi := cpu.Read(uint16(cpu.t + uint16(cpu.x) + 1) & 0x00FF)

	cpu.addr_abs = (hi << 8) | lo

	return 0
}

// Address Mode: Indirect Y
func (cpu *CPU) IZY() uint8 {
	t := cpu.Read(cpu.pc)
	cpu.pc++

	lo := cpu.Read(t & 0x00FF)
	hi := cpu.Read((t + 1) & 0x00FF)

	cpu.addr_abs = (hi << 8) | lo
	cpu.addr_abs += cpu.y

	if (cpu.addr_abs & 0xFF00) != (hi << 8) {
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

	cpu.SetFlag(V, (~(uint16(cpu.a) ^ uint16(cpu.fetched)) & (uint16(cpu.a) ^ uint16(temp))) & 0x0080)

	cpu.SetFlag(N, temp & 0x80)

	cpu.a = temp & 0x00FF

	return 1
}

// Instruction: Subtraction with Borrow In
func (cpu *CPU) SBC() uint8 {
	cpu.Fetch()

	value := uint16(cpu.fetched) ^ 0x00FF

	temp := uint16(cpu.a) + value + uint16(cpu.GetFlag(C))
	cpu.SetFlag(C, temp & 0xFF00)
	cpu.SetFlag(Z, ((temp & 0x00FF) == 0))
	cpu.SetFlag(V, (temp ^ uint16(cpu.a)) & (temp ^ value) & 0x0080)
	cpu.SetFlag(N, temp & 0x0080)
	cpu.a = temp & 0x00FF

	return 1
}

// Instruction: Bitwise logic AND
func (cpu *CPU) AND() uint8 {
	cpu.Fetch()
	cpu.a = cpu.a & cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, cpu.a & 0x00)

	return 1
}

// Instruction: Bitwise Shift Left
func (cpu *CPU) ASL() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.fetched << 1)
	cpu.SetFlag(C, (temp & 0xFF00) > 0)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
	cpu.SetFlag(N, temp & 0x80)

	if cpu.lookup[cpu.opcode].AddrMode == (*CPU)IMP {
		cpu.a = temp & 0x00FF
	}
	else {
		cpu.Write(cpu.addr_abs, temp & 0x00FF)
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
	cpu.SetFlag(N, cpu.fetched & (1 << 7))
	cpu.SetFlag(V, fetched & (1 << 6))

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

	cpu.SetFlag(I, 1)
	cpu.Write(0x0100 + cpu.stkp, (cpu.pc >> 8) & 0x00FF)
	cpu.stkp--
	cpu.Write(0x0100 + cpu.stkp, cpu.pc & 0x00FF)
	cpu.stkp--

	cpu.SetFlag(B, 1)
	cpu.Write(0x0100 + cpu.stkp, cpu.status)
	cpu.stkp--
	cpu.SetFlag(B, 0)

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
	SetFlag(N, temp & 0x0080)

	return 1
}

// Instruction: Compare X Register
func (cpu *CPU) CPX() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.x) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.x >= cpu.fetched)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	SetFlag(N, temp & 0x0080)

	return 0
}

// Instruction: Comapre Y Register
func (cpu *CPU) CPY() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.y) - uint16(cpu.fetched)
	cpu.SetFlag(C, cpu.y >= cpu.fetched)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	SetFlag(N, temp & 0x0080)

	return 0
}

// Instruction: Decrement Value at Memory Location
func (cpu *CPU) DEC() uint8 {
	cpu.Fetch()
	temp := cpu.fetched - 1
	cpu.Write(cpu.addr_abs, temp & 0x00FF)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, temp & 0x0080)

	return 0
}

// Instruction: Decrement X Register
func (cpu *CPU) DEX() uint8 {
	cpu.x--
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, cpu.x & 0x80)

	return 0
}

// Instruction: Decrement Y Register
func (cpu *CPU) DEY() uint8 {
	cpu.y--
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, cpu.y & 0x80)

	return 0
}

// Instruction: Bitwise Logic XOR
func (cpu *CPU) EOR() uint8 {
	cpu.Fetch()
	cpu.a = cpu.a ^ cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, cpu.a & 0x80)

	return 1
}

// Instruction: Increment Value at Memory Location
func (cpu *CPU) INC() uint8 {
	cpu.Fetch()
	temp := cpu.fetched + 1
	cpu.Write(cpu.addr_abs, temp & 0x00FF)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, temp & 0x0080)

	return 0
}

// Instruction: Increment X Register
func (cpu *CPU) INX() uint8 {
	cpu.x++
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, cpu.x & 0x80)

	return 0
}

// Instruction: Increment Y Register
func (cpu *CPU) INY() uint8 {
	cpu.y++
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, cpu.y & 0x80)

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

	cpu.Write(0x0100 + cpu.stkp, (cpu.pc >> 8) & 0x00FF)
	cpu.stkp--
	cpu.Write(0x0100 + cpu.stkp, cpu.pc & 0x00FF)
	cpu.stkp--

	cpu.pc = cpu.addr_abs

	return 0
}

// Instruction: Load The Accumulator
func (cpu *CPU) LDA() uint8 {
	cpu.Fetch()
	cpu.a = cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, cpu.a & 0x80)

	return 1
}

// Instruction: Load the X Register
func (cpu *CPU) LDX() uint8 {
	cpu.Fetch()
	cpu.x = cpu.fetched
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, cpu.x & 0x80)

	return 1
}

// Instruction: Load the Y Register
func (cpu *CPU) LDY() uint8 {
	cpu.Fetch()
	cpu.y = cpu.fetched
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, cpu.y & 0x80)

	return 1
}

func (cpu *CPU) NOP() uint8 {
	switch cpu.opcode {
	case 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
		return 1
	}
	return 0
}

func (cpu *CPU) ORA() uint8 {
	cpu.Fetch()
	cpu.a = cpu.a | cpu.fetched
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, cpu.a & 0x80)

	return 0
}

func (cpu *CPU) PHA() uint8 {
	cpu.Write(0x0100 + cpu.stkp, cpu.a)
	cpu.stkp--

	return 0
}

func (cpu *CPU) PHP() uint8 {
	cpu.Write(0x0100 + cpu.stkp, cpu.status | B | U)
	cpu.SetFlag(B, 0)
	cpu.SetFlag(U, 0)
	cpu.stkp--

	return 0
}

func (cpu *CPU) PLA() uint8 {
	cpu.stkp++
	cpu.a = cpu.Read(0x0100 + cpu.stkp)
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, cpu.a & 0x80)

	return 0
}

func (cpu *CPU) PLP() uint8 {
	cpu.stkp++
	cpu.status = cpu.Read(0x0100 + cpu.stkp)
	cpu.SetFlag(U, 1)

	return 0
}

func (cpu *CPU) ROL() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.fetched << 1) | cpu.GetFlag(C)
	cpu.SetFlag(C, temp & 0xFF00)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x0000)
	cpu.SetFlag(N, temp & 0x0080)

	if cpu.lookup[cpu.opcode].AddrMode == (*CPU)IMP {
		cpu.a = temp & 0x00FF
	}
	else {
		cpu.Write(cpu.addr_abs, temp & 0x00FF)
	}

	return 0
}

func (cpu *CPU) ROR() uint8 {
	cpu.Fetch()
	temp := uint16(cpu.GetFlag(C) << 7) | (cpu.fetched >> 1)
	cpu.SetFlag(C, cpu.fetched & 0x01)
	cpu.SetFlag(Z, (temp & 0x00FF) == 0x00)
	cpu.SetFlag(N, temp & 0x0080)

	if cpu.lookup[cpu.opcode].AddrMode == (*CPU)IMP {
		cpu.a = temp & 0x00FF
	}
	else {
		cpu.Write(cpu.addr_abs, temp & 0x00FF)
	}

	return 0
}

func (cpu *CPU) RTI() uint8 {
	cpu.stkp++
	cpu.status = cpu.Read(0x0100 + cpu.stkp)
	cpu.status &= ^B
	cpu.status &= ^U

	cpu.stkp++
	cpu.pc = uint16(cpu.Read(0x0100 + cpu.stkp))
	cpu.stkp++
	cpu.pc |= uint16(cpu.Read(0x0100 + cpu.stkp)) << 8

	return 0
}

func (cpu *CPU) RTS() uint8 {
	cpu.stkp++
	cpu.pc = uint16(cpu.Read(0x0100 + cpu.stkp))
	cpu.stkp++
	cpu.pc |= uint16(cpu.Read(0x0100 + cpu.stkp)) << 8

	cpu.pc++

	return 0
}

func (cpu *CPU) SEC() uint8 {
	cpu.SetFlag(C, true)

	return 0
}

func (cpu *CPU) SED() uint8 {
	cpu.SetFlag(D, true)

	return 0
}

func (cpu *CPU) SEI() uint8 {
	cpu.SetFlag(I, true)

	return 0
}

func (cpu *CPU) STA() uint8 {
	cpu.Write(cpu.addr_abs, cpu.a)

	return 0
}

func (cpu *CPU) STX() uint8 {
	cpu.Write(cpu.addr_abs, cpu.x)

	return 0
}

func (cpu *CPU) STY() uint8 {
	cpu.Write(cpu.addr_abs, cpu.y)

	return 0
}

func (cpu *CPU) TAX() uint8 {
	cpu.x = cpu.a
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, cpu.x & 0x80)

	return 0
}

func (cpu *CPU) TAY() uint8 {
	cpu.y = cpu.a
	cpu.SetFlag(Z, cpu.y == 0x00)
	cpu.SetFlag(N, cpu.y & 0x80)

	return 0
}

func (cpu *CPU) TSX() uint8 {
	cpu.x = cpu.stkp
	cpu.SetFlag(Z, cpu.x == 0x00)
	cpu.SetFlag(N, cpu.x & 0x80)

	return 0
}

func (cpu *CPU) TXA() uint8 {
	cpu.a = cpu.x
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, cpu.a & 0x80)

	return 0
}

func (cpu *CPU) TXS() uint8 {
	cpu.stkp = cpu.x

	return 0
}

func (cpu *CPU) TYA() uint8 {
	cpu.a = cpu.y
	cpu.SetFlag(Z, cpu.a == 0x00)
	cpu.SetFlag(N, cpu.a & 0x80)

	return 0
}

func (cpu *CPU) XXX() uint8 {
	return 0
}
