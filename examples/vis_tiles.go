package main

import (
	"LunaNES/emu"
	"log"
	"time"
	"LunaNES/pixelengine"
)

func GetTilePixels(ppu *emu.PPU, patternTableBase uint16, tileIndex uint8) [8][8]uint8 {
	var tile [8][8]uint8
	baseAddress := patternTableBase + uint16(tileIndex)*16
	for y := 0; y < 8; y++ {
		plane0 := ppu.PpuRead(baseAddress+uint16(y), true)
		plane1 := ppu.PpuRead(baseAddress+uint16(y)+8, true)
		for x := 0; x < 8; x++ {
			bit0 := (plane0 >> (7 - x)) & 1
			bit1 := (plane1 >> (7 - x)) & 1
			tile[y][x] = (bit1 << 1) | bit0
		}
	}
	return tile
}

func main() {
	bus := emu.NewBus()
	cart := emu.NewCartridge("../ROMS/SuperMarioBros.nes")

	if cart == nil {
		log.Fatalln("Error: cartridge could not be loaded")
	}

	bus.InsertCartridge(cart)
	bus.Reset()

	go func() {
		ticker := time.NewTicker(time.Second / 120)
		for range ticker.C {
			for {
				bus.Clock()
				if bus.Ppu.FrameComplete {
					bus.Ppu.FrameComplete = false
					break
				}
			}

			pixelengine.Clear()

			tileSize := 8
			tilesPerRow := 16

			for table := 0; table < 2; table++ {
				base := uint16(table * 0x1000)

				for i := 0; i < 256; i++ {
					tile := GetTilePixels(&bus.Ppu, base, uint8(i))
					tileX := (i % tilesPerRow) * tileSize
					tileY := (i / tilesPerRow) * tileSize

					// Draw each tile pixel
					for y := 0; y < tileSize; y++ {
						for x := 0; x < tileSize; x++ {
							colorID := tile[y][x]
							color := bus.Ppu.GetColourFromPaletteRam(0, colorID)
							// Pattern Table 0 starts at x = 0
							// Pattern Table 1 starts at x = 128 (16 tiles * 8 px)
							pixelengine.SetPixel(
								table*128+tileX+x, // x position: shift table 1 to the right
								tileY+y,
								color.R, color.G, color.B,
							)
						}
					}
				}
			}
		}
	}()

	pixelengine.Start()
}
