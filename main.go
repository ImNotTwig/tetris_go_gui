package main

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	SideWindowHorizontalPadding = BoardWidth + Padding*2 + BorderWidth*2
	SideWindowVerticalPadding   = BoardHeight + Padding + BorderWidth - BoardHeight/4
	BorderWidth                 = 3
	BoardWidth                  = 300
	BoardHeight                 = 600
	Padding                     = BoardWidth / 20
	PixelScale                  = BoardWidth / 10 // 50
)

func run() {
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

	// looping to get keyboard inputs, line clears, and falling tetrominos
	win.Clear(color.RGBA{30, 30, 46, 255})
	imd.Color = color.RGBA{100, 100, 100, 100}
	imd.Push(pixel.V(Padding, Padding))
	imd.Push(pixel.V(BoardWidth+Padding+BorderWidth, BoardHeight+Padding+BorderWidth))
	imd.Rectangle(BorderWidth)

	imd.Push(pixel.V(SideWindowHorizontalPadding, SideWindowVerticalPadding))
	imd.Push(pixel.V((SideWindowHorizontalPadding)+BoardWidth/2, BoardHeight+Padding+BorderWidth))
	imd.Rectangle(BorderWidth)

	for !win.Closed() {
		if game.GameOver {
			break
		}
		imd.Reset()

		// TODO: GET THE NEXT PIECE TO DISPLAY PROPERLY,
		// I PROBABLY NEED TO MAKE A CONSTANT FOR DIFFERENT TETROMINOS
		if game.Current7Bag != nil && len(game.Current7Bag) > 0 {

		}

		if win.Pressed(pixelgl.KeyRight) {
			if !game.CheckIfSomethingRight(game.CurrentPiece.Shape) && time.Now().After(move_time.Add(time.Millisecond*time.Duration(50))) && !hard_dropped {
				move_time = time.Now()
				game.MoveRight()
			}
		}
		if win.JustPressed(pixelgl.KeyRight) || win.JustPressed(pixelgl.KeyLeft) || win.JustPressed(pixelgl.KeyDown) {
			lock_time = time.Now()
		}
		if win.Pressed(pixelgl.KeyLeft) {
			if !game.CheckIfSomethingLeft(game.CurrentPiece.Shape) && time.Now().After(move_time.Add(time.Millisecond*time.Duration(50))) && !hard_dropped {
				move_time = time.Now()
				game.MoveLeft()
			}
		}
		if win.Pressed(pixelgl.KeyDown) {
			can_drop = game.GravityDrop()
			if !can_drop && time.Now().After(lock_time.Add(time.Millisecond*time.Duration(100))) && !hard_dropped {
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
			game.RotateClockWise()
		}
		// setting all the pixels
		for i := 0; i < 24; i++ {
			for j := 0; j < 10; j++ {
				if i < 20 {
					imd.Color = pixel.ToRGBA(Tetro(game.PlayingBoard[i][j]).TetroToColor())
				} else {
					imd.Color = pixel.ToRGBA(color.Transparent)
				}
				imd.Push(pixel.V(float64(PixelScale*j+Padding+(BorderWidth*2)), float64(PixelScale*i+Padding+(BorderWidth*2))))

				imd.Push(pixel.V(float64((PixelScale*j)+PixelScale+10+BorderWidth/2), float64((PixelScale*i)+PixelScale+10+BorderWidth/2)))

				imd.Rectangle(0)
			}
		}
		// if now is after the move timer
		if time.Now().After(drop_time.Add(time.Millisecond * time.Duration(game.FallingSpeedMillis))) {
			can_drop = game.GravityDrop()
			drop_time = time.Now()
			if !can_drop && time.Now().After(lock_time.Add(time.Millisecond*time.Duration(200))) {
				for i := 0; i < len(game.PlayingBoard[20]); i++ {
					if game.PlayingBoard[20][i] != Pixel(0) && time.Now().After(lock_time.Add(time.Millisecond*time.Duration(100))) {
						game.GameOver = true
					}
				}
				game.SetNextTetroFromBag()
				hard_dropped = false
				lock_time = time.Now()
			}
		}
		// checking for lines that need to be cleared
		for i := 0; i < len(game.PlayingBoard); i++ {
			if !can_drop {
				line_cleared := true
				var starting_line int
				for j := 0; j < len(game.PlayingBoard[i]); j++ {
					starting_line = i
					if game.PlayingBoard[i][j] == Pixel(0) {
						line_cleared = false
						break
					}
				}

				if line_cleared {
					for j := 0; j < len(game.PlayingBoard[i]); j++ {
						game.PlayingBoard[i][j] = Pixel(0)
					}
					for j := starting_line; j < len(game.PlayingBoard)-starting_line; j++ {
						if j+1 < 20 {
							var line_pixel_list []Pixel
							for h := 0; h < len(game.PlayingBoard[j+1]); h++ {
								line_pixel_list = append(line_pixel_list, game.PlayingBoard[j+1][h])
							}
							for h := 0; h < len(game.PlayingBoard[j]); h++ {

								game.PlayingBoard[j][h] = line_pixel_list[h]

							}
						}
					}
				}
			}
		}

		imd.Draw(win)
		win.Update()
		imd.Clear()
	}

}

func main() {
	pixelgl.Run(run)
}
