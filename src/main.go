package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Board struct {
	width, height int
	state         []bool
}

func (b *Board) CreateEmpty() *Board {
	newState := make([]bool, b.width*b.height)
	newBoard := Board{
		b.width, b.height,
		newState,
	}

	return &newBoard
}

func (b *Board) Get(x, y int) (bool, error) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return false, errors.New("Out of bounds")
	}

	idx := (y * b.width) + x
	return b.state[idx], nil
}

func (b *Board) Set(x, y int, val bool) error {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return errors.New("Out of bounds")
	}

	idx := (y * b.width) + x
	b.state[idx] = val
	return nil
}

func (b *Board) CountLiveNeighbors(x, y int) int {
	numLive := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			neighbor, _ := b.Get(x+dx, y+dy)
			if neighbor == true {
				numLive++
			}
		}
	}

	return numLive
}

type Game struct {
	board *Board
	step  int
}

func (g *Game) RandomInit(width, height int, fill float32) {
	startingBoard := Board{
		width, height,
		make([]bool, width*height),
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rand.Float32() < fill {
				startingBoard.Set(x, y, true)
			}
		}
	}

	g.board = &startingBoard
	g.step = 0
}

func (g *Game) NextStep() bool {
	curBoard := g.board
	nextBoard := curBoard.CreateEmpty()

	gameAlive := false

	for y := 0; y < curBoard.height; y++ {
		for x := 0; x < curBoard.width; x++ {
			curCell, _ := curBoard.Get(x, y)
			numLiveNeighbors := curBoard.CountLiveNeighbors(x, y)

			if curCell == true && (numLiveNeighbors == 2 || numLiveNeighbors == 3) {
				nextBoard.Set(x, y, true)
				gameAlive = true
			} else if curCell == false && numLiveNeighbors == 3 {
				nextBoard.Set(x, y, true)
				gameAlive = true
			}
		}
	}

	g.board = nextBoard
	g.step++

	return gameAlive
}

func (g *Game) PrintState() {
	curBoard := g.board
	for y := 0; y < curBoard.height; y++ {
		for x := 0; x < curBoard.width; x++ {
			curCell, _ := curBoard.Get(x, y)
			if curCell == true {
				fmt.Print("o")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	width := 60
	height := 60

	game := Game{}
	game.RandomInit(width, height, 0.1)

	gameAlive := true

	for gameAlive == true {
		game.PrintState()
		gameAlive = game.NextStep()

		time.Sleep(250 * time.Millisecond)
	}

	fmt.Println("yo")
}
