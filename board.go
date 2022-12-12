package main

import (
	"fmt"
	"image/color"
)

type Point struct {
	Row, Col int // y, x
}

// a Pixel is an integer signifying what color the pixel on the board should be
type Pixel int

// a Tetro or otherwise known as a tetromino are the pieces that are dropped from the top of the board,
// the int value determines what type of tetro it is
type Tetro int

// types of Tetrominos:
// O: 1
// L: 2
// J: 3
// I: 4
// T: 5
// S: 6
// Z: 7

// a shape is a list of points, usually 4, that makes up a tetromino
type Shape [4]Point

type Tetromino struct {
	Shape Shape
	Tetro Tetro
}

// the playing board which is a 2-D array that has all the data for the pixels on it
type Board [][]Pixel

// returns a new board with an initialized 2-D array
func NewBoard() Board {
	var b Board
	for i := 0; i < 24; i++ {
		b = append(b, make([]Pixel, 0))
		for j := 0; j < 10; j++ {
			b[i] = append(b[i], Pixel(0))
		}
	}

	return b
}

// convert the int tetro to a shape
func (t Tetro) TetroToNewShape() Shape {
	var retShape Shape
	switch t {
	case 2: // L
		retShape = Shape{
			Point{Row: 1 + 20, Col: 0 + 4},
			Point{Row: 1 + 20, Col: 1 + 4},
			Point{Row: 1 + 20, Col: 2 + 4},
			Point{Row: 0 + 20, Col: 0 + 4},
		}
	case 4: // I
		retShape = Shape{
			Point{Row: 1 + 20, Col: 0 + 4},
			Point{Row: 1 + 20, Col: 1 + 4},
			Point{Row: 1 + 20, Col: 2 + 4},
			Point{Row: 1 + 20, Col: 3 + 4},
		}
	case 1: // O
		retShape = Shape{
			Point{Row: 1 + 20, Col: 0 + 4},
			Point{Row: 1 + 20, Col: 1 + 4},
			Point{Row: 0 + 20, Col: 0 + 4},
			Point{Row: 0 + 20, Col: 1 + 4},
		}
	case 5: // T
		retShape = Shape{
			Point{Row: 1 + 20, Col: 0 + 4},
			Point{Row: 1 + 20, Col: 1 + 4},
			Point{Row: 1 + 20, Col: 2 + 4},
			Point{Row: 0 + 20, Col: 1 + 4},
		}
	case 6: // S
		retShape = Shape{
			Point{Row: 0 + 20, Col: 0 + 4},
			Point{Row: 0 + 20, Col: 1 + 4},
			Point{Row: 1 + 20, Col: 1 + 4},
			Point{Row: 1 + 20, Col: 2 + 4},
		}
	case 7: // Z
		retShape = Shape{
			Point{Row: 1 + 20, Col: 0 + 4},
			Point{Row: 1 + 20, Col: 1 + 4},
			Point{Row: 0 + 20, Col: 1 + 4},
			Point{Row: 0 + 20, Col: 2 + 4},
		}
	case 3: // J
		retShape = Shape{
			Point{Row: 1 + 20, Col: 0 + 4},
			Point{Row: 0 + 20, Col: 1 + 4},
			Point{Row: 0 + 20, Col: 0 + 4},
			Point{Row: 0 + 20, Col: 2 + 4},
		}
	default:
		panic(fmt.Sprintf("Invalid integer passed into TetroToShape: %v", t))
	}
	return retShape
}

// convert the int tetro to a color
func (t Tetro) TetroToColor() color.RGBA {
	switch t {
	// nothing
	case 0:
		return color.RGBA{30, 30, 46, 255}

		// O
	case 1:
		return color.RGBA{249, 226, 175, 255}

		// L
	case 2:
		return color.RGBA{250, 179, 135, 255}

		// J
	case 3:
		return color.RGBA{137, 180, 250, 255}

		// I
	case 4:
		return color.RGBA{148, 226, 213, 255}

		// T
	case 5:
		return color.RGBA{203, 166, 247, 255}

		// S
	case 6:
		return color.RGBA{166, 227, 161, 255}

		// Z
	case 7:
		return color.RGBA{243, 139, 168, 255}
	}

	panic(fmt.Sprintf("Invalid integer passed into TetroToColor: %v", t))
}
