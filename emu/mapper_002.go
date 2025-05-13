package emu

import ("log")

type Mapper002 struct {
	Mapper
	nPRGBankSelectLo uint8
	nPRGBankSelectHi uint8
}

func NewMapper_002(prgBanks uint8, chrBanks uint8) *Mapper002 {
	log.Printf("Initializing Mapper 002 with %d PRG banks, %d CHR banks\n", prgBanks, chrBanks)
	mapper := Mapper002{}
	mapper.numPrgBanks = prgBanks
	mapper.numChrBanks = chrBanks

	return &mapper
}

func (mapper *Mapper002) CpuMapRead(addr uint16, mapped_addr *uint32, data *uint8) bool {
	if addr >= 0x8000 && addr <= 0xBFFF {
		*mapped_addr = uint32(mapper.nPRGBankSelectLo)*0x4000 + uint32(addr&0x3FFF)
		log.Printf("[CpuMapRead] addr: %04X → Lo bank: %d, mapped_addr: %d\n", addr, mapper.nPRGBankSelectLo, *mapped_addr)
		return true
	}
	if addr >= 0xC000 && addr <= 0xFFFF {
		*mapped_addr = uint32(mapper.nPRGBankSelectHi)*0x4000 + uint32(addr&0x3FFF)
		log.Printf("[CpuMapRead] addr: %04X → Hi bank: %d, mapped_addr: %d\n", addr, mapper.nPRGBankSelectHi, *mapped_addr)
		return true
	}
	return false
}


func (mapper *Mapper002) CpuMapWrite(addr uint16, mapped_addr *uint32, data uint8) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		mapper.nPRGBankSelectLo = data & 0x0F
	}
	return false
}


func (mapper *Mapper002) PpuMapRead(addr uint16, mapped_addr *uint32) bool {
	if addr < 0x2000 {
		*mapped_addr = uint32(addr)
		return true
	}
	return false
}


func (mapper *Mapper002) PpuMapWrite(addr uint16, mapped_addr *uint32) bool {
	if addr < 0x2000 {
		if mapper.numChrBanks == 0 {
			*mapped_addr = uint32(addr)
			return true
		}
	}
	return false
}

func (mapper *Mapper002) Reset() {
	mapper.nPRGBankSelectLo = 0
	mapper.nPRGBankSelectHi = mapper.numPrgBanks - 1
}
