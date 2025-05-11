package emu

import (
	"fmt"
)


type Bus struct {
	cpuRam [0x1FFF + 1]uint8
	Cpu CPU
	Ppu PPU
	Controller [2]uint8
	cart Cartridge
	nSystemClockCounter uint32  // count of how many clock cycles have passed
	controllerState[2] uint8
	dmaPage uint8
	dmaAddr uint8
	dmaData uint8
	dmaTransfer bool  // flag indicating if DMA is happening
	dmaDummy bool
}


func NewBus() *Bus {
	bus := Bus{}
	
	bus.Cpu = *NewCPU()
	bus.Cpu.ConnectBus(&bus)
	
	bus.Ppu = *NewPPU()
	
	bus.nSystemClockCounter = 0

	for i := range bus.cpuRam {
		bus.cpuRam[i] = 0x00
	}

	return &bus
}


func (b *Bus) Reset() {
	b.Ppu.Reset()
	b.Cpu.Reset()
	b.nSystemClockCounter = 0
	b.dmaPage = 0x00
	b.dmaAddr = 0x00
	b.dmaData = 0x00
	b.dmaDummy = true
	b.dmaTransfer = false
}


func (b *Bus) Clock() {
	b.Ppu.Clock()

	// clock CPU 3 times slower then PPU
	if b.nSystemClockCounter % 3 == 0 {
		// lock CPU during DMA transfer operation
		if b.dmaTransfer {
			if b.dmaDummy {  // wait for correct clock cycle to begin DMA transfer
				if b.nSystemClockCounter % 2 == 1 {
					b.dmaDummy = false
				}
			} else {
				if b.nSystemClockCounter % 2 == 0 { // even clock cycles
					addr := (uint16(b.dmaPage) << 8) | uint16(b.dmaAddr)
					b.dmaData = b.CpuRead(addr, true)
				} else {  // odd clock cycles
					b.Ppu.Oam[b.dmaAddr] = b.dmaData
					b.dmaAddr++

					if b.dmaAddr == 0x00 {  // DMA transfer complete
						b.dmaTransfer = false
						b.dmaDummy = true
					}
				}
			}
		} else {  // clock CPU if DMA transfer is not taking place
			b.Cpu.Clock()
		}
	}

	if b.Ppu.Nmi {
		b.Ppu.Nmi = false
		b.Cpu.NMI()
	}

	b.nSystemClockCounter++
}


func (b *Bus) InsertCartridge(cartridge *Cartridge) {
	b.cart = *cartridge
	b.Ppu.ConnectCartridge(cartridge)
}


// Writes a chunk of bytes to the bus
func (b *Bus) WriteBytes(addr uint16, data []uint8) {
	for i, byteData := range data {
		addr := addr + uint16(i)
		if addr > 0x1FFF {
			break // Stop writing if we reach the end of RAM
		}
		b.CpuWrite(addr, byteData)
	}
}


// All write operations sent out to the bus get processed here
func (b *Bus) CpuWrite(addr uint16, data uint8) {
	if b.cart.CpuWrite(addr, data) {
		// cartridge address range
	} else if addr >= 0x0000 && addr <= 0x1FFF {
		b.cpuRam[addr & 0x07FF] = data
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		b.Ppu.CpuWrite(addr & 0x0007, data)
	} else if addr == 0x4014 {
		b.dmaPage = data
		b.dmaAddr = 0x00
		b.dmaTransfer = true
	} else if addr >= 0x4016 && addr <= 0x4017 {
		b.controllerState[addr & 0x0001] = b.Controller[addr & 0x0001]
	}
}


// All read operations sent out to the bus get processed here
func (b *Bus) CpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)

	if b.cart.CpuRead(addr, &data) {
		// cartridge address range
	} else if addr >= 0x0000 && addr <= 0x1FFF {
		data = b.cpuRam[addr & 0x07FF]
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		data = b.Ppu.CpuRead(addr & 0x0007, bReadOnly)
	} else if addr >= 0x4016 && addr <= 0x4017 {
		if (b.controllerState[addr & 0x0001] & 0x80) != 0 {
		    data = 1
		} else {
		    data = 0
		}
		b.controllerState[addr & 0x0001] = b.controllerState[addr & 0x0001] << 1
	}

	return data
}


// Prints out the CPU RAM
// RAM get printed out in pages
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
			fmt.Printf(" %02X", b.CpuRead(uint16(i+j), true))
			// Collect ASCII representation, if printable
			if b.CpuRead(uint16(i+j), true) >= 0x20 && b.CpuRead(uint16(i+j), true) <= 0x7E {
				ascii += string(b.CpuRead(uint16(i+j), true))
			} else {
				ascii += "."
			}
		}
		fmt.Printf(" | %s\n", ascii)
	}
}
