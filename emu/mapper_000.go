package emu


type Mapper000 struct {
	Mapper
}

func NewMapper_000(prgBanks uint8, chrBanks uint8) *Mapper000 {
	mapper := Mapper000{}
	mapper.numPrgBanks = prgBanks
	mapper.numChrBanks = chrBanks

	return &mapper
}

func (mapper *Mapper) CpuMapRead(addr uint16, mapped_addr *uint32) boolean {
	if addr >= 0x8000 && addr <= 0xFFFF {
		if mapper.numPrgBanks > 1 {
			mapped_addr = addr & 0x7FFF
		} else {
			mapped_addr = addr & 0x3FFF
		}
		return true
	}
	return false
}


func (mapper *Mapper) CpuMapWrite(addr uint16, mapped_addr *uint32) boolean {
	if addr >= 0x8000 && addr <= 0xFFFF {
		if mapper.numPrgBanks > 1 {
			mapped_addr = addr & 0x7FFF
		} else {
			mapped_addr = addr & 0x3FFF
		}
		return true
	}
	return false
}


func (mapper *Mapper) PpuMapRead(addr uint16, mapped_addr *uint32) boolean {
	if addr >= 0x0000 && addr <= 0x1FFF {
		mapped_addr = addr
		return true
	}
	return false
}


func (mapper *Mapper) PpuMapWrite(addr uint16, mapped_addr *uint32) boolean {
	if addr >= 0x0000 && addr <= 0x1FFF {
		if mappper.numChrBanks == 0 {
			mapped_addr = addr
			return true
		}
	}
	return false
}
