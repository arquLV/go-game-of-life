package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type gameState int

const (
	startup gameState = iota
	gameplay
)

// Game structure holds the general state, current board and the current step of the game
type Game struct {
	state gameState

	board *Board
	step  int
}

// Initializes a board of size `width x height` and sets
// a fraction `fill` (0.0 to 1.0) of its cells to alive
func (g *Game) randomInit(width, height int, fill float32) {
	startingBoard := Board{
		width, height,
		make([]bool, width*height),
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rand.Float32() < fill {
				startingBoard.set(x, y, true)
			}
		}
	}

	g.board = &startingBoard
	g.step = 0
}

// Advances the current gameboard to the next step
func (g *Game) nextStep() bool {
	curBoard := g.board
	nextBoard := curBoard.createEmpty()

	gameAlive := false

	for y := 0; y < curBoard.height; y++ {
		for x := 0; x < curBoard.width; x++ {
			curCell, _ := curBoard.get(x, y)
			numLiveNeighbors := curBoard.countLiveNeighbors(x, y)

			if curCell == true && (numLiveNeighbors == 2 || numLiveNeighbors == 3) {
				nextBoard.set(x, y, true)
				gameAlive = true
			} else if curCell == false && numLiveNeighbors == 3 {
				nextBoard.set(x, y, true)
				gameAlive = true
			}
		}
	}

	g.board = nextBoard
	g.step++

	return gameAlive
}

func (g *Game) printState() {
	curBoard := g.board
	for y := 0; y < curBoard.height; y++ {
		for x := 0; x < curBoard.width; x++ {
			curCell, _ := curBoard.get(x, y)
			if curCell == true {
				fmt.Print("o")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

// Executes given function **op** for each live cell on the board, providing its x and y coords
func (g *Game) forEachLiveCell(op func(int, int)) {
	curBoard := g.board
	for y := 0; y < curBoard.height; y++ {
		for x := 0; x < curBoard.width; x++ {
			curCell, _ := curBoard.get(x, y)
			if curCell == true {
				op(x, y)
			}
		}
	}
}

func startGameloop() {
	const FrameTime = 250 // ms

	winWidth := 800
	winHeight := 800

	// Creates the Pixel window for rendering with the given width/height.
	windowConfig := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, float64(winWidth), float64(winHeight)),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(windowConfig)
	if err != nil {
		panic(err)
	}

	// Actual size of the gameboard (number of cells in each dimension)
	width := 400
	height := 400

	game := Game{state: startup}
	go game.randomInit(width, height, 0.3)

	imd := imdraw.New(nil)
	imd.Color = color.RGBA{0xff, 0x00, 0x00, 0xff}

	scaleX := winWidth / width
	scaleY := winHeight / height

	// time.Sleep(1 * time.Second)

	prevUpdate := time.Now()

	evHandler := NewEventHandler(win)
	evHandler.OnClick(pixel.R(0, 0, 300, 300), func() {
		fmt.Println("yolooo")
	})

	for !win.Closed() {

		dt := time.Now().Sub(prevUpdate)

		switch game.state {
		case startup:
			elapsed := dt.Milliseconds()
			if elapsed >= 2000 {
				game.state = gameplay
			} else {
				smoothColor := uint8(255 * float64(elapsed) / 2000)
				win.Clear(color.RGBA{smoothColor, smoothColor, smoothColor, 0xff})
			}

		case gameplay:
			if dt.Milliseconds() >= FrameTime {
				imd.Clear()

				game.forEachLiveCell(func(x, y int) {
					imd.Push(pixel.V(float64(x*scaleX), float64(y*scaleY)), pixel.V(float64((x+1)*scaleX), float64((y+1)*scaleY)))
					imd.Rectangle(0)
				})

				game.nextStep()

				win.Clear(color.RGBA{0xff, 0xff, 0xff, 0xff})
				imd.Draw(win)

				prevUpdate = time.Now()
			}
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(startGameloop)
}
