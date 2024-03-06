package emu


type Sprite struct {
	Pixels [][]Pixel
	Rows uint8
	Cols uint8
}


func CreateSprite(rows uint8, cols uint8) *Sprite {
	sprite = Sprite{}
	sprite.Pixels = [rows][cols]Pixel{}
	sprite.Rows = rows
	sprite.Cols = cols

	return &sprite
}


func (s *Sprite) SetPixel(row uint8, col uint8, value uint8) {
	if 0 <= row <= s.Rows && 0 <= col <= s.Cols {
		s.Pixels[row][col] = value
	}
}
