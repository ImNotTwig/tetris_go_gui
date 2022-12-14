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

	cfg := pixelgl.WindowConfig{
		Title:  "Tetris",
		Bounds: pixel.R(0, 0, BoardWidth+100, BoardHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
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
		if !line_cleared && !hard_dropped {
			if win.Pressed(pixelgl.KeyRight) {
				if !game.CheckIfSomethingRight(game.CurrentPiece.Shape) && time.Now().After(move_time.Add(time.Millisecond*time.Duration(50))) {

					move_time = time.Now()
					game.MoveRight()
				}
			}
			if win.JustPressed(pixelgl.KeyRight) || win.JustPressed(pixelgl.KeyLeft) || win.JustPressed(pixelgl.KeyDown) {
				lock_time = time.Now()
			}
			if win.Pressed(pixelgl.KeyLeft) {
				if !game.CheckIfSomethingLeft(game.CurrentPiece.Shape) && time.Now().After(move_time.Add(time.Millisecond*time.Duration(50))) {
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
				imd.Push(pixel.V(float64(PixelScale*j+Padding+(BorderWidth*2)), float64(PixelScale*i+Padding+(BorderWidth*2))))

				imd.Push(pixel.V(float64((PixelScale*j)+PixelScale+Padding/2+BorderWidth/2), float64((PixelScale*i)+PixelScale+Padding/2+BorderWidth/2)))

				imd.Rectangle(0)
			}
		}
		if !can_drop {
			can_drop = game.check_lines()
		}

		imd.Color = color.RGBA{100, 100, 100, 100}
		imd.Push(pixel.V(Padding, Padding))
		imd.Push(pixel.V(BoardWidth+Padding+BorderWidth, BoardHeight+Padding+BorderWidth))
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
				imd.Push(pixel.V(float64(SideWindowHorizontalPadding+shape[i].Col*PixelScale+PixelScale+Padding), float64(SideWindowVerticalPadding+PixelScale+shape[i].Row*PixelScale+Padding)))
				imd.Push(pixel.V(float64(PixelScale+PixelScale+SideWindowHorizontalPadding+shape[i].Col*PixelScale), float64(PixelScale+PixelScale+SideWindowVerticalPadding+shape[i].Row*PixelScale)))
				imd.Rectangle(0)
			}
		}
		txt := text.New(pixel.V(float64(SideWindowHorizontalPadding+PixelScale), float64(PixelScale+PixelScale+SideWindowVerticalPadding+2*PixelScale)), atlas)
		fmt.Fprint(txt, "Next")
		txt.Draw(win, pixel.IM)

		imd.Draw(win)
		win.Update()
		win.Clear(color.RGBA{30, 30, 46, 255})
		imd.Clear()
	}
}

func main() {
	pixelgl.Run(run)
}
