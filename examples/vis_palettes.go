package main

import (
	"LunaNES/emu"
	"fmt"
	"log"
	"time"
    "LunaNES/pixelengine"
    "sort"
)


func printCodeWindow(keys []uint16, mainKey uint16, dissasMap map[uint16]string) {
	var mainKeyIndex int
    for i, key := range keys {
        if key == mainKey {
            mainKeyIndex = i
            break
        }
    }

    startIndex := mainKeyIndex - 5
    if startIndex < 0 {
        startIndex = 0
    }

    endIndex := mainKeyIndex + 10
    if endIndex >= len(dissasMap) {
        endIndex = len(dissasMap) - 1
    }

    // Print 5 elements before the mainKey, the mainKey's element, and 10 elements after the mainKey
    for i := startIndex; i <= endIndex; i++ {
        if i == mainKeyIndex {
            fmt.Printf("\n> %s", dissasMap[keys[i]])
        } else {
            fmt.Printf("\n  %s", dissasMap[keys[i]])
        }
    }
    fmt.Println()
}


func main() {
	bus := emu.NewBus()
	cart := emu.NewCartridge("../ROMS/nestest.nes")

	if cart == nil {
        log.Println("Error: cart is nil")
    }

	bus.InsertCartridge(cart)

	bus.Reset()

	dissasMap := bus.Cpu.Disassemble(0x8000, 0xCFFF)

	// Get all the keys of the map and sort them
	var keys []uint16
	for key := range dissasMap {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	
	// Print out each instruction
	fmt.Println("\nDisassembled Code\n---")
	for _, key := range keys {
		fmt.Println(dissasMap[key])
	}

	prev_pc := uint16(0)

	go func() {
    ticker := time.NewTicker(time.Second / 60)  // 60Hz
    for {
        select {
        case <-ticker.C:
            // Clock until a full PPU frame is complete
            for {
                bus.Clock()

                if bus.Ppu.FrameComplete {
                    bus.Ppu.FrameComplete = false
                    break // Finished a frame
                }
            }

            if prev_pc != bus.Cpu.Pc {
                prev_pc = bus.Cpu.Pc
                nSwatchSize := 6

                //printCodeWindow(keys, bus.Cpu.Pc, dissasMap)
                //bus.Cpu.PrintStatusFlags()
                //bus.Cpu.PrintRAM(0x00, 1)
                //bus.Cpu.PrintCPU()

                for p := 0; p < 8; p++ {
                    for s := 0; s < 4; s++ {
                        pix := bus.Ppu.GetColourFromPaletteRam(uint8(p), uint8(s))
                        pixelengine.SetRect(p*(nSwatchSize*5)+s*nSwatchSize, 0, 6, 6, pix.R, pix.G, pix.B)
                    }
                }
            }
        }
    }
}()

    pixelengine.Start()
}