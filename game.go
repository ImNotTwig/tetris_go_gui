package main

import (
	"math/rand"
	"time"
)

func ContainsShape(sl Shape, point *Point) bool {
	for _, v := range sl {
		if v == *point {
			return true
		}
	}
	return false
}
func ContainsSlice(sl []*Point, point *Point) bool {
	for _, v := range sl {
		if v == point {
			return true
		}
	}
	return false
}

type Game struct {
	PlayingBoard Board

	CurrentPiece *Tetromino

	Current7Bag []*Tetromino

	Score int

	Level int

	GameOver bool

	Paused bool

	FallingSpeedMillis int
}

func NewGame() Game {
	return Game{
		PlayingBoard:       NewBoard(),
		CurrentPiece:       nil,
		Current7Bag:        nil,
		Score:              0,
		Level:              0,
		GameOver:           false,
		Paused:             false,
		FallingSpeedMillis: 200,
	}
}

func (g *Game) GenerateNewBag() {
	if g.Current7Bag == nil || len(g.Current7Bag) == 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		tetro_list := []*Tetromino{
			{
				Tetro: Tetro(1),
				Shape: Tetro(1).TetroToNewShape(),
			},
			{
				Tetro: Tetro(2),
				Shape: Tetro(2).TetroToNewShape(),
			},
			{
				Tetro: Tetro(3),
				Shape: Tetro(3).TetroToNewShape(),
			},
			{
				Tetro: Tetro(4),
				Shape: Tetro(4).TetroToNewShape(),
			},
			{
				Tetro: Tetro(5),
				Shape: Tetro(5).TetroToNewShape(),
			},
			{
				Tetro: Tetro(6),
				Shape: Tetro(6).TetroToNewShape(),
			},
			{
				Tetro: Tetro(7),
				Shape: Tetro(7).TetroToNewShape(),
			},
		}

		r.Shuffle(len(tetro_list), func(i, j int) {
			tetro_list[i], tetro_list[j] = tetro_list[j], tetro_list[i]
		})

		g.Current7Bag = tetro_list
	}
}

func (g *Game) SetNextTetroFromBag() {
	if len(g.Current7Bag) > 0 {
		g.CurrentPiece = g.Current7Bag[0]
		for i := 0; i < len(g.CurrentPiece.Shape); i++ {
			g.PlayingBoard[g.CurrentPiece.Shape[i].Row][g.CurrentPiece.Shape[i].Col] = Pixel(g.CurrentPiece.Tetro)
		}
		g.Current7Bag = g.Current7Bag[1:]
	} else {
		g.GenerateNewBag()
	}
}

func (g *Game) GetRandomTetromino() {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random_number := r.Intn(7) + 1

	shape := Tetro(random_number).TetroToNewShape()

	g.CurrentPiece = &Tetromino{
		Shape: shape,
		Tetro: Tetro(random_number),
	}

	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		g.PlayingBoard[g.CurrentPiece.Shape[i].Row][g.CurrentPiece.Shape[i].Col] = Pixel(g.CurrentPiece.Tetro)
	}
}

func (g *Game) CheckIfSomethingUnder(s Shape) bool {
	for i := 0; i < len(s); i++ {
		x := s[i].Col
		y := s[i].Row
		if y != 0 {
			if g.PlayingBoard[y-1][x] != Pixel(0) && !ContainsShape(s, &Point{Row: y - 1, Col: x}) {
				return true
			}
		}
	}
	return false
}
func (g *Game) CheckIfSomethingRight(s Shape) bool {
	for i := 0; i < len(s); i++ {
		x := s[i].Col
		y := s[i].Row
		if x+1 < 10 {
			if g.PlayingBoard[y][x+1] != Pixel(0) && !ContainsShape(s, &Point{Row: y, Col: x + 1}) {
				return true
			}
		}
	}
	return false
}
func (g *Game) CheckIfSomethingLeft(s Shape) bool {
	for i := 0; i < len(s); i++ {
		x := s[i].Col
		y := s[i].Row
		if x-1 > 0 {
			if g.PlayingBoard[y][x-1] != Pixel(0) && !ContainsShape(s, &Point{Row: y, Col: x - 1}) {
				return true
			}
		}
	}
	return false
}

