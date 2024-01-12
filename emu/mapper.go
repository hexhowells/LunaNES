package emu


type MapperInterface interface {
	CpuMapRead(addr uint16, mapped_addr *uint32) bool 
	CpuMapWrite(addr uint16, mapped_addr *uint32) bool 
	PpuMapRead(addr uint16, mapped_addr *uint32) bool 
	PpuMapWrite(addr uint16, mapped_addr *uint32) bool 
}


type Mapper struct {
	numPrgBanks uint8
	numChrBanks uint8
}
