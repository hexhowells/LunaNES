package main

import (
	"LunaNES/emu"
	"LunaNES/pixelengine"
	"log"
	"time"

	"github.com/karalabe/usb"
)

func main() {
	bus := emu.NewBus()
	cart := emu.NewCartridge("../ROMS/IceClimber.nes")
	if cart == nil {
		log.Fatalln("Error: cartridge could not be loaded")
	}

	bus.InsertCartridge(cart)
	bus.Reset()

	go func() {
		devices, err := usb.Enumerate(0x081f, 0xe401)
		if err != nil {
			log.Fatalf("Enumeration error: %v", err)
		}
		if len(devices) == 0 {
			log.Fatal("NES controller not found")
		}

		device, err := devices[0].Open()
		if err != nil {
			log.Fatalf("Open error: %v", err)
		}
		defer device.Close()

		buffer := make([]byte, 8)

		for {
			count, err := device.Read(buffer)
			if err != nil {
				log.Printf("Read error: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if count > 0 {
				var controllerState byte = 0x00

				if buffer[5] == 47 || buffer[5] == 63 { // A or A+B
					controllerState |= 0x80
				}
				if buffer[5] == 31 || buffer[5] == 63 { // B or A+B
					controllerState |= 0x40
				}
				if buffer[6] == 16 || buffer[6] == 48 { // Select or Select+Start
					controllerState |= 0x20
				}
				if buffer[6] == 32 || buffer[6] == 48 { // Start or Select+Start
					controllerState |= 0x10
				}
				if buffer[1] == 0 { // Up
					controllerState |= 0x08
				}
				if buffer[1] == 255 { // Down
					controllerState |= 0x04
				}
				if buffer[0] == 0 { // Left
					controllerState |= 0x02
				}
				if buffer[0] == 255 { // Right
					controllerState |= 0x01
				}

				bus.Controller[0] = controllerState
			}
		}
	}()

	// Start emulation loop
	go func() {
		ticker := time.NewTicker(time.Second / 60)
		for range ticker.C {
			// Run until the next frame is complete
			for {
				bus.Clock()
				if bus.Ppu.FrameComplete {
					bus.Ppu.FrameComplete = false
					break
				}
			}

			// Render the screen from the PPU's framebuffer
			screen := bus.Ppu.Screen()
			for y := 0; y < 240; y++ {
				for x := 0; x < 256; x++ {
					p := screen[x][y]
					pixelengine.SetPixel(x, y, p.R, p.G, p.B)
				}
			}
		}
	}()

	pixelengine.Start()
}
