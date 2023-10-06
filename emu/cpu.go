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
	Operate func(cpu) uint8
	AddrMode func(cpu) uint8
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

	lookup := []INSTRUCTION{
		{}
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
	;
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
	;
}
