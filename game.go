package main

import (
	"math/rand"
	"time"
)

// this function checks if a point is in a shape
func ContainsShape(sl Shape, point *Point) bool {
	for _, v := range sl {
		if v == *point {
			return true
		}
	}
	return false
}

// the main game struct for the game loop, and controlling score and whatnot, everything to do with the game is in this struct
type Game struct {
	// this playing board is a map of Points, I used this instead of a 2d array
	PlayingBoard Board

	// the current tetro thats falling
	CurrentPiece *Tetromino

	// the held piece, this is an int because we convert it to a tetro when it starts falling
	HeldPiece int

	// can we hold a piece, aka, have we already held a piece since a piece has fallen
	CanHold bool

	// the current list of tetros in the bag, which acts as a queue
	Current7Bag []*Tetromino

	// the current score of the game
	Score int

	// the total amount of lines that have been cleared, this is used to determine the level
	LinesCleared int

	// the current level of the game
	Level int

	// is the game over, have we places a tetro above the board, and does every line have a taken pixel
	GameOver bool

	// is the game paused
	Paused bool

	// how many milliseconds should it take for the piece to fall a pixel,
	// this is divided by the level so when we go up in level the speed of the falling pieces also goes up
	FallingSpeedMillis int
}

// returns a new game with defaults
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
		FallingSpeedMillis: 600,
	}
}

// generates a new 7bag, this is used when the game is started, and when the bag is empty
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

// gets the next tetro from the bag and sets it as the current tetro, then pops it from the bag
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

// gets a random tetro, this is seperate from the 7bag
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

// checks if something is under the current tetro
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

// checks if something is to the right of the current tetro
func (g *Game) CheckIfSomethingRight() bool {
	for i := 0; i < len(g.CurrentPiece.Shape); i++ {
		if g.CurrentPiece.Shape[i].Col+1 < WidthOfBoardInPixels {
			if g.PlayingBoard[Point{g.CurrentPiece.Shape[i].Row, g.CurrentPiece.Shape[i].Col + 1}] != Pixel(0) &&
				!ContainsShape(g.CurrentPiece.Shape, &Point{Row: g.CurrentPiece.Shape[i].Row, Col: g.CurrentPiece.Shape[i].Col + 1}) {
				return true
			}
		}
	}
	return false
}

// checks if something is to the left of the current tetro
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

// moves the current tetro down
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

// moves the current tetro right
func (g *Game) MoveRight() bool {
	for j := 0; j < len(g.CurrentPiece.Shape); j++ {
		if g.CurrentPiece.Shape[j].Col+1 >= WidthOfBoardInPixels {
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

// moves the current tetro to the left
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

// rotates the falling piece clockwise, if it can
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
		if retShape[i].Col >= WidthOfBoardInPixels {
			sub := retShape[i].Col
			for sub >= WidthOfBoardInPixels {
				for j := 0; j < len(retShape); j++ {
					retShape[j].Col -= 1
				}
				sub--
			}
			break
		}
		if retShape[i].Row >= HeightOfBoardInPixels {
			sub := retShape[i].Row
			for sub >= HeightOfBoardInPixels {
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

// check for lines that should be cleared
func (game *Game) check_lines() bool {
	lines := make([]int, 0)
	for i := 0; i < len(game.PlayingBoard); i++ {
		line_cleared = true
		for j := 0; j < WidthOfBoardInPixels; j++ {
			if game.PlayingBoard[Point{i, j}] == Pixel(0) {
				line_cleared = false
			}
		}
		if line_cleared {
			lines = append(lines, i)
		}
	}

	if len(lines) > 0 {
		line_cleared = true
		for i := lines[len(lines)-1] + 1; i < HeightOfBoardInPixels; i++ {
			if i < 21 {
				a := 0
				for j := 0; j < WidthOfBoardInPixels; j++ {
					if i-len(lines)+a > -1 {
						game.PlayingBoard[Point{i - len(lines), j}] = game.PlayingBoard[Point{i, j}]
					}
				}
				for j := 0; j < WidthOfBoardInPixels; j++ {
					if i-len(lines) > -1 {
						game.PlayingBoard[Point{i, j}] = Pixel(0)
					}
				}
			}
		}
	}

	game.LinesCleared += len(lines)

	game.Level = int(game.LinesCleared / 10)

	if game.Level == 0 {
		game.Level = 1
	}

	switch len(lines) {
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

// swap the current piece with the held piece
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
