package main

import (
	"LunaNES/emu"
	"log"
)

/*
Runs the nestest ROM to test the correctness of the NES CPU
Docs can be found here: https://www.qmtpro.com/~nes/misc/nestest.txt
*/
func main() {
	cpu := emu.NewCPU()
	bus := emu.NewBus()
	cart := emu.NewCartridge("../ROMS/nestest.nes")

	if cart == nil {
        log.Println("Error: cartridge not loaded correctly")
    }

	bus.InsertCartridge(cart)

	// Set the reset vector
	bus.CpuWrite(0xFFFC, 0x00)
	bus.CpuWrite(0xFFFD, 0x80)

	cpu.ConnectBus(bus)

	cpu.Reset()

	cpu.Pc = 0xC000  // run the test rom in automation mode

	// Run the program for 10 million clock cycles
	new_inst := false
	for i := 0; i < 10_000_000; {
		new_inst = cpu.Clock()
		if new_inst {
			i++
		}
	}

	// print final state of the CPU
	// RAM addresses 0x02 & 0x03 should both have the value 0x00 to indicate all tests passed successfully
	// any other value indicates a specific error state which is detailed in the docs
	cpu.PrintCPU()
	cpu.PrintStatusFlags()
	cpu.PrintRAM(0x00, 1)
}