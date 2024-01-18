package main

import (
	"LunaNES/emu"
	"encoding/hex"
	"strings"
	"fmt"
	"sort"
	"log"
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
	cart := emu.NewCartridge("../ROMS/SuperMarioBros.nes")

	if cart == nil {
        log.Println("Error: cart is nil")
    }

	bus.InsertCartridge(cart)

	// Store a program into memory
	//hexString := "A2 0A 8E 00 00 A2 03 8E 01 00 AC 00 00 A9 00 18 6D 01 00 88 D0 FA 8D 02 00 EA EA EA"
	hexString := "A9 10 C9 A0 F0 07 90 0A B0 0D 4C 19 00 A9 01 4C 19 00 A9 FF 4C 19 00 A9 2A 00"
	hexString = strings.ReplaceAll(hexString, " ", "")
	bytes, _ := hex.DecodeString(hexString)

	bus.WriteBytes(0x0000, bytes)

	// Set the reset vector
	bus.CpuWrite(0xFFFC, 0x00)
	bus.CpuWrite(0xFFFD, 0x80)

	cpu.ConnectBus(bus)

	// Print initial state of the CPU
	cpu.PrintCPU()
	cpu.PrintStatusFlags()
	cpu.PrintRAM(0x80, 1)

	cpu.Reset()

	// manually set program counter after reset to force it to run code at 0x0000
	// this isnt correct for emulating the nes since the cartridge address space
	// starts at 0x8000 but works for simple programs that dont use the allocated space
	cpu.Pc = 0x0000

	new_inst := false

	dissasMap := cpu.Disassemble(0x0000, 0x0018)

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

	// Run the program by clocking the cpu
	for i := 0; i < 41; {
		new_inst = cpu.Clock()
		if new_inst {
			i++
			printCodeWindow(keys, cpu.Pc, dissasMap)
			cpu.PrintCPU()
			cpu.PrintStatusFlags()
			cpu.PrintRAM(0x00, 1)
		}
	}
}