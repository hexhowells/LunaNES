package emu

import (
	"os"
	"encoding/binary"
	"log"
)


const (
	HORIZONTAL = iota
	VERTICAL
	ONESCREEN_LO
	ONESCREEN_HI
)


type Cartridge struct {
	prgMemory []uint8  // program memory
	chrMemory []uint8  // character memory
	mapperID uint8
	numPrgBanks uint8  // how many banks of memory for the program
	numChrBanks uint8  // how many banks of memory for the characters
	header sHeader  // INES file header
	mapper *Mapper000  // onboard mapper
	imageValid bool
	mirror int
}


type sHeader struct {
	Name [4]byte  // Header name
	PrgRomChunks uint8  // number of program rom pages
	ChrRomChunks uint8  // number of character rom pages
	Mapper1 uint8
	Mapper2 uint8
	PrgRamSize uint8  // rarely used
	TvSystem1 uint8  // rarely used
	TvSystem2 uint8  // rarely used
	Unused [5]byte  // padding
}


func NewCartridge(filename string) *Cartridge {
	cart := Cartridge{}
	cart.mapperID = 0
	cart.numPrgBanks = 0
	cart.numChrBanks = 0

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error: could not open ROM file")
		log.Println(err)
		return nil
	}
	defer file.Close()

	// Read file header
	err = binary.Read(file, binary.LittleEndian, &cart.header)
	if err != nil {
		log.Println("Error: could not read ROM file header")
		log.Println(err)
		return nil
	}

	// Skip training info if present
	if cart.header.Mapper1 & 0x04 != 0 {
		file.Seek(512, 1)
	}

	// Determine the mapper ID
	cart.mapperID = ((cart.header.Mapper2 >> 4) << 4) | (cart.header.Mapper1 >> 4)

	if cart.header.Mapper1 & 0x01 != 0 {
		cart.mirror = 1  // vertical
	} else {
		cart.mirror = 0  // horizontal
	}

	nFileType := 1  // I NES file type (there are 3 types of file)

	if nFileType == 0 {

	}

	if nFileType == 1 {
		// Load program memory
		cart.numPrgBanks = cart.header.PrgRomChunks
		cart.prgMemory = make([]uint8, uint16(cart.numPrgBanks) * 16384)
		_, err = file.Read(cart.prgMemory)
		if err != nil {
			log.Println("Error: could not read program memory from ROM")
			log.Println(err)
			return nil
		}

		// Load character memory
		cart.numChrBanks = cart.header.ChrRomChunks
		cart.chrMemory = make([]uint8, uint16(cart.numChrBanks) * 8192)
		_, err = file.Read(cart.chrMemory) //
		if err != nil {
			log.Println("Error: could not read character memory from ROM")
			log.Println(err)
			return nil
		}

		// Load the mapper
		switch cart.mapperID {
			case 0:
				cart.mapper = NewMapper_000(cart.numPrgBanks, cart.numChrBanks)
			default:
				log.Println("Mapper not supported for this cartridge")
				return nil
		}

		cart.imageValid = true
	}

	if nFileType == 2 {

	}

	return &cart
}


func (cart *Cartridge) ImageValid() bool {
	return cart.imageValid
}


func (cart *Cartridge) CpuRead(addr uint16, data *uint8) bool {
	mappedAddr := uint32(0)

	if cart.mapper.CpuMapRead(addr, &mappedAddr) {
		*data = cart.prgMemory[mappedAddr]
		return true
	} else {
		return false
	}
}


func (cart *Cartridge) CpuWrite(addr uint16, data uint8) bool {
	mappedAddr := uint32(0)

	if cart.mapper.CpuMapWrite(addr, &mappedAddr) {
		cart.prgMemory[mappedAddr] = data
		return true
	} else {
		return false
	} 
}


func (cart *Cartridge) PpuRead(addr uint16, data *uint8) bool {
	mappedAddr := uint32(0)

	if cart.mapper.PpuMapRead(addr, &mappedAddr) {
		*data = cart.chrMemory[mappedAddr]
		return true
	} else {
		return false
	}
}


func (cart *Cartridge) PpuWrite(addr uint16, data uint8) bool {
	mappedAddr := uint32(0)

	if cart.mapper.PpuMapWrite(addr, &mappedAddr) {
		cart.chrMemory[mappedAddr] = data
		return true
	} else {
		return false
	} 
}
