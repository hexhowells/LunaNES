import "errors"


type Bus struct {
	ram [0xFFFF + 1]byte
	cpu CPU
}


func NewBus() *Bus {
	bus := Bus{}
	bus.cpu.ConnectBus(bus)

	for i := range bus.ram {
		bus.ram[i] = 0x00
	}

	return &bus
}


func (b *Bus) Write(addr uint16, data byte) error {
	if add >= 0x0000 && addr <= 0xFFFF {
		b.ram[addr] = data
		return nil
	}
	
	return errors.New("address out of bounds")
}


func (b *Bus) Read(addr uint16, bReadOnly boolean) (byte, error) {
	if add >= 0x0000 && addr <= 0xFFFF {
		return b.ram[addr], nil
	}

	return 0x00, errors.New("address out of bounds")
}