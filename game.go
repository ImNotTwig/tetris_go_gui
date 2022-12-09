package main

import (
	"math/rand"
	"time"
)

func remove(slice []*Point, s int) []*Point {
	return append(slice[:s], slice[s+1:]...)
}

func ContainsShape(sl Shape, point *Point) bool {
	for _, v := range sl {
		if &v == point {
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
		Score:              0,
		Level:              0,
		GameOver:           false,
		Paused:             false,
		FallingSpeedMillis: 500,
	}
}

func (g *Game) GetRandomTetromino() {
	if g.CurrentPiece == nil {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		random_number := r.Intn(7) + 1

		shape := Tetro(random_number).TetroToNewShape()

		g.CurrentPiece = &Tetromino{
			Color: Tetro(random_number).TetroToColor(),
			Shape: shape,
			Tetro: Tetro(random_number),
		}
		println(Tetro(random_number))
	}
	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		g.PlayingBoard[g.CurrentPiece.Shape[i].Row][g.CurrentPiece.Shape[i].Col] = Pixel(g.CurrentPiece.Tetro)
	}
}

func (g *Game) GravityDrop() bool {
	// check if it can move down
	// if it can remove the previous pixels on the board
	// the add the new ones
	var lowest_points []*Point
	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		if lowest_points == nil {
			lowest_points = []*Point{&g.CurrentPiece.Shape[i]}
		} else {
			for j := 0; j < len(lowest_points); j++ {
				if lowest_points[j].Row > g.CurrentPiece.Shape[i].Row {
					remove(lowest_points, j)
					lowest_points = append(lowest_points, &g.CurrentPiece.Shape[i])
					break
				} else if lowest_points[j].Row == g.CurrentPiece.Shape[i].Row {
					lowest_points = append(lowest_points, &g.CurrentPiece.Shape[i])
					break
				}
			}
		}
	}

	for i := 0; i < len(lowest_points); i++ {
		for j := 0; j < len(g.CurrentPiece.Shape); j++ {
			if g.CurrentPiece.Shape[j].Row-1 < 0 {
				return false 
			}
		}
		if g.PlayingBoard[lowest_points[i].Row-1][lowest_points[i].Col] != Pixel(0) && !ContainsShape(g.CurrentPiece.Shape, lowest_points[i]) && !ContainsSlice(lowest_points, lowest_points[i]) {
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
