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
	cpu := emu.NewCPU()
	bus := emu.NewBus()
	cart := emu.NewCartridge("../ROMS/palette_fill_novblank.nes")
	new_inst := false

	if cart == nil {
        log.Println("Error: cart is nil")
    }

	bus.InsertCartridge(cart)

	// Set the reset vector
	bus.CpuWrite(0xFFFC, 0x00)
	bus.CpuWrite(0xFFFD, 0x80)

	cpu.ConnectBus(bus)

	cpu.Reset()

	cpu.Pc = 0x8100

	dissasMap := cpu.Disassemble(0x8000, 0x81FF)

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

	go func() {
        ticker := time.NewTicker(time.Second / 1000)
        for {
            select {
            case <-ticker.C:
                // Run the program by clocking the cpu
				new_inst = cpu.Clock()
				if new_inst {
					nSwatchSize := 6

					for p := 0; p < 8; p++ {
						for s := 0; s < 4; s++ {
							//printCodeWindow(keys, cpu.Pc, dissasMap)
							//cpu.PrintStatusFlags()
							//cpu.PrintRAM(0x00, 1)
							//cpu.PrintCPU()
							pix := bus.Ppu.GetColourFromPaletteRam(uint8(p), uint8(s))
							pixelengine.SetRect(p * (nSwatchSize * 5) + s * nSwatchSize, 0, 6, 6, pix.R, pix.G, pix.B)
						}
					}
				}
				
            }
        }
    }()

    pixelengine.Start()
}