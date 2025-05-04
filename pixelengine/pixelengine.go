package pixelengine

import (
    "image/color"
    "time"
    "LunaNES/emu"
    "github.com/hajimehoshi/ebiten/v2"
)


type Pixel struct {
    R, G, B uint8
}

const (
    ScreenWidth  = 256
    ScreenHeight = 240
    WindowWidth = 720
    WindowHeight = 675
    fontSize = 12
)

var pixels [ScreenWidth][ScreenHeight]Pixel


func init() {

}


// Window struct and functions (required for ebiten)
type Window struct {
    lastUpdateTime time.Time
}


func (g *Window) Update() error {
    // Update logic can be added here if needed
    return nil
}


func Clear() {
    for x := 0; x < ScreenWidth; x++ {
        for y := 0; y < ScreenHeight; y++ {
            pixels[x][y] = Pixel{R: 0, G: 0, B: 0}
        }
    }
}


func (g *Window) Draw(screen *ebiten.Image) {
    for x := 0; x < ScreenWidth; x++ {
        for y := 0; y < ScreenHeight; y++ {
            p := pixels[x][y]
            screen.Set(x, y, color.RGBA{R: p.R, G: p.G, B: p.B, A: 255})
        }
    }
}


func (g *Window) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ScreenWidth, ScreenHeight
}


// Start the pixel engine
func Start() {
    window := &Window{lastUpdateTime: time.Now()}
    ebiten.SetWindowSize(WindowWidth, WindowHeight)
    ebiten.SetWindowTitle("LunaNES Emulator")

    if err := ebiten.RunGame(window); err != nil {
        panic(err)
    }
}


// Set a pixel to an given RGB value
func SetPixel(x, y int, r, g, b uint8) {
    if x >= 0 && x < WindowWidth && y >= 0 && y < WindowHeight {
        pixels[x][y] = Pixel{R: r, G: g, B: b}
    } else {
        panic("Pixel coordinate value out of range")
    }
}


// Set a rectangle starting at (x,y) with size (h,w) of pixel value (r,g,b)
func SetRect(x, y, h, w int, r, g, b uint8) {
    for dx := 0; dx < w; dx++ {
        for dy := 0; dy < h; dy++ {
            SetPixel(x+dx, y+dy, r, g, b)
        }
    }
}


// Set a sprite pixel map starting at (x,y)
func SetSprite(x, y int, s emu.Sprite) {
    for dx := 0; dx < int(s.Rows); dx++ {
        for dy := 0; dy < int(s.Cols); dy++ {
            pixel := s.Pixels[dx][dy]
            SetPixel(x+dx, y+dy, pixel.R, pixel.G, pixel.B)
        }
    }
}
