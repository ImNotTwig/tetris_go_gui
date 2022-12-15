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
	SideWindowHorizontalPadding = BoardWidth + Padding*2 + BorderWidth*2
	SideWindowVerticalPadding   = BoardHeight + Padding + BorderWidth - BoardHeight/4
	BorderWidth                 = 3
	BoardWidth                  = 300
	BoardHeight                 = BoardWidth * 2
	Padding                     = BoardWidth / 20
	PixelScale                  = BoardWidth / 10 // 50
)

var (
	line_cleared bool

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
		atlas        = text.NewAtlas(face, text.ASCII)
		game         = NewGame()
		imd          = imdraw.New(nil)
		can_drop     = false
		hard_dropped = false
		lock_time    time.Time
		drop_time    time.Time
		move_time    time.Time
	)

	game.GenerateNewBag()
	game.SetNextTetroFromBag()

	for !win.Closed() {
		WidthSubForFullScreen := (win.Bounds().W() / 2) - (BoardWidth / 2)
		HeightSubForFullScreen := (win.Bounds().H() / 2) - (BoardHeight / 2)

		if win.JustPressed(pixelgl.KeyEscape) {
			game.Paused = !game.Paused
		}

		if !game.Paused {
			if game.GameOver {
				break
			}
			imd.Reset()

			if !line_cleared && !hard_dropped {
				// if any of the movement keys were just pressed set the lock_time to when that key was pressed
				if win.JustPressed(pixelgl.KeyRight) ||
					win.JustPressed(pixelgl.KeyLeft) ||
					win.JustPressed(pixelgl.KeyDown) {

					lock_time = time.Now()
				}

				if win.Pressed(pixelgl.KeyRight) {
					if !game.CheckIfSomethingRight() &&
						time.Now().After(move_time.Add(time.Millisecond*time.Duration(75))) {

						move_time = time.Now()
						game.MoveRight()
					}
				}
				if win.Pressed(pixelgl.KeyLeft) {
					if !game.CheckIfSomethingLeft() &&
						time.Now().After(move_time.Add(time.Millisecond*time.Duration(75))) {

						move_time = time.Now()
						game.MoveLeft()
					}
				}
				if win.Pressed(pixelgl.KeyDown) {
					can_drop = game.GravityDrop()
					if !can_drop &&
						time.Now().After(lock_time.Add(time.Millisecond*time.Duration(100))) {

						can_drop = false
					}
				}
				if win.JustPressed(pixelgl.KeySpace) {
					for game.GravityDrop() {
					}
					can_drop = false
					hard_dropped = true
				}
				if win.JustPressed(pixelgl.KeyUp) {
					if game.RotateClockWise() &&
						time.Now().After(lock_time.Add(time.Millisecond*time.Duration(75))) &&
						!hard_dropped {

						lock_time = time.Now()
					}
				}
				if win.JustPressed(pixelgl.KeyC) {
					if time.Now().After(move_time.Add(time.Millisecond * time.Duration(75))) {
						game.HoldTetro()
						move_time = time.Now()
						lock_time = time.Now()
						game.CanHold = false
					}
				}
			} else if line_cleared && hard_dropped {
				line_cleared = false
				hard_dropped = false
				lock_time = time.Now()
				drop_time = time.Now()
				game.SetNextTetroFromBag()
				game.CanHold = true
			}
			// if now is after the move timer
			if time.Now().After(drop_time.Add(time.Millisecond * time.Duration(game.FallingSpeedMillis))) {
				if !line_cleared {
					can_drop = game.GravityDrop()
					drop_time = time.Now()
				}
				if !can_drop && time.Now().After(lock_time.Add(time.Millisecond*time.Duration(200))) {
					game.CanHold = true
					all_lines_filled := 0
					for i := 0; i < len(game.PlayingBoard); i++ {
						for j := 0; j < 10; j++ {
							if game.PlayingBoard[Point{i, j}] != Pixel(0) {
								all_lines_filled = i
							}
						}
						if all_lines_filled >= 20 &&
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
			line_cleared = false
			if !can_drop {
				line_cleared = game.check_lines()
			}
		}

		for i := 0; i < 24; i++ {
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
		for i := 0; i < 24; i++ {
			for j := 0; j < 10; j++ {
				if ContainsShape(ghost_tetro, &Point{i, j}) && !ContainsShape(game.CurrentPiece.Shape, &Point{i, j}) {
					imd.Color = pixel.ToRGBA(Tetro(8).TetroToColor())
				} else if i < 20 {
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
		txt := text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale+HeightSubForFullScreen)), atlas)
		fmt.Fprint(txt, "Next")
		txt.Draw(win, pixel.IM)

		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale)), atlas)
		fmt.Fprint(txt, "Score")
		txt.Draw(win, pixel.IM)
		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale-PixelScale)), atlas)
		fmt.Fprint(txt, strconv.Itoa(game.Score))
		txt.Draw(win, pixel.IM)

		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale-PixelScale-PixelScale)), atlas)
		fmt.Fprint(txt, "Level")
		txt.Draw(win, pixel.IM)
		txt = text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale+WidthSubForFullScreen), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale-PixelScale-PixelScale-PixelScale)), atlas)
		fmt.Fprint(txt, strconv.Itoa(game.Level))
		txt.Draw(win, pixel.IM)

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

		imd.Draw(win)
		win.Update()
		win.Clear(color.RGBA{30, 30, 46, 255})
		imd.Clear()
	}
}

func main() {
	pixelgl.Run(run)
}
