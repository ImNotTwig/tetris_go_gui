package main

type Point struct {
	Row, Col int // y, x
}

// a Pixel is an integer signifying what color the pixel on the board should be
type Pixel int

type Board struct {
	Field [][]Pixel
}

func NewBoard() Board {
	var b Board
	for i := 0; i < 20; i++ {
		b.Field = append(b.Field, make([]Pixel, 0))
		for j := 0; j < 10; j++ {
			b.Field[i] = append(b.Field[i], Pixel(0))
		}
	}

	return b
}
