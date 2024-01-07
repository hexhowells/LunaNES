package emu

type PPU struct {

}


func (p *PPU) cpuRead(addr uint16, bReadOnly bool) uint8 {
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


func (p *PPU) cpuWrite(addr uint16, data uint8) {
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


func (p *PPU) ppuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)
	addr &= 0x3FFF

	return data
}


func (p *PPU) ppuWrite(addr uint16, data uint8) {
	addr &= 0x3FFF
}
