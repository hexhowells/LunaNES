package emu

import (
	"fmt"
)


type Bus struct {
	ram [0xFFFF + 1]uint8
	cpu CPU
}


func NewBus() *Bus {
	bus := Bus{}
	bus.cpu.ConnectBus(&bus)

	for i := range bus.ram {
		bus.ram[i] = 0x00
	}

	return &bus
}


func (b *Bus) Write(addr uint16, data uint8) {
	if addr >= 0x0000 && addr <= 0xFFFF {
		b.ram[addr] = data
	}
}


func (b *Bus) Read(addr uint16, bReadOnly bool) uint8 {
	if addr >= 0x0000 && addr <= 0xFFFF {
		return b.ram[addr]
	}

	return 0x00
}


func (b *Bus) PrintRAM() {
	const bytesPerRow = 16
	fmt.Println("Address  | 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F | ASCII")
	fmt.Println("---------+-------------------------------------------------+-----------------")
	for i := 0; i < len(b.ram); i += bytesPerRow {
		// Print address
		fmt.Printf("%04X     |", i)
		ascii := ""
		for j := 0; j < bytesPerRow; j++ {
			// Print byte
			fmt.Printf(" %02X", b.ram[i+j])
			// Collect ASCII representation, if printable
			if b.ram[i+j] >= 0x20 && b.ram[i+j] <= 0x7E {
				ascii += string(b.ram[i+j])
			} else {
				ascii += "."
			}
		}
		fmt.Printf(" | %s\n", ascii)
	}
}
