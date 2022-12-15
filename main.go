package main

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	// padding to display the next piece next to the board (horizontal)
	SideWindowHorizontalPadding = BoardWidth + Padding*2 + BorderWidth*2

	// padding to display the next piece next to the board (vertical)
	SideWindowVerticalPadding = BoardHeight + Padding + BorderWidth - BoardHeight/4

	// how wide is the border around the board
	BorderWidth = 3

	// how wide the board is in screen pixels
	BoardWidth = 300

	// How tall is the board in screen pixels
	BoardHeight = BoardWidth * 2

	// padding for extra space between in-game pixels
	Padding = BoardWidth / 20

	// how many screen pixels wide/tall is a in-game pixel
	PixelScale = BoardWidth / WidthOfBoardInPixels

	// How many pixels wide is the board.
	// pixels being the pixels defined in this game, not screen pixels
	WidthOfBoardInPixels = 10

	// How many pixels tall is the board.
	// pixels being the pixels defined in this game, not screen pixels
	HeightOfBoardInPixels = 24

	// How many pixels tall is the non-hidden part of the board.
	// pixels being the pixels defined in this game, not screen pixels
	NonHiddenPixelHeight = 20
)

var (
	// line_cleared determines if a line was cleared or not
	line_cleared bool

	// the ghost_tetro contains the coordinates of the ghost pixel on screen
	ghost_tetro Shape
)

