package emu

type PPU struct {
	cart Cartridge
	nameTable [2][1024]uint8
	paletteTable [32]uint8
	scanline int16
	cycle int16
	frameComplete bool
}


func (p *PPU) CpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)

	switch addr {
		case 0x0000:  // control
			break
		case 0x0001:  // mask
			break
		case 0x0002:  // status
			break
		case 0x0003:  // OAM address
			break
		case 0x0004:  // OAM data
			break
		case 0x0005:  // scroll
			break
		case 0x0006:  // PPU address
			break
		case 0x0007:  // PPU data
			break
	}

	return data
}


func (p *PPU) CpuWrite(addr uint16, data uint8) {
	switch addr {
		case 0x0000:  // control
			break
		case 0x0001:  // mask
			break
		case 0x0002:  // status
			break
		case 0x0003:  // OAM address
			break
		case 0x0004:  // OAM data
			break
		case 0x0005:  // scroll
			break
		case 0x0006:  // PPU address
			break
		case 0x0007:  // PPU data
			break
	}
}


func (p *PPU) PpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)
	addr &= 0x3FFF

	if p.cart.PpuRead(addr, &data) {
		// cartridge address range
	}

	return data
}


func (p *PPU) PpuWrite(addr uint16, data uint8) {
	addr &= 0x3FFF

	if p.cart.PpuWrite(addr, data) {
		// cartridge address range
	}
}


func (p *PPU) ConnectCartridge(cartridge *Cartridge) {
	p.cart = *cartridge
}


func (p *PPU) Clock() {
	p.cycle++

	if p.cycle >= 341 {
		p.cycle = 0
		p.scanline++

		if p.scanline >= 261 {
			p.scanline = -1
			p.frameComplete = true
		}
	}
}
