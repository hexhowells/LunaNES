package emu


type Pixel struct {
	R int
	G int
	B int
}


type status struct {
	spriteOverflow bool
	spriteZeroHit bool
	verticalBlank bool
}

func (reg *status) setRegisters(data uint8) {
	reg.spriteOverflow = data & (1<<5) != 0
	reg.spriteZeroHit = data & (1<<6) != 0
	reg.verticalBlank = data & (1<<7) != 0
}

func (reg *status) getRegisters() uint8 {
	return uint8(reg.spriteOverflow) << 5 |
			uint8(reg.spriteZeroHit) << 6 |
			uint8(reg.verticalBlank) << 7
}


type mask struct {
	grayscale bool
	renderBackgroundLeft bool
	renderSpritesLeft bool
	renderBackground bool
	renderSprites bool
	enhanceRed bool
	enhanceGreen bool
	enhanceBlue bool
}

func (reg *mask) setRegisters(data uint8) {
	reg.grayscale = data & (1) != 0
	reg.renderBackgroundLeft = data & (1<<1) != 0
	reg.renderSpritesLeft = data & (1<<2) != 0
	reg.renderBackground = data & (1<<3) != 0
	reg.renderSprites = data & (1<<4) != 0
	reg.enhanceRed = data & (1<<5) != 0
	reg.enhanceGreen = data & (1<<6) != 0
	reg.enhanceBlue = data & (1<<7) != 0
}

func (reg *mask) getRegisters() uint8 {
	return uint8(grayscale) |
			uint8(renderBackgroundLeft) << 1 |
			uint8(renderSpritesLeft) << 2 |
			uint8(renderBackground) << 3 |
			uint8(renderSprites) << 4 |
			uint8(enhanceRed) << 5 |
			uint8(enhanceGreen) << 6 |
			uint8(enhanceBlue) << 7
}


type control struct {
	nametableX bool
	nametableY bool
	incrementMode bool
	patternSprite bool
	patternBackground bool
	spriteSize bool
	slaveMode bool
	enableNmi bool
}

func (reg *control) setRegisters(data uint8) {
	reg.nametableX = data & (1) != 0
	reg.nametableY = data & (1<<1) != 0
	reg.incrementMode = data & (1<<2) != 0
	reg.patternSprite = data & (1<<3) != 0
	reg.patternBackground = data & (1<<4) != 0
	reg.spriteSize = data & (1<<5) != 0
	reg.slaveMode = data & (1<<6) != 0
	reg.enableNmi = data & (1<<7) != 0
}

func (reg *control) getRegisters() uint8 {
	return uint8(nametableX) |
			uint8(nametableY) << 1 |
			uint8(incrementMode) << 2 |
			uint8(patternSprite) << 3 |
			uint8(patternBackground) << 4 |
			uint8(spriteSize) << 5 |
			uint8(slaveMode) << 6 |
			uint8(enableNmi) << 7
}


type PPU struct {
	cart Cartridge

	nameTable [2][1024]uint8
	patternTable[2][4096]uint8
	paletteTable [32]uint8

	scanline int16
	cycle int16
	
	FrameComplete bool

	colourPalette [0x40]Pixel  // stores the colour palettes
	screen [256, 240]Pixel  // stores the pixels to display on the screen
	sprNameTable[2]Sprite  // stores the sprites from the name table
	sprPatternTable[2]Sprite  // stores the sprites from the pattern table

	status status
	mask mask
	control control

	addressLatch uint8  // indicates if high or low byte is being written to
	ppuDataBuffer uint8  // data to ppu is delayed by 1 cycle, so need to buffer the data
	ppuAddress uint16  // stores the compiled address

	Nmi bool
}


func NewPPU() *PPU{
	ppu := PPU{}

	ppu.status = &status{}
	ppu.mask = &mask{}
	ppu.control = &control{}

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


func (p *PPU) GetColourFromPaletteRam(palette uint8, pixel uint8) {
	return p.colourPalette[p.PpuRead(0x3F00 + (palette << 2) + pixel) & 0x3F]
}


func (p *PPU) CpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)

	switch addr {
		case 0x0000:  // control
			break
		case 0x0001:  // mask
			break
		case 0x0002:  // status
			data = (ppu.status.getRegisters() & 0xE0) | (ppu.ppuDataBuffer & 0x1F)
			ppu.status.verticalBlank = false
			ppu.addressLatch = 0
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
			data = ppu.ppuDataBuffer
			ppu.ppuDataBuffer = ppu.PpuRead(ppu.ppuAddress)

			if ppu.ppuAddress > 0x3f00 {
				data = ppu.ppuDataBuffer
			}
			ppu.ppuAddress++
			break
	}

	return data
}