func run() {
	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 20,
	})

	monitor_width, monitor_height := pixelgl.PrimaryMonitor().PhysicalSize()
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Tetris",
		Bounds: pixel.R(0, 0, monitor_width, monitor_height),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	defer win.Destroy()

	var (
		// atlas for displaying text on the screen
		atlas = text.NewAtlas(face, text.ASCII)

		// main game struct
		game = NewGame()

		// imdraw struct to draw shapes on the screen
		imd = imdraw.New(nil)

		// can_drop determines if we should lock the piece
		can_drop = false

		// hard_dropped determines if we just hard dropped a piece
		hard_dropped = false

		// lock time determines how many milliseconds shouldve passed before we lock the piece
		lock_time time.Time

		// drop time determines how many milliseconds shouldve passed before we drop the piece a pixel
		drop_time time.Time

		// move time determines how many milliseconds shouldve passed before we can move the piece
		move_time time.Time
	)

	game.GenerateNewBag()
	game.SetNextTetroFromBag()

	for !win.Closed() {
		// calculating the center of the screen every time we loop
		WidthSubForFullScreen := (win.Bounds().W() / 2) - (BoardWidth / 2)
		HeightSubForFullScreen := (win.Bounds().H() / 2) - (BoardHeight / 2)

		// checking if we paused/unpaused the game
		if win.JustPressed(pixelgl.KeyEscape) {
			game.Paused = !game.Paused
		}
		// if the game is not paused go through the normal game loop
		if !game.Paused {
			// if we lost, break this loop to close the window
			if game.GameOver {
				break
			}
			// resetting the graphics
			imd.Reset()

			// if a line was not just cleared and we didnt just hard drop a piece
			// we check these so we cant like rotate a piece if its supposed to be locked in place
			if !line_cleared && !hard_dropped {
				// if any of the movement keys were just pressed set the lock_time to when that key was pressed
				if win.JustPressed(pixelgl.KeyRight) ||
					win.JustPressed(pixelgl.KeyLeft) ||
					win.JustPressed(pixelgl.KeyDown) {

					lock_time = time.Now()
				}

				// if we pressed right, move the piece right if it can
				if win.Pressed(pixelgl.KeyRight) {
					if !game.CheckIfSomethingRight() &&
						time.Now().After(move_time.Add(time.Millisecond*time.Duration(75))) {

						move_time = time.Now()
						game.MoveRight()
					}
				}
				// if we pressed left, move the piece left if it can
				if win.Pressed(pixelgl.KeyLeft) {
					if !game.CheckIfSomethingLeft() &&
						time.Now().After(move_time.Add(time.Millisecond*time.Duration(75))) {

						move_time = time.Now()
						game.MoveLeft()
					}
				}
				// if we're pressing down, start falling down faster
				if win.Pressed(pixelgl.KeyDown) {
					can_drop = game.GravityDrop()
					if !can_drop &&
						time.Now().After(lock_time.Add(time.Millisecond*time.Duration(200))) {

						can_drop = false
					}
				}
				// if we just pressed space, hard drop
				if win.JustPressed(pixelgl.KeySpace) {
					for game.GravityDrop() {
					}
					can_drop = false
					hard_dropped = true
				}
				// if we just pressed up, rotate the piece if it can
				if win.JustPressed(pixelgl.KeyUp) {
					if game.RotateClockWise() &&
						time.Now().After(lock_time.Add(time.Millisecond*time.Duration(75))) &&
						!hard_dropped {

						lock_time = time.Now()
					}
				}
				// if we just pressed C then hold the current piece
				if win.JustPressed(pixelgl.KeyC) {
					if time.Now().After(move_time.Add(time.Millisecond * time.Duration(75))) {
						game.HoldTetro()
						move_time = time.Now()
						lock_time = time.Now()
						game.CanHold = false
					}
				}
				// if a line was cleared and we just hard dropped then reset timers and set a new current piece
			} else if line_cleared && hard_dropped {
				line_cleared = false
				hard_dropped = false
				lock_time = time.Now()
				drop_time = time.Now()
				game.SetNextTetroFromBag()
				game.CanHold = true
			}
			// if now is after the move timer, then move the piece down naturally
			if time.Now().After(drop_time.Add(time.Millisecond * time.Duration(game.FallingSpeedMillis/game.Level))) {
				// if a line has not been cleared, drop the piece
				if !line_cleared {
					can_drop = game.GravityDrop()
					drop_time = time.Now()
				}
				// otherwise check if the game should end, and if it shouldnt set a new piece
				if !can_drop && time.Now().After(lock_time.Add(time.Millisecond*time.Duration(200))) {
					game.CanHold = true
					// we are checking if every row on the board has a pixel thats taken
					// if thats true then the game is over, because we've filled the board
					all_lines_filled := 0
					for i := 0; i < len(game.PlayingBoard); i++ {
						for j := 0; j < WidthOfBoardInPixels; j++ {
							if game.PlayingBoard[Point{i, j}] != Pixel(0) {
								all_lines_filled = i
							}
						}
						if all_lines_filled >= NonHiddenPixelHeight &&
							!line_cleared &&
							time.Now().After(lock_time.Add(time.Millisecond*time.Duration(200))) {
							game.GameOver = true
						}
					}

					game.SetNextTetroFromBag()
					hard_dropped = false
					lock_time = time.Now()
				}

			}
			// checking if any lines were cleared
			line_cleared = false
			if !can_drop {
				line_cleared = game.check_lines()
			}
		}

		// getting the coordinates of the ghost tetro
		for i := 0; i < HeightOfBoardInPixels; i++ {
			shape := make(Shape, 0)
			for j := 0; j < len(game.CurrentPiece.Shape); j++ {
				shape = append(shape, Point{Col: game.CurrentPiece.Shape[j].Col, Row: game.CurrentPiece.Shape[j].Row - i})
			}

			if game.CheckIfSomethingUnder(&shape) {
				ghost_tetro = shape
				break
			}
		}

		// setting all the pixels
		for i := 0; i < HeightOfBoardInPixels; i++ {
			for j := 0; j < WidthOfBoardInPixels; j++ {
				if ContainsShape(ghost_tetro, &Point{i, j}) && !ContainsShape(game.CurrentPiece.Shape, &Point{i, j}) {
					imd.Color = pixel.ToRGBA(Tetro(8).TetroToColor())
				} else if i < NonHiddenPixelHeight {
					imd.Color = pixel.ToRGBA(Tetro(game.PlayingBoard[Point{i, j}]).TetroToColor())
				} else {
					imd.Color = pixel.ToRGBA(color.Transparent)
				}
				imd.Push(pixel.V(float64(PixelScale*j+Padding+(BorderWidth*2)+int(WidthSubForFullScreen)), float64(PixelScale*i+Padding+(BorderWidth*2)+int(HeightSubForFullScreen))))

				imd.Push(pixel.V(float64((PixelScale*j)+PixelScale+Padding/2+BorderWidth/2+int(WidthSubForFullScreen)), float64((PixelScale*i)+PixelScale+Padding/2+BorderWidth/2+int(HeightSubForFullScreen))))

				imd.Rectangle(0)
			}
		}

		// showing the border of the board
		imd.Color = color.RGBA{100, 100, 100, 100}
		imd.Push(pixel.V(Padding+WidthSubForFullScreen, Padding+HeightSubForFullScreen))
		imd.Push(pixel.V(BoardWidth+Padding+BorderWidth+WidthSubForFullScreen, BoardHeight+Padding+BorderWidth+HeightSubForFullScreen))
		imd.Rectangle(BorderWidth)

		// checking if the bag is emtpy so we can show the next piece
		if len(game.Current7Bag) < 1 || game.Current7Bag == nil {
			game.GenerateNewBag()
		}

		// showing the next piece
		shape := game.Current7Bag[0].Tetro.TetroToNewShape()
		for i := 0; i < len(shape); i++ {
			shape[i].Col -= 4
			shape[i].Row -= 22
		}
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				imd.Color = pixel.ToRGBA(game.Current7Bag[0].Tetro.TetroToColor())
				imd.Push(pixel.V(float64(SideWindowHorizontalPadding+shape[i].Col*PixelScale+PixelScale+Padding+int(WidthSubForFullScreen)), float64(SideWindowVerticalPadding+PixelScale+shape[i].Row*PixelScale+Padding+int(HeightSubForFullScreen))))
				imd.Push(pixel.V(float64(PixelScale+PixelScale+SideWindowHorizontalPadding+shape[i].Col*PixelScale+int(WidthSubForFullScreen)), float64(PixelScale+PixelScale+SideWindowVerticalPadding+shape[i].Row*PixelScale+int(HeightSubForFullScreen))))
				imd.Rectangle(0)
			}
		}
		// displaying next for the next piece
		txt := text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale+HeightSubForFullScreen)), atlas)
		fmt.Fprint(txt, "Next")
		txt.Draw(win, pixel.IM)

		// displaying the score text
		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale)), atlas)
		fmt.Fprint(txt, "Score")
		txt.Draw(win, pixel.IM)
		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale-PixelScale)), atlas)
		fmt.Fprint(txt, strconv.Itoa(game.Score))
		txt.Draw(win, pixel.IM)

		// displaying the level text
		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale-PixelScale-PixelScale)), atlas)
		fmt.Fprint(txt, "Level")
		txt.Draw(win, pixel.IM)
		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale-PixelScale-PixelScale-PixelScale)), atlas)
		fmt.Fprint(txt, strconv.Itoa(game.Level))
		txt.Draw(win, pixel.IM)

		// showing the held piece
		if game.HeldPiece != 0 {
			shape := Tetro(game.HeldPiece).TetroToNewShape()
			for i := 0; i < len(shape); i++ {
				shape[i].Col -= 4
				shape[i].Row -= 22
			}

			for i := 0; i < 4; i++ {
				for j := 0; j < 4; j++ {
					imd.Color = pixel.ToRGBA(Tetro(game.HeldPiece).TetroToColor())
					imd.Push(pixel.V(float64(-(SideWindowHorizontalPadding/2)+(shape[i].Col*PixelScale)+(PixelScale+Padding)+int(WidthSubForFullScreen)), float64((SideWindowVerticalPadding)+(PixelScale+Padding)+(shape[i].Row*PixelScale)+int(HeightSubForFullScreen))))
					imd.Push(pixel.V(float64(PixelScale+PixelScale-SideWindowHorizontalPadding/2+shape[i].Col*PixelScale+int(WidthSubForFullScreen)), float64(PixelScale+PixelScale+SideWindowVerticalPadding+shape[i].Row*PixelScale+int(HeightSubForFullScreen))))
					imd.Rectangle(0)
				}
			}
			txt := text.New(pixel.V(float64(-SideWindowHorizontalPadding/2+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale+HeightSubForFullScreen)), atlas)
			fmt.Fprint(txt, "Held")
			txt.Draw(win, pixel.IM)
		}

		// clearing the screen for the next frame
		imd.Draw(win)
		win.Update()
		win.Clear(color.RGBA{30, 30, 46, 255})
		imd.Clear()
	}
}

func main() {
	pixelgl.Run(run)
}
