package emu


type Pixel struct {
	R uint8
	G uint8
	B uint8
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
	return Btoi(reg.spriteOverflow) << 5 |
			Btoi(reg.spriteZeroHit) << 6 |
			Btoi(reg.verticalBlank) << 7
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
	return Btoi(reg.grayscale) |
			Btoi(reg.renderBackgroundLeft) << 1 |
			Btoi(reg.renderSpritesLeft) << 2 |
			Btoi(reg.renderBackground) << 3 |
			Btoi(reg.renderSprites) << 4 |
			Btoi(reg.enhanceRed) << 5 |
			Btoi(reg.enhanceGreen) << 6 |
			Btoi(reg.enhanceBlue) << 7
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
	return Btoi(reg.nametableX) |
			Btoi(reg.nametableY) << 1 |
			Btoi(reg.incrementMode) << 2 |
			Btoi(reg.patternSprite) << 3 |
			Btoi(reg.patternBackground) << 4 |
			Btoi(reg.spriteSize) << 5 |
			Btoi(reg.slaveMode) << 6 |
			Btoi(reg.enableNmi) << 7
}

type loopyRegister struct {
	coarseX uint16  // 5 bits
	coarseY uint16  // 5 bits
	nametableX bool  // 1 bit
	nametableY bool  // 1 bit
	fineY uint16  // 3 bits
	unused bool  // 1 bit
}

func (lr *loopyRegister) SetRegisters(value uint16) {
	lr.coarseX = value & 0x001F                  // bits 0–4
	lr.coarseY = (value >> 5) & 0x001F           // bits 5–9
	lr.nametableX = (value>>10)&1 != 0           // bit 10
	lr.nametableY = (value>>11)&1 != 0           // bit 11
	lr.fineY = (value >> 12) & 0x0007            // bits 12–14
	lr.unused = (value>>15)&1 != 0               // bit 15
}

func (lr *loopyRegister) GetRegisters() uint16 {
	return lr.coarseX & 0x001F |
			(lr.coarseY & 0x001F) << 5 |
			Btoi16(lr.nametableX) << 10 |
			Btoi16(lr.nametableY) << 11 |
			(lr.fineY & 0x0007) << 12 |
			Btoi16(lr.unused) << 15
}


type sObjectAttributeEntry struct {
	y uint8
	id uint8
	attribute uint8
	x uint8
}


func (p *PPU) GetOAMEntry(index int) sObjectAttributeEntry {
	base := index * 4
	return sObjectAttributeEntry{
		y:         p.Oam[base+0],
		id:        p.Oam[base+1],
		attribute: p.Oam[base+2],
		x:         p.Oam[base+3],
	}
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
	screen [256][240]Pixel  // stores the pixels to display on the screen
	sprNameTable[2]Sprite  // stores the sprites from the name table
	sprPatternTable[2]Sprite  // stores the sprites from the pattern table

	status *status
	mask *mask
	control *control

	Oam [256]uint8
	OamAddr uint8

	spriteScanLine [8]sObjectAttributeEntry
	spriteCount uint8

	vramAddr *loopyRegister
	tramAddr *loopyRegister

	fineX uint8

	addressLatch uint8  // indicates if high or low byte is being written to
	ppuDataBuffer uint8  // data to ppu is delayed by 1 cycle, so need to buffer the data

	Nmi bool

	bgNextTileID      uint8
	bgNextTileAttrib  uint8
	bgNextTileLsb     uint8
	bgNextTileMsb     uint8

	bgShifterPatternLo uint16
	bgShifterPatternHi uint16
	bgShifterAttribLo  uint16
	bgShifterAttribHi  uint16

	spriteShifterPatternLo [8]uint8
	spriteShifterPatternHi [8]uint8
	bSpriteZeroHitPossible bool
	bSpriteZeroBeingRendered bool
}


func NewPPU() *PPU{
	ppu := PPU{}

	ppu.status = &status{}
	ppu.mask = &mask{}
	ppu.control = &control{}
	ppu.vramAddr = &loopyRegister{}
	ppu.tramAddr = &loopyRegister{}

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
				tileLsb := p.PpuRead(uint16(i) * 0x1000 + uint16(nOffset) + uint16(row) + 0x0000, true)
				tileMsb := p.PpuRead(uint16(i) * 0x1000 + uint16(nOffset) + uint16(row) + 0x0008, true)

				for col := 0; col < 8; col++ {
					pixel := (tileLsb & 0x01) + (tileMsb & 0x01)

					tileLsb >>= 1
					tileMsb >>= 1

					// Set the pixel value of the tile
					p.sprPatternTable[i].SetPixel(
						uint8(nTileX * 8 + (7 - col)), 
						uint8(nTileY * 8 + row),
						p.GetColourFromPaletteRam(palette, pixel),
					)
				}
			}
		}
	}
}


