package main

import "errors"

type Board struct {
	width, height int
	state         []bool
}

func (b *Board) createEmpty() *Board {
	newState := make([]bool, b.width*b.height)
	newBoard := Board{
		b.width, b.height,
		newState,
	}

	return &newBoard
}

func (b *Board) get(x, y int) (bool, error) {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return false, errors.New("Out of bounds")
	}

	idx := (y * b.width) + x
	return b.state[idx], nil
}

func (b *Board) set(x, y int, val bool) error {
	if x < 0 || x >= b.width || y < 0 || y >= b.height {
		return errors.New("Out of bounds")
	}

	idx := (y * b.width) + x
	b.state[idx] = val
	return nil
}

func (b *Board) countLiveNeighbors(x, y int) int {
	numLive := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			neighbor, _ := b.get(x+dx, y+dy)
			if neighbor == true {
				numLive++
			}
		}
	}

	return numLive
}
