package emu

import (
	"os"
	"encoding/binary"
	"io"
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
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error: could not open ROM file")
		log.Println(err)
		return nil
	}
	defer file.Close()

	// Read 16-byte iNES header
	err = binary.Read(file, binary.LittleEndian, &cart.header)
	if err != nil {
		log.Println("Error: could not read ROM file header")
		log.Println(err)
		return nil
	}

	// Seek to byte 16 (past header)
	_, err = file.Seek(16, 0)
	if err != nil {
		log.Println("Error: failed to seek past header")
		log.Println(err)
		return nil
	}

	// If trainer is present, skip 512 bytes
	if cart.header.Mapper1&0x04 != 0 {
		log.Println("Trainer detected, skipping 512 bytes")
		file.Seek(512, 1)
	}

	// Determine mapper ID
	cart.mapperID = ((cart.header.Mapper2 >> 4) << 4) | (cart.header.Mapper1 >> 4)

	// Determine mirroring mode
	if cart.header.Mapper1&0x01 != 0 {
		cart.mirror = VERTICAL
	} else {
		cart.mirror = HORIZONTAL
	}

	// Determine file type
	nFileType := 1
	if (cart.header.Mapper2 & 0x0C) == 0x08 {
		nFileType = 2
	}

	if nFileType == 1 {
		// Load PRG-ROM
		cart.numPrgBanks = cart.header.PrgRomChunks
		cart.prgMemory = make([]uint8, uint32(cart.numPrgBanks)*16384)
		_, err = io.ReadFull(file, cart.prgMemory)
		if err != nil {
			log.Println("Error: could not fully read PRG-ROM")
			log.Println(err)
			return nil
		}

		// Load CHR-ROM or allocate CHR-RAM
		cart.numChrBanks = cart.header.ChrRomChunks
		if cart.numChrBanks == 0 {
			cart.chrMemory = make([]uint8, 8192)
		} else {
			cart.chrMemory = make([]uint8, uint32(cart.numChrBanks)*8192)
			_, err = io.ReadFull(file, cart.chrMemory)
			if err != nil {
				log.Println("Error: could not fully read CHR-ROM")
				log.Println(err)
				return nil
			}
		}
	}

	// Mapper setup
	switch cart.mapperID {
	case 0:
		cart.mapper = NewMapper_000(cart.numPrgBanks, cart.numChrBanks)
	case 2:
		cart.mapper = NewMapper_002(cart.numPrgBanks, cart.numChrBanks)
	default:
		log.Println("Mapper not supported:", cart.mapperID)
		return nil
	}

	cart.imageValid = true
	return &cart
}

func (cart *Cartridge) ImageValid() bool {
	return cart.imageValid
}

func (cart *Cartridge) CpuRead(addr uint16, data *uint8) bool {
	mappedAddr := uint32(0)
	if cart.mapper.CpuMapRead(addr, &mappedAddr, data) {
		if int(mappedAddr) >= len(cart.prgMemory) {
			log.Printf("OUT OF BOUNDS READ: mappedAddr=%d, prgMemory size=%d", mappedAddr, len(cart.prgMemory))
			return false
		}
		*data = cart.prgMemory[mappedAddr]
		return true
	}
	return false
}

func (cart *Cartridge) CpuWrite(addr uint16, data uint8) bool {
	mappedAddr := uint32(0)
	if cart.mapper.CpuMapWrite(addr, &mappedAddr, data) {
		if int(mappedAddr) >= len(cart.prgMemory) {
			log.Printf("OUT OF BOUNDS WRITE: mappedAddr=%d, prgMemory size=%d", mappedAddr, len(cart.prgMemory))
			return false
		}
		cart.prgMemory[mappedAddr] = data
		return true
	}
	return false
}

func (cart *Cartridge) PpuRead(addr uint16, data *uint8) bool {
	mappedAddr := uint32(0)
	if cart.mapper.PpuMapRead(addr, &mappedAddr) {
		*data = cart.chrMemory[mappedAddr]
		return true
	}
	return false
}

func (cart *Cartridge) PpuWrite(addr uint16, data uint8) bool {
	mappedAddr := uint32(0)
	if cart.mapper.PpuMapWrite(addr, &mappedAddr) {
		cart.chrMemory[mappedAddr] = data
		return true
	}
	return false
}

func (cart *Cartridge) Reset() {
	cart.mapper.Reset()
}
