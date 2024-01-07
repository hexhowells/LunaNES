package emu

import (
	"fmt"
)


type Bus struct {
	cpuRam [0x1FFF + 1]uint8
	cpu CPU
	ppu PPU
	nSystemClockCounter uint32  // count of how many clock cycles have passed
}


func NewBus() *Bus {
	bus := Bus{}
	bus.cpu.ConnectBus(&bus)
	bus.nSystemClockCounter = 0

	for i := range bus.cpuRam {
		bus.cpuRam[i] = 0x00
	}

	return &bus
}


func (b *Bus) Reset() {
	b.cpu.reset()
	b.nSystemClockCounter = 0
}


func (b *Bus) Clock() {
	b.ppu.Clock()

	if b.nSystemClockCounter % 3 == 0 {
		b.cpu.Clock()
	}

	b.nSystemClockCounter++
}


func (b *Bus) WriteBytes(addr uint16, data []uint8) {
	for i, byteData := range data {
		addr := addr + uint16(i)
		if addr > 0x1FFF {
			break // Stop writing if we reach the end of RAM
		}
		b.CpuWrite(addr & 0x07FF, byteData)
	}
}


func (b *Bus) CpuWrite(addr uint16, data uint8) {
	if addr >= 0x0000 && addr <= 0x1FFF {
		b.cpuRam[addr & 0x07FF] = data
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		ppu.CpuWrite(addr & 0x07FF, data)
	}
}


func (b *Bus) CpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)

	if addr >= 0x0000 && addr <= 0x1FFF {
		data = b.cpuRam[addr & 0x07FF]
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		ppu.CpuRead(addr & 0x07FF, bReadOnly)
	}

	return data
}


func (b *Bus) PrintRAM(startPage int, pages int) {
	const bytesPerRow = 16

	startPage = 16 * 16 * startPage

	if pages == 0 {
		pages = len(b.cpuRam)
	} else {
		pages = (16 * 16 * pages) + startPage
	}

	fmt.Println("\nAddress  | 00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F | ASCII")
	fmt.Println("---------+-------------------------------------------------+-----------------")
	for i := startPage; i < pages; i += bytesPerRow {
		// Print address
		fmt.Printf("%04X     |", i)
		ascii := ""
		for j := 0; j < bytesPerRow; j++ {
			// Print byte
			fmt.Printf(" %02X", b.cpuRam[i+j])
			// Collect ASCII representation, if printable
			if b.cpuRam[i+j] >= 0x20 && b.cpuRam[i+j] <= 0x7E {
				ascii += string(b.cpuRam[i+j])
			} else {
				ascii += "."
			}
		}
		fmt.Printf(" | %s\n", ascii)
	}
}
