package emu


type Sprite struct {
	pixels [][]Pixel
	rows uint8
	cols uint8
}


func CreateSprite(rows uint8, cols uint8) *Sprite {
	sprite = Sprite{}
	sprite.pixels = [rows][cols]Pixel{}
	sprite.rows = rows
	sprite.cols = cols

	return &sprite
}


func (s *Sprite) SetPixel(row uint8, col uint8, value uint8) {
	if 0 <= row <= s.rows && 0 <= col <= s.cols {
		s.pixels[row][col] = value
	}
}
