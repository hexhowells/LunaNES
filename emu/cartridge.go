package emu

type Cartridge struct {
	prgMemory []uint8  // program memory
	chrMemory []uint8  // character memory
	mapperID uint8
	numPrgBanks uint8  // how many banks of memory for the program
	numChrBanks uint8  // how many banks of memory for the characters
	header sHeader  // INES file header
	mapper Mapper  // onboard mapper
	imageValid bool  // 
}


type sHeader struct {
	name [4]byte
	prgRomChunks uint8
	chrRomChunks uint8
	mapper1 uint8
	mapper2 uint8
	prgRamSize uint8
	tvSystem1 uint8
	tvSystem2 uint8
	unused [5]byte
}


func NewCartridge(filename string) *Cartridge {
	cart :=- Cartridge{}
	cart.mapperID = 0
	cart.numPrgBanks = 0
	cart.numChrBanks = 0

	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Read file header
	err = binary.Read(file, binary.LittleEndian, cart.header)
	if err != nil {
		return nil
	}

	// Skip training info if present
	if cart.header.mapper1 & 0x04 != 0 {
		file.Seek(512, 1)
	}

	// Determine the mapper ID
	cart.mapperID = ((cart.header.mapper2 >> 4) << 4) | (cart.header.mapper1 >> 4)

	if cart.header.mapper1 & 0x01 != 0 {
		cart.mirror = 1  // vertical
	} else {
		cart.mirror = 0  // horizontal
	}

	nFileType := 1  // I NES file type (there are 3 types of file)

	if nFileType == 0 {

	}

	if nFileType == 1 {
		// Load program memory
		cart.numPrgBanks = cart.header.prgRomChunks
		cart.prgMemory = make([]uint8, cart.numPrgBanks * 16384)
		_, err = file.Read(cart.prgMemory)
		if err != nil {
			return nil
		}

		// Load character memory
		cart.numChrBanks = cart.header.chrRomChunks
		cart.chrMemory = make([]uint8, cart.numChrBanks * 8192)
		_, err = file.Read(cart.chrMemory) //
		if err != nil {
			return nil
		}

		// Load the mapper
		switch cart.mapperID {
			case 0:
				cart.mapper = NewMapper(cart.numPrgBanks, cart.numChrBanks)
		}

		cart.imageValid = true
	}

	if nFileType == 2 {

	}

	return &cart
}


func (cart *Cartridge) ImageValid() {
	return cart.imageValid
}


func (cart *Cartridge) CpuRead(addr uint16, &data uint8) bool {
	mapperAddr := uint32(0)

	if cart.mapper.CpuMapRead(addr, mappedAddr) {
		data := cart.prgMemory[mappedAddr]
		return true
	} else {
		return false
	}
}


func (cart *Cartridge) CpuWrite(addr uint16, data uint8) bool {
	mapperAddr := uint32(0)

	if cart.mapper.CpuMapWrite(addr, mappedAddr) {
		cart.prgMemory[mappedAddr] = data
		return true
	} else {
		return false
	} 
}


func (cart *Cartridge) PpuRead(addr uint16, &data uint8) bool {
	mapperAddr := uint32(0)

	if cart.mapper.PpuMapRead(addr, mappedAddr) {
		data := cart.prgMemory[mappedAddr]
		return true
	} else {
		return false
	}
}


func (cart *Cartridge) PpuWrite(addr uint16, data uint8) bool {
	mapperAddr := uint32(0)

	if cart.mapper.PpuMapWrite(addr, mappedAddr) {
		cart.prgMemory[mappedAddr] = data
		return true
	} else {
		return false
	} 
}
