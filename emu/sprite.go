package emu


type Sprite struct {
	Pixels [][]Pixel
	Rows uint8
	Cols uint8
}


func CreateSprite(rows uint8, cols uint8) *Sprite {
	sprite := Sprite{}

	sprite.Pixels = make([][]Pixel, rows)
	for i := range sprite.Pixels {
		sprite.Pixels[i] = make([]Pixel, cols)
	}

	sprite.Rows = rows
	sprite.Cols = cols

	return &sprite
}


func (s *Sprite) SetPixel(row uint8, col uint8, value Pixel) {
	if row < s.Rows && col < s.Cols {
		s.Pixels[row][col] = value
	}
}
