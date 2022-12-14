package main

import (
	"fmt"
	"image/color"
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
)

func run() {

	ttf, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: 20,
	})
	atlas := text.NewAtlas(face, text.ASCII)

	game := NewGame()

	monitor_width, monitor_height := pixelgl.PrimaryMonitor().PhysicalSize()
	if err != nil {
		panic(err)
	}

	cfg := pixelgl.WindowConfig{
		Monitor: pixelgl.PrimaryMonitor(),
		Title:   "Tetris",
		Bounds:  pixel.R(0, 0, monitor_width, monitor_height),
		VSync:   true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		WidthSubForFullScreen  = (win.Bounds().W() / 2) - (BoardWidth / 2)
		HeightSubForFullScreen = (win.Bounds().H() / 2) - (BoardHeight / 2)
	)

	win.SetSmooth(true)

	imd := imdraw.New(nil)

	game.GenerateNewBag()
	game.SetNextTetroFromBag()
	can_drop := false

	var lock_time time.Time
	var drop_time time.Time
	var move_time time.Time

	hard_dropped := false

	for !win.Closed() {
		if game.GameOver {
			break
		}
		imd.Reset()

		if line_cleared {
			hard_dropped = false
			lock_time = time.Now()
			line_cleared = false
			drop_time = time.Now()
		}

		if !line_cleared && !hard_dropped {
			if win.Pressed(pixelgl.KeyRight) {
				if !game.CheckIfSomethingRight() && time.Now().After(move_time.Add(time.Millisecond*time.Duration(50))) {
					move_time = time.Now()
					game.MoveRight()
				}
			}
			if win.JustPressed(pixelgl.KeyRight) || win.JustPressed(pixelgl.KeyLeft) || win.JustPressed(pixelgl.KeyDown) {
				lock_time = time.Now()
			}
			if win.Pressed(pixelgl.KeyLeft) {
				if !game.CheckIfSomethingLeft() && time.Now().After(move_time.Add(time.Millisecond*time.Duration(50))) {
					move_time = time.Now()
					game.MoveLeft()
				}
			}
			if win.Pressed(pixelgl.KeyDown) {
				can_drop = game.GravityDrop()
				if !can_drop && time.Now().After(lock_time.Add(time.Millisecond*time.Duration(100))) {
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
				if game.RotateClockWise() && time.Now().After(lock_time.Add(time.Millisecond*time.Duration(50))) && !hard_dropped {
					lock_time = time.Now()
				}
			}
			if win.JustPressed(pixelgl.KeyC) {
				if time.Now().After(move_time.Add(time.Millisecond * time.Duration(50))) {
					game.HoldTetro()
					move_time = time.Now()
					lock_time = time.Now()
					game.CanHold = false
				}
			}
		} else if line_cleared && hard_dropped {
			lock_time = time.Now()
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

		// setting all the pixels
		for i := 0; i < 24; i++ {
			for j := 0; j < 10; j++ {
				if i < 20 {
					imd.Color = pixel.ToRGBA(Tetro(game.PlayingBoard[Point{i, j}]).TetroToColor())
				} else {
					imd.Color = pixel.ToRGBA(color.Transparent)
				}
				imd.Push(pixel.V(float64(PixelScale*j+Padding+(BorderWidth*2)+int(WidthSubForFullScreen)), float64(PixelScale*i+Padding+(BorderWidth*2)+int(HeightSubForFullScreen))))

				imd.Push(pixel.V(float64((PixelScale*j)+PixelScale+Padding/2+BorderWidth/2+int(WidthSubForFullScreen)), float64((PixelScale*i)+PixelScale+Padding/2+BorderWidth/2+int(HeightSubForFullScreen))))

				imd.Rectangle(0)
			}
		}
		if !can_drop {
			line_cleared = game.check_lines()
		}

		imd.Color = color.RGBA{100, 100, 100, 100}
		imd.Push(pixel.V(Padding+WidthSubForFullScreen, Padding+HeightSubForFullScreen))
		imd.Push(pixel.V(BoardWidth+Padding+BorderWidth+WidthSubForFullScreen, BoardHeight+Padding+BorderWidth+HeightSubForFullScreen))
		imd.Rectangle(BorderWidth)

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
