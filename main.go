package main

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	BorderWidth = 3
	BoardWidth  = 500
	BoardHeight = 1000
	Padding     = 20
	PixelScale  = BoardWidth / 10 // 50
) // Padding + PixelScale = 70

func run() {
	game := NewGame()

	cfg := pixelgl.WindowConfig{
		Title:  "Tetris",
		Bounds: pixel.R(0, 0, 1920, 1080),
		VSync:  false,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	// main playing board
	imd.Color = pixel.ToRGBA(color.RGBA{100, 100, 100, 100})
	imd.Push(pixel.V(Padding, Padding))
	imd.Push(pixel.V(BoardWidth+Padding+BorderWidth, BoardHeight+Padding+BorderWidth))
	imd.Rectangle(BorderWidth)

	game.GenerateNewBag()
	can_drop := false

	// looping to get keyboard inputs, line clears, and falling tetrominos
	for !win.Closed() {
		time.Sleep(time.Duration(game.FallingSpeedMillis) * time.Millisecond)
		win.Clear(color.RGBA{30, 30, 46, 255})
		for i := 0; i < 24; i++ {
			for j := 0; j < 10; j++ {
				if i < 20 {
					imd.Color = pixel.ToRGBA(Tetro(game.PlayingBoard[i][j]).TetroToColor())
				} else {
					imd.Color = pixel.ToRGBA(color.Transparent)
				}
				imd.Push(pixel.V(float64(PixelScale*j+Padding+(BorderWidth/2)), float64(PixelScale*i+Padding+(BorderWidth/2))))
				imd.Push(pixel.V(float64((PixelScale*j)+PixelScale+Padding+(BorderWidth/2)), float64((PixelScale*i)+PixelScale+Padding+(BorderWidth/2))))
				imd.Rectangle(0)
			}
		}
		if !can_drop {
			game.SetNextTetroFromBag()
		}
		can_drop = game.GravityDrop()

		imd.Draw(win)
		win.Update()

	}
}

func main() {
	pixelgl.Run(run)
}
