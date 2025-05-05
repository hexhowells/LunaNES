package main

import (
    "fmt"
    "log"
    "time"

    "github.com/karalabe/usb"
)

func main() {
    // Enumerate device with NES controller's VID/PID
    devices, err := usb.Enumerate(0x081f, 0xe401)
    if err != nil {
        log.Fatalf("Enumeration error: %v", err)
    }

    if len(devices) == 0 {
        log.Fatal("NES controller not found")
    }

    // Open first found device
    device, err := devices[0].Open()
    if err != nil {
        log.Fatalf("Open error: %v", err)
    }
    defer device.Close()

    // Read 8 bytes buffer for controller input
    buffer := make([]byte, 8)

    for {
        // Read controller state
        count, err := device.Read(buffer)
        if err != nil {
            log.Printf("Read error: %v", err)
            time.Sleep(1 * time.Second)
            continue
        }

        if count > 0 {
            printButtons(buffer)
        }
    }
}

func printButtons(data []byte) {
    var pressed []string

    // D-pad
    if data[0] == 0 {
        pressed = append(pressed, "Left")
    } else if data[0] == 255 {
        pressed = append(pressed, "Right")
    }
    if data[1] == 0 {
        pressed = append(pressed, "Up")
    } else if data[1] == 255 {
        pressed = append(pressed, "Down")
    }

    // A
    if data[5] == 47 {
        pressed = append(pressed, "A")
    }
    // B
    if data[5] == 31 {
        pressed = append(pressed, "B")
    }
    // A + B
    if data[5] == 63 {
        pressed = append(pressed, "A")
        pressed = append(pressed, "B")
    }

    // Select
    if data[6] == 16 {
        pressed = append(pressed, "Select")
    }
    // Start
    if data[6] == 32 {
        pressed = append(pressed, "Start")
    }
    // Select + Start
    if data[6] == 48 {
        pressed = append(pressed, "Select")
        pressed = append(pressed, "Start")
    }

    fmt.Println(data)
    if len(pressed) > 0 {
        fmt.Printf("Buttons pressed: %v\n", pressed)
    } else {
        fmt.Println("No buttons pressed")
    }
}