func (p *PPU) BackgroundPatternTableBase() uint16 {
	if p.control.patternBackground {
		return 0x1000
	}
	return 0x0000
}


func (p *PPU) GetColourFromPaletteRam(palette uint8, pixel uint8) Pixel {
	return p.colourPalette[p.PpuRead(0x3F00 + uint16(palette << 2) + uint16(pixel), true) & 0x3F]
}


func (p *PPU) CpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)

	switch addr {
		case 0x0000:  // control
			break
		case 0x0001:  // mask
			break
		case 0x0002:  // status
			data = (p.status.getRegisters() & 0xE0) | (p.ppuDataBuffer & 0x1F)
			p.status.verticalBlank = false
			p.addressLatch = 0
			break
		case 0x0003:  // OAM address
			break
		case 0x0004:  // OAM data
			data = p.Oam[p.OamAddr]
			break
		case 0x0005:  // scroll
			break
		case 0x0006:  // PPU address
			break
		case 0x0007:  // PPU data
			data = p.ppuDataBuffer
			p.ppuDataBuffer = p.PpuRead(p.vramAddr.GetRegisters(), true)

			if p.vramAddr.GetRegisters() > 0x3f00 {
				data = p.ppuDataBuffer
			}
			if p.control.incrementMode {
				p.vramAddr.SetRegisters(p.vramAddr.GetRegisters() + 32)
			} else {
				p.vramAddr.SetRegisters(p.vramAddr.GetRegisters() + 1)
			}
			break
	}

	return data
}


func (p *PPU) CpuWrite(addr uint16, data uint8) {
	switch addr {
		case 0x0000:  // control
			p.control.setRegisters(data)
			p.tramAddr.nametableX = p.control.nametableX
			p.tramAddr.nametableY = p.control.nametableY
			break
		case 0x0001:  // mask
			p.mask.setRegisters(data)
			break
		case 0x0002:  // status
			break
		case 0x0003:  // OAM address
			p.OamAddr = data
			break
		case 0x0004:  // OAM data
			p.Oam[p.OamAddr] = data
			p.OamAddr = (p.OamAddr + 1) & 0xFF
			break
		case 0x0005:  // scroll
			if p.addressLatch == 0 {
				p.fineX = data & 0x07
				p.tramAddr.coarseX = uint16(data >> 3)
				p.addressLatch = 1
			} else {
				p.tramAddr.fineY = uint16(data & 0x07)
				p.tramAddr.coarseY = uint16(data >> 3)
				p.addressLatch = 0
			}
			break
		case 0x0006:  // PPU address
			if p.addressLatch == 0 {  // store the lower 8 bits of the ppu address
				p.tramAddr.SetRegisters((p.tramAddr.GetRegisters() & 0x00FF) | (uint16(data) << 8))
				p.addressLatch = 1
			} else {
				p.tramAddr.SetRegisters((p.tramAddr.GetRegisters() & 0xFF00) | uint16(data))
				*p.vramAddr = *p.tramAddr
				p.addressLatch = 0
			}
			break
		case 0x0007:  // PPU data
			p.PpuWrite(p.vramAddr.GetRegisters(), data)
			
			if p.control.incrementMode {
				p.vramAddr.SetRegisters(p.vramAddr.GetRegisters() + 32)
			} else {
				p.vramAddr.SetRegisters(p.vramAddr.GetRegisters() + 1)
			}
			break
	}
}


