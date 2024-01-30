package emu


type Pixel struct {
	R int
	G int
	B int
}


type PPU struct {
	cart Cartridge
	nameTable [2][1024]uint8
	patternTable[2][4096]uint8
	paletteTable [32]uint8
	scanline int16
	cycle int16
	frameComplete bool

	colourPalette [0x40]Pixel  // stores the colour palettes
	screen [256, 240]Pixel  // stores the pixels to display on the screen
	sprNameTable[2]Sprite  // stores the sprites from the name table
	sprPatternTable[2]Sprite  // stores the sprites from the pattern table
}


func NewPPU() *PPU{
	ppu := PPU{}

	ppu.colourPalette[0x00] = Pixel{84, 84, 84}
	ppu.colourPalette[0x00] = Pixel{84, 84, 84}
	ppu.colourPalette[0x01] = Pixel{0, 30, 116}
	ppu.colourPalette[0x02] = Pixel{8, 16, 144}
	ppu.colourPalette[0x03] = Pixel{48, 0, 136}
	ppu.colourPalette[0x04] = Pixel{68, 0, 100}
	ppu.colourPalette[0x05] = Pixel{92, 0, 48}
	ppu.colourPalette[0x06] = Pixel{84, 4, 0}
	ppu.colourPalette[0x07] = Pixel{60, 24, 0}
	ppu.colourPalette[0x08] = Pixel{32, 42, 0}
	ppu.colourPalette[0x09] = Pixel{8, 58, 0}
	ppu.colourPalette[0x0A] = Pixel{0, 64, 0}
	ppu.colourPalette[0x0B] = Pixel{0, 60, 0}
	ppu.colourPalette[0x0C] = Pixel{0, 50, 60}
	ppu.colourPalette[0x0D] = Pixel{0, 0, 0}
	ppu.colourPalette[0x0E] = Pixel{0, 0, 0}
	ppu.colourPalette[0x0F] = Pixel{0, 0, 0}

	ppu.colourPalette[0x10] = Pixel{152, 150, 152}
	ppu.colourPalette[0x11] = Pixel{8, 76, 196}
	ppu.colourPalette[0x12] = Pixel{48, 50, 236}
	ppu.colourPalette[0x13] = Pixel{92, 30, 228}
	ppu.colourPalette[0x14] = Pixel{136, 20, 176}
	ppu.colourPalette[0x15] = Pixel{160, 20, 100}
	ppu.colourPalette[0x16] = Pixel{152, 34, 32}
	ppu.colourPalette[0x17] = Pixel{120, 60, 0}
	ppu.colourPalette[0x18] = Pixel{84, 90, 0}
	ppu.colourPalette[0x19] = Pixel{40, 114, 0}
	ppu.colourPalette[0x1A] = Pixel{8, 124, 0}
	ppu.colourPalette[0x1B] = Pixel{0, 118, 40}
	ppu.colourPalette[0x1C] = Pixel{0, 102, 120}
	ppu.colourPalette[0x1D] = Pixel{0, 0, 0}
	ppu.colourPalette[0x1E] = Pixel{0, 0, 0}
	ppu.colourPalette[0x1F] = Pixel{0, 0, 0}

	ppu.colourPalette[0x20] = Pixel{236, 238, 236}
	ppu.colourPalette[0x21] = Pixel{76, 154, 236}
	ppu.colourPalette[0x22] = Pixel{120, 124, 236}
	ppu.colourPalette[0x23] = Pixel{176, 98, 236}
	ppu.colourPalette[0x24] = Pixel{228, 84, 236}
	ppu.colourPalette[0x25] = Pixel{236, 88, 180}
	ppu.colourPalette[0x26] = Pixel{236, 106, 100}
	ppu.colourPalette[0x27] = Pixel{212, 136, 32}
	ppu.colourPalette[0x28] = Pixel{160, 170, 0}
	ppu.colourPalette[0x29] = Pixel{116, 196, 0}
	ppu.colourPalette[0x2A] = Pixel{76, 208, 32}
	ppu.colourPalette[0x2B] = Pixel{56, 204, 108}
	ppu.colourPalette[0x2C] = Pixel{56, 180, 204}
	ppu.colourPalette[0x2D] = Pixel{60, 60, 60}
	ppu.colourPalette[0x2E] = Pixel{0, 0, 0}
	ppu.colourPalette[0x2F] = Pixel{0, 0, 0}

	ppu.colourPalette[0x30] = Pixel{236, 238, 236}
	ppu.colourPalette[0x31] = Pixel{168, 204, 236}
	ppu.colourPalette[0x32] = Pixel{188, 188, 236}
	ppu.colourPalette[0x33] = Pixel{212, 178, 236}
	ppu.colourPalette[0x34] = Pixel{236, 174, 236}
	ppu.colourPalette[0x35] = Pixel{236, 174, 212}
	ppu.colourPalette[0x36] = Pixel{236, 180, 176}
	ppu.colourPalette[0x37] = Pixel{228, 196, 144}
	ppu.colourPalette[0x38] = Pixel{204, 210, 120}
	ppu.colourPalette[0x39] = Pixel{180, 222, 120}
	ppu.colourPalette[0x3A] = Pixel{168, 226, 144}
	ppu.colourPalette[0x3B] = Pixel{152, 226, 180}
	ppu.colourPalette[0x3C] = Pixel{160, 214, 228}
	ppu.colourPalette[0x3D] = Pixel{160, 162, 160}
	ppu.colourPalette[0x3E] = Pixel{0, 0, 0}
	ppu.colourPalette[0x3F] = Pixel{0, 0, 0}

	return &ppu
}


func (p *PPU) GetPatternTable(i uint8, palette uint8) {
	// Loop through all 16x16 tiles
	for nTileY := 0; nTileY < 16; nTileY++ {
		for nTileX := 0; nTileX < 16; nTileX++ {
			nOffset := uint8(nTileY * 256 + nTileX * 16)

			// Loop through 8x8 grid of pixels per tile
			// And set each pixel value for the tile
			for row := 0; row < 8; row++ {
				tileLsb := p.PpuRead(i * 0x1000 + nOffset + row + 0x0000)
				tileMsb := p.PpuRead(i * 0x1000 + nOffset + row + 0x0008)

				for col := 0; col < 8; col++ {
					pixel := (tileLsb & 0x01) + (tileMsb & 0x01)

					tileLsb >>= 1
					tileMsb >>= 1

					// Set the pixel value of the tile
					p.sprPatternTable[i].SetPixel(
						nTileY * 8 + row, 
						nTileX * 8 + (7 - col), 
						p.GetColourFromPaletteRam(palette, pixel)
					)
				}
			}
		}
	}
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
