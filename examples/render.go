package main

import (
	"LunaNES/emu"
	"LunaNES/pixelengine"
	"log"
	"time"
)

func main() {
	bus := emu.NewBus()
	cart := emu.NewCartridge("../ROMS/nestest.nes")
	if cart == nil {
		log.Fatalln("Error: cartridge could not be loaded")
	}

	bus.InsertCartridge(cart)
	bus.Reset()

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