func (p *PPU) PpuRead(addr uint16, bReadOnly bool) uint8 {
	data := uint8(0x00)
	addr &= 0x3FFF

	if p.cart.PpuRead(addr, &data) {
		// cartridge address range
	} else if addr >= 0x0000 && addr <= 0x1FFF {  // pattern table
		data = p.patternTable[(addr & 0x1000) >> 12][addr & 0x0FFF]
	} else if addr >= 0x2000 && addr <= 0x3EFF {  // nametable
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
	} else if addr >= 0x3F00 && addr <= 0x3FFF { // palette memory
		addr &= 0x001F

		if addr == 0x0010 {addr = 0x0000}
		if addr == 0x0014 {addr = 0x0004}
		if addr == 0x0018 {addr = 0x0008}
		if addr == 0x001C {addr = 0x000C}

		if p.mask.grayscale {
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
	} else if addr >= 0x0000 && addr <= 0x1FFF { // pattern table
		p.patternTable[(addr & 0x1000) >> 12][addr & 0x0FFF] = data
	} else if addr >= 0x2000 && addr <= 0x3EFF { // nametable
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
	} else if addr >= 0x3F00 && addr <= 0x3FFF { // palette memory
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


func (p *PPU) Reset() {
	p.scanline = 0
	p.cycle = 0
	p.FrameComplete = false
	p.Nmi = false

	// Reset PPU registers
	p.status = &status{}
	p.mask = &mask{}
	p.control = &control{}
	p.vramAddr = &loopyRegister{}
	p.tramAddr = &loopyRegister{}
	p.OamAddr = 0x00

	p.addressLatch = 0x00
	p.ppuDataBuffer = 0x00
}


func (p *PPU) Screen() *[256][240]Pixel {
	return &p.screen
}


func (p *PPU) Clock() {

	//--------------------------------------------------------------------
	// inline helpers
	//--------------------------------------------------------------------

	// increment scroll horizontally 1 tile (coarse X)
	incrementScrollX := func() {
		if p.mask.renderBackground || p.mask.renderSprites {
			if p.vramAddr.coarseX == 31 {
				p.vramAddr.coarseX = 0
				p.vramAddr.nametableX = !p.vramAddr.nametableX
			} else {
				p.vramAddr.coarseX++
			}
		}
	}

	// increment scroll vertically 1 scanline (fine Y + coarse Y)
	incrementScrollY := func() {
		if p.mask.renderBackground || p.mask.renderSprites {
			if p.vramAddr.fineY < 7 {
				p.vramAddr.fineY++
			} else {
				p.vramAddr.fineY = 0
				switch p.vramAddr.coarseY {
				case 29:
					p.vramAddr.coarseY = 0
					p.vramAddr.nametableY = !p.vramAddr.nametableY
				case 31:
					p.vramAddr.coarseY = 0
				default:
					p.vramAddr.coarseY++
				}
			}
		}
	}

	transferAddressX := func() {
		if p.mask.renderBackground || p.mask.renderSprites {
			p.vramAddr.coarseX = p.tramAddr.coarseX
			p.vramAddr.nametableX = p.tramAddr.nametableX
		}
	}
	transferAddressY := func() {
		if p.mask.renderBackground || p.mask.renderSprites {
			p.vramAddr.fineY      = p.tramAddr.fineY
			p.vramAddr.coarseY    = p.tramAddr.coarseY
			p.vramAddr.nametableY = p.tramAddr.nametableY
		}
	}

	loadBackgroundShifters := func() {
		p.bgShifterPatternLo = (p.bgShifterPatternLo & 0xFF00) | uint16(p.bgNextTileLsb)
		p.bgShifterPatternHi = (p.bgShifterPatternHi & 0xFF00) | uint16(p.bgNextTileMsb)
		p.bgShifterAttribLo  = (p.bgShifterAttribLo & 0xFF00) | func() uint16 {
			if (p.bgNextTileAttrib & 0x01) != 0 { return 0xFF }
			return 0
		}()
		p.bgShifterAttribHi  = (p.bgShifterAttribHi & 0xFF00) | func() uint16 {
			if (p.bgNextTileAttrib & 0x02) != 0 { return 0xFF }
			return 0
		}()
	}

	updateShifters := func() {
		if p.mask.renderBackground {
			p.bgShifterPatternLo <<= 1
			p.bgShifterPatternHi <<= 1
			p.bgShifterAttribLo  <<= 1
			p.bgShifterAttribHi  <<= 1
		}

		if p.mask.renderSprites && p.cycle >= 1 && p.cycle < 258 {
			for i := 0; i < int(p.spriteCount); i++ {
				if p.spriteScanLine[i].x > 0 {
					p.spriteScanLine[i].x--
				} else {
					p.spriteShifterPatternLo[i] <<= 1
					p.spriteShifterPatternHi[i] <<= 1
				}
			}
		}
	}

	//--------------------------------------------------------------------
	// Visible scanlines + pre‑render (-1 .. 239)
	//--------------------------------------------------------------------
	if p.scanline >= -1 && p.scanline < 240 {

		// odd‑frame idle‑cycle skip
		if p.scanline == 0 && p.cycle == 0 {
			p.cycle = 1
		}

		// clear VBlank at first pre‑render line
		if p.scanline == -1 && p.cycle == 1 {
			p.status.verticalBlank = false
		}

		// ----------------------------------------------------------------
		// Background pipeline – fetches & shifters
		// ----------------------------------------------------------------
		if (p.cycle >= 2 && p.cycle < 258) || (p.cycle >= 321 && p.cycle < 338) {

			updateShifters()

			switch (p.cycle - 1) % 8 {
			case 0:
				loadBackgroundShifters()
				p.bgNextTileID = p.PpuRead(0x2000|(p.vramAddr.GetRegisters()&0x0FFF), true)

			case 2:
				attrAddr := uint16(0x23C0 |
					(Btoi16(p.vramAddr.nametableY) << 11) |
					(Btoi16(p.vramAddr.nametableX) << 10) |
					((p.vramAddr.coarseY >> 2) << 3) |
					(p.vramAddr.coarseX >> 2))
				attr := p.PpuRead(attrAddr, true)
				if p.vramAddr.coarseY&0x02 != 0 { attr >>= 4 }
				if p.vramAddr.coarseX&0x02 != 0 { attr >>= 2 }
				p.bgNextTileAttrib = attr & 0x03

			case 4:
				base := p.BackgroundPatternTableBase()
				p.bgNextTileLsb = p.PpuRead(base+uint16(p.bgNextTileID)*16+uint16(p.vramAddr.fineY), true)

			case 6:
				base := p.BackgroundPatternTableBase()
				p.bgNextTileMsb = p.PpuRead(base+uint16(p.bgNextTileID)*16+uint16(p.vramAddr.fineY)+8, true)

			case 7:
				incrementScrollX()
			}
		}

		if p.cycle == 256 {
			incrementScrollY()
		}
		if p.cycle == 257 {
			loadBackgroundShifters()
			transferAddressX()

			// Sprite evaluation initialise
			for i := range p.spriteScanLine {
				p.spriteScanLine[i] = sObjectAttributeEntry{}
				p.spriteShifterPatternLo[i] = 0
				p.spriteShifterPatternHi[i] = 0
			}
		}

		if p.cycle == 338 || p.cycle == 340 {
			p.bgNextTileID = p.PpuRead(0x2000|(p.vramAddr.GetRegisters()&0x0FFF), true)
		}

		if p.scanline == -1 && p.cycle >= 280 && p.cycle < 305 {
			transferAddressY()
		}

		// ----------------------------------------------------------------
		// Sprite evaluation (cycle 257, visible scanlines only)
		// ----------------------------------------------------------------
		if p.cycle == 257 && p.scanline >= 0 {

			p.spriteCount = 0
			p.bSpriteZeroHitPossible = false
			for i := range p.spriteShifterPatternLo {
				p.spriteShifterPatternLo[i] = 0
				p.spriteShifterPatternHi[i] = 0
			}

			height := int16(8)
			if p.control.spriteSize {
				height = 16
			}

			for o := 0; o < 64 && p.spriteCount < 8; o++ {
				y := int16(p.Oam[o*4+0])
				diff := int16(p.scanline) - y - 1
				if diff >= 0 && diff < height {
					if o == 0 {
						p.bSpriteZeroHitPossible = true
					}
					p.spriteScanLine[p.spriteCount] = p.GetOAMEntry(o)
					p.spriteCount++
				}
			}
			p.status.spriteOverflow = p.spriteCount > 8
		}

		// ----------------------------------------------------------------
		// Load sprite pattern shifters (cycle 340)
		// ----------------------------------------------------------------
		if p.cycle == 340 {
			for i := 0; i < int(p.spriteCount); i++ {
				e := p.spriteScanLine[i]
				row := int16(p.scanline) - int16(e.y) - 1

				if e.attribute&0x80 != 0 { // vertical flip
					row ^= 0x07
					if p.control.spriteSize && row >= 8 {
						row ^= 0x07
					}
				}

				var addr uint16
				if !p.control.spriteSize { // 8×8
					addr = Btoi16(p.control.patternSprite)<<12 |
						uint16(e.id)<<4 | uint16(row&0x07)
				} else { // 8×16
					tile := e.id & 0xFE
					if row >= 8 {
						tile++
					}
					table := e.id & 0x01
					addr = uint16(table)<<12 | uint16(tile)<<4 | uint16(row&0x07)
				}

				p.spriteShifterPatternLo[i] = p.PpuRead(addr, false)
				p.spriteShifterPatternHi[i] = p.PpuRead(addr+8, false)

				if e.attribute&0x40 != 0 { // horizontal flip
					p.spriteShifterPatternLo[i] = flipByte(p.spriteShifterPatternLo[i])
					p.spriteShifterPatternHi[i] = flipByte(p.spriteShifterPatternHi[i])
				}
			}
		}

		//--------------------------------------------------------------------
		// Visible pixel compositing            (dots 1‑256 on visible lines)
		//--------------------------------------------------------------------
		if p.cycle > 0 && p.cycle <= 256 && p.scanline >= 0 {

			// ------------ background pixel -------------
			var bgPixel, bgPalette uint8
			if p.mask.renderBackground {
				bit := uint16(0x8000 >> p.fineX)
				if (p.bgShifterPatternLo & bit) != 0 { bgPixel |= 1 }
				if (p.bgShifterPatternHi & bit) != 0 { bgPixel |= 2 }
				if (p.bgShifterAttribLo  & bit) != 0 { bgPalette |= 1 }
				if (p.bgShifterAttribHi  & bit) != 0 { bgPalette |= 2 }
			}

			// ------------ foreground / sprite pixel -------------
			var fgPixel, fgPalette uint8
			fgPriorityFront := false // true -> sprite in front of BG
			if p.mask.renderSprites {
				for i := 0; i < int(p.spriteCount); i++ {
					if p.spriteScanLine[i].x == 0 {
						lo := (p.spriteShifterPatternLo[i] & 0x80) >> 7
						hi := (p.spriteShifterPatternHi[i] & 0x80) >> 6
						fgPixel = hi | lo
						if fgPixel != 0 {
							fgPalette = (p.spriteScanLine[i].attribute & 0x03) + 0x04
							fgPriorityFront = (p.spriteScanLine[i].attribute & 0x20) == 0

							if i == 0 && p.bSpriteZeroHitPossible &&
								bgPixel != 0 && fgPixel != 0 {
								if p.mask.renderBackground && p.mask.renderSprites {
									leftOK := !(p.mask.renderBackgroundLeft || p.mask.renderSpritesLeft)
									if (leftOK && p.cycle >= 9 && p.cycle < 258) ||
										(!leftOK && p.cycle >= 1 && p.cycle < 258) {
										p.status.spriteZeroHit = true
									}
								}
							}
							break
						}
					}
				}
			}

			// ------------ final pixel selection -------------
			var finalPixel, finalPalette uint8
			switch {
			case bgPixel == 0 && fgPixel == 0:
				finalPixel, finalPalette = 0, 0
			case bgPixel == 0 && fgPixel > 0:
				finalPixel, finalPalette = fgPixel, fgPalette
			case bgPixel > 0 && fgPixel == 0:
				finalPixel, finalPalette = bgPixel, bgPalette
			default: // both non‑zero
				if fgPriorityFront {
					finalPixel, finalPalette = fgPixel, fgPalette
				} else {
					finalPixel, finalPalette = bgPixel, bgPalette
				}
			}

			p.screen[p.cycle-1][p.scanline] =
				p.GetColourFromPaletteRam(finalPalette, finalPixel)
		}
	}

	

	//--------------------------------------------------------------------
	// VBlank (scanlines 241‑260)
	//--------------------------------------------------------------------
	if p.scanline == 241 && p.cycle == 1 {
		p.status.verticalBlank = true
		if p.control.enableNmi {
			p.Nmi = true
		}
	}

	//--------------------------------------------------------------------
	// advance clocks
	//--------------------------------------------------------------------
	p.cycle++
	if p.cycle >= 341 {
		p.cycle = 0
		p.scanline++
		if p.scanline >= 261 {
			p.scanline = -1
			p.FrameComplete = true
		}
	}
}

// flipByte reverses the bit order of b (helper for H‑flip)
func flipByte(b uint8) uint8 {
	b = (b&0xF0)>>4 | (b&0x0F)<<4
	b = (b&0xCC)>>2 | (b&0x33)<<2
	b = (b&0xAA)>>1 | (b&0x55)<<1
	return b
}