func (g *Game) GravityDrop() bool {
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		if g.CurrentPiece.Shape[j].Row-1 < 0 {
			return false
		}
	}

	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		if g.CheckIfSomethingUnder(g.CurrentPiece.Shape) {
			return false
		}
	}

	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[g.CurrentPiece.Shape[j].Row][g.CurrentPiece.Shape[j].Col] = Pixel(0)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[g.CurrentPiece.Shape[j].Row-1][g.CurrentPiece.Shape[j].Col] = Pixel(g.CurrentPiece.Tetro)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.CurrentPiece.Shape[j].Row -= 1
	}
	return true
}

func (g *Game) MoveRight() bool {
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		if g.CurrentPiece.Shape[j].Col+1 >= 10 {
			return false
		}
	}

	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		if g.CheckIfSomethingRight(g.CurrentPiece.Shape) {
			return false
		}
	}

	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[g.CurrentPiece.Shape[j].Row][g.CurrentPiece.Shape[j].Col] = Pixel(0)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[g.CurrentPiece.Shape[j].Row][g.CurrentPiece.Shape[j].Col+1] = Pixel(g.CurrentPiece.Tetro)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.CurrentPiece.Shape[j].Col += 1
	}
	return true
}

func (g *Game) MoveLeft() bool {
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		if g.CurrentPiece.Shape[j].Col-1 < 0 {
			return false
		}
	}

	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		if g.CheckIfSomethingLeft(g.CurrentPiece.Shape) {
			return false
		}
	}

	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[g.CurrentPiece.Shape[j].Row][g.CurrentPiece.Shape[j].Col] = Pixel(0)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[g.CurrentPiece.Shape[j].Row][g.CurrentPiece.Shape[j].Col-1] = Pixel(g.CurrentPiece.Tetro)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.CurrentPiece.Shape[j].Col -= 1
	}
	return true
}

func (g *Game) RotateClockWise() {
	var retShape Shape
	pivot := g.CurrentPiece.Shape[1]
	retShape[1] = pivot
	for i := 0; i < 4; i++ {
		// Index 1 is the pivot point
		if i == 1 {
			continue
		}
		dRow := pivot.Row - g.CurrentPiece.Shape[i].Row
		dCol := pivot.Col - g.CurrentPiece.Shape[i].Col
		retShape[i].Row = pivot.Row + (dCol * -1)
		retShape[i].Col = pivot.Col + (dRow)
	}

	for i := 0; i < len(retShape); i++ {
		if retShape[i].Row < 0 {
			sub := retShape[i].Row
			for sub < 0 {
				for j := 0; j < len(retShape); j++ {
					retShape[j].Row += 1
				}
				sub++
			}
			break
		} else if retShape[i].Col < 0 {
			sub := retShape[i].Col
			for sub < 0 {
				for j := 0; j < len(retShape); j++ {
					retShape[j].Col += 1
				}
				sub++
			}
			break
		}
		if retShape[i].Col >= 10 {
			sub := retShape[i].Col
			for sub >= 10 {
				for j := 0; j < len(retShape); j++ {
					retShape[j].Col -= 1
				}
				sub--
			}
			break
		}
		if g.PlayingBoard[retShape[i].Row][retShape[i].Col] != Pixel(0) &&
			!ContainsShape(g.CurrentPiece.Shape, &Point{
				Row: retShape[i].Row,
				Col: retShape[i].Col,
			}) {
			return
		}
	}

	for i := 0; i < 4; i++ {
		g.PlayingBoard[g.CurrentPiece.Shape[i].Row][g.CurrentPiece.Shape[i].Col] = Pixel(0)
	}

	g.CurrentPiece.Shape = retShape

	for i := 0; i < 4; i++ {
		g.PlayingBoard[g.CurrentPiece.Shape[i].Row][g.CurrentPiece.Shape[i].Col] = Pixel(g.CurrentPiece.Tetro)
	}

}
