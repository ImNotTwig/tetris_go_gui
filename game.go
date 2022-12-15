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

	HeldPiece int

	CanHold bool

	Current7Bag []*Tetromino

	Score int

	LinesCleared int

	Level int

	GameOver bool

	Paused bool

	FallingSpeedMillis int
}

func NewGame() Game {
	return Game{
		PlayingBoard:       NewBoard(),
		CurrentPiece:       nil,
		HeldPiece:          0,
		CanHold:            true,
		Current7Bag:        nil,
		Score:              0,
		LinesCleared:       0,
		Level:              1,
		GameOver:           false,
		Paused:             false,
		FallingSpeedMillis: 300,
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
			g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col}] = Pixel(g.CurrentPiece.Tetro)
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
		g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col}] = Pixel(g.CurrentPiece.Tetro)
	}
}

func (g *Game) CheckIfSomethingUnder(s *Shape) bool {
	if s == nil {
		s = &g.CurrentPiece.Shape
	}
	for i := 0; i < len(*s); i++ {
		if (*s)[i].Row != 0 {
			if g.PlayingBoard[Point{(*s)[i].Row - 1, (*s)[i].Col}] != Pixel(0) &&
				!ContainsShape(*s, &Point{Row: (*s)[i].Row - 1, Col: (*s)[i].Col}) {
				return true
			}
		} else {
			return true
		}
	}
	return false
}
func (g *Game) CheckIfSomethingRight() bool {
	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		if g.CurrentPiece.Shape[i].Col+1 < 10 {
			if g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col + 1}] != Pixel(0) &&
				!ContainsShape(g.CurrentPiece.Shape, &Point{Row: g.CurrentPiece.Shape[i].Row, Col: g.CurrentPiece.Shape[i].Col + 1}) {
				return true
			}
		}
	}
	return false
}
func (g *Game) CheckIfSomethingLeft() bool {
	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		if g.CurrentPiece.Shape[i].Col-1 > -1 {
			if g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col - 1}] != Pixel(0) &&
				!ContainsShape(g.CurrentPiece.Shape, &Point{Row: g.CurrentPiece.Shape[i].Row, Col: g.CurrentPiece.Shape[i].Col - 1}) {
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

	if g.CheckIfSomethingUnder(nil) {
		return false
	}

	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[j].Row, g.CurrentPiece.Shape[j].Col}] = Pixel(0)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[j].Row - 1, g.CurrentPiece.Shape[j].Col}] = Pixel(g.CurrentPiece.Tetro)
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

	if g.CheckIfSomethingRight() {
		return false
	}

	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[j].Row, g.CurrentPiece.Shape[j].Col}] = Pixel(0)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[j].Row, g.CurrentPiece.Shape[j].Col + 1}] = Pixel(g.CurrentPiece.Tetro)
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

	if g.CheckIfSomethingLeft() {
		return false
	}

	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[j].Row, g.CurrentPiece.Shape[j].Col}] = Pixel(0)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[j].Row, g.CurrentPiece.Shape[j].Col - 1}] = Pixel(g.CurrentPiece.Tetro)
	}
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		g.CurrentPiece.Shape[j].Col -= 1
	}
	return true
}

func (g *Game) ClearShape() {
	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col}] = Pixel(0)
	}
	g.CurrentPiece.Shape = Shape{}
}

func (g *Game) RotateClockWise() bool {
	var retShape Shape
	for i := 0; i < 4; i++ {
		retShape = append(retShape, Point{0, 0})
	}
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
			for sub <= 0 {
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
		if retShape[i].Row >= 24 {
			sub := retShape[i].Row
			for sub >= 24 {
				for j := 0; j < len(retShape); j++ {
					retShape[j].Row -= 1
				}
				sub--
			}
			break
		}
		if g.PlayingBoard[Point{retShape[i].Row, retShape[i].Col}] != Pixel(0) &&
			!ContainsShape(g.CurrentPiece.Shape, &Point{
				Row: retShape[i].Row,
				Col: retShape[i].Col,
			}) {
			return false
		}
	}

	for i := 0; i < 4; i++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col}] = Pixel(0)
	}

	g.CurrentPiece.Shape = retShape

	for i := 0; i < 4; i++ {
		g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col}] = Pixel(g.CurrentPiece.Tetro)
	}
	return true
}

func (game *Game) check_lines() bool {
	lines := make([]int, 0)
	for i := 0; i < len(game.PlayingBoard); i++ {
		line_cleared = true
		for j := 0; j < 10; j++ {
			if game.PlayingBoard[Point{i, j}] == Pixel(0) {
				line_cleared = false
			}
		}
		if line_cleared {
			lines = append(lines, i)
		}
	}

	lines_length := len(lines)

	if lines_length > 0 {
		line_cleared = true
		for i := lines[lines_length-1] + 1; i < 24; i++ {
			if i < 21 {
				a := 0
				for j := 0; j < 10; j++ {
					if i-lines_length+a > -1 {
						game.PlayingBoard[Point{i - lines_length, j}] = game.PlayingBoard[Point{i, j}]
					}
				}
				for j := 0; j < 10; j++ {
					if i-lines_length > -1 {
						game.PlayingBoard[Point{i, j}] = Pixel(0)
					}
				}
			}
		}
	}

	game.LinesCleared += lines_length

	game.Level = int(game.LinesCleared / 10)
	if game.Level == 0 {
		game.Level = 1
	}

	switch lines_length {
	case 1:
		game.Score += 40 * game.Level
	case 2:
		game.Score += 100 * game.Level
	case 3:
		game.Score += 300 * game.Level
	case 4:
		game.Score += 1200 * game.Level
	}

	return line_cleared
}

func (g *Game) HoldTetro() {
	if !g.CanHold {
		return
	}
	if g.HeldPiece != 0 {
		for i := 0; i < len(g.CurrentPiece.Shape); i++ {
			g.PlayingBoard[g.CurrentPiece.Shape[i]] = Pixel(0)
		}
		temp := g.HeldPiece
		g.HeldPiece = int(g.CurrentPiece.Tetro)
		g.CurrentPiece = &Tetromino{
			Tetro: Tetro(temp),
			Shape: Tetro(temp).TetroToNewShape(),
		}
		//g.HeldPiece = int(g.CurrentPiece.Tetro)
	} else {
		for i := 0; i < len(g.CurrentPiece.Shape); i++ {
			g.PlayingBoard[g.CurrentPiece.Shape[i]] = Pixel(0)
		}
		g.HeldPiece = int(g.CurrentPiece.Tetro)

		g.SetNextTetroFromBag()
	}

}