func (p *PPU) CpuWrite(addr uint16, data uint8) {
	switch addr {
		case 0x0000:  // control
			ppu.control.setRegisters(data)
			break
		case 0x0001:  // mask
			ppu.mask.setRegisters(data)
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
			if ppu.addressLatch == 0 {  // store the lower 8 bits of the ppu address
				ppu.ppuAddress = (ppu.ppuAddress 7 0x00FF) | (data << 8)
				ppu.addressLatch = 1
			} else {
				ppu.ppuAddress = (ppu.ppuAddress & 0xFF00) | data
				ppu.addressLatch = 0
			}
			break
		case 0x0007:  // PPU data
			ppu.PpuWrite(ppu.ppuAddress, data)
			ppu.ppuAddress++
			break
	}
}


func (p *PPU) PpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)
	addr &= 0x3FFF

	if p.cart.PpuRead(addr, &data) {
		// cartridge address range
	} else if addr >= 0x0000 && addr <= 0x1FFF {
		data = p.patternTable[(addr & 0x1000) >> 12][addr & 0x0FFF]
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		addr &= 0x0FFF

		if p.cart.mirror == VERTICAL {
			if addr >= 0x0000 && addr <= 0x03FF {data = p.nameTable[0][addr & 0x03FF]}
			if addr >= 0x0400 && addr <= 0x07FF {data = p.nameTable[1][addr & 0x03FF]}
			if addr >= 0x0800 && addr <= 0x0BFF {data = p.nameTable[0][addr & 0x03FF]}
			if addr >= 0x0C00 && addr <= 0x0FFF {data = p.nameTable[1][addr & 0x03FF]}
		} else if p.cart.mirror == HORIZONTAL {
			if addr >= 0x0000 && addr <= 0x03FF {data = p.nameTable[0][addr & 0x03FF]}
			if addr >= 0x0400 && addr <= 0x07FF {data = p.nameTable[0][addr & 0x03FF]}
			if addr >= 0x0800 && addr <= 0x0BFF {data = p.nameTable[1][addr & 0x03FF]}
			if addr >= 0x0C00 && addr <= 0x0FFF {data = p.nameTable[1][addr & 0x03FF]}
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F

		if addr == 0x0010 {addr = 0x0000}
		if addr == 0x0014 {addr = 0x0004}
		if addr == 0x0018 {addr = 0x0008}
		if addr == 0x001C {addr = 0x000C}

		if p.maskRegister.grayscale {
			data = p.paletteTable[addr] & 0x30
		} else {
			data = p.paletteTable[addr] & 0x3F
		}
	}

	return data
}


func (p *PPU) PpuWrite(addr uint16, data uint8) {
	addr &= 0x3FFF

	if p.cart.PpuWrite(addr, data) {
		// cartridge address range
	} else if addr >= 0x0000 && addr <= 0x1FFF {
		p.patternTable[(addr & 0x1000) >> 12][addr & 0x0FFF] = data
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		addr &= 0x0FFF

		if p.cart.mirror == VERTICAL {
			if addr >= 0x0000 && addr <= 0x03FF {p.nameTable[0][addr & 0x03FF] = data}
			if addr >= 0x0400 && addr <= 0x07FF {p.nameTable[1][addr & 0x03FF] = data}
			if addr >= 0x0800 && addr <= 0x0BFF {p.nameTable[0][addr & 0x03FF] = data}
			if addr >= 0x0C00 && addr <= 0x0FFF {p.nameTable[1][addr & 0x03FF] = data}
		} else if p.cart.mirror == HORIZONTAL {
			if addr >= 0x0000 && addr <= 0x03FF {p.nameTable[0][addr & 0x03FF] = data}
			if addr >= 0x0400 && addr <= 0x07FF {p.nameTable[0][addr & 0x03FF] = data}
			if addr >= 0x0800 && addr <= 0x0BFF {p.nameTable[1][addr & 0x03FF] = data}
			if addr >= 0x0C00 && addr <= 0x0FFF {p.nameTable[1][addr & 0x03FF] = data}
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F

		if addr == 0x0010 {addr = 0x0000}
		if addr == 0x0014 {addr = 0x0004}
		if addr == 0x0018 {addr = 0x0008}
		if addr == 0x001C {addr = 0x000C}
		p.paletteTable[addr] = data
	}
}


func (p *PPU) ConnectCartridge(cartridge *Cartridge) {
	p.cart = *cartridge
}


func (p *PPU) Clock() {
	if ppu.scanline == -1 && cycle == 1 {
		ppu.status.verticalBlank = false
	}

	if ppu.scanline == 241 && ppu.cycle == 1 {
		ppu.status.verticalBlank = true
		if ppu.control.enableNmi {
			ppu.Nmi = true
		}
	}

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
