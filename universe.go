package main

import (
	"math/rand"
)

type Universe struct {
	cells [][]bool
	w, h  int
}

func NewUniverse(w, h int) Universe {
	m := Universe{cells: make([][]bool, h), w: w, h: h}
	for i := 0; i < m.h; i++ {
		m.cells[i] = make([]bool, w)
	}
	return m
}

func (univ *Universe) neighbours(x, y int) (alive int) {
	for a := -1; a <= 1; a++ {
		for b := -1; b <= 1; b++ {
			if a == 0 && b == 0 {
				continue
			}
			if univ.Exists(x+a, y+b) && univ.Alive(x+a, y+b) {
				alive++
			}
		}
	}
	return
}

func (univ *Universe) Exists(x, y int) bool {
	return x < univ.h && x >= 0 && y < univ.w && y >= 0
}

func (univ *Universe) Alive(x, y int) bool {
	return univ.cells[x][y]
}

func (univ *Universe) Set(x, y int, alive bool) {
	univ.cells[x][y] = alive
}

func (univ *Universe) Randomize() {
	for i := 0; i < univ.h; i++ {
		for j := 0; j < univ.w; j++ {
			univ.cells[i][j] = rand.Intn(5) == 1
		}
	}
}

func (univ *Universe) Tick() {
	after := NewUniverse(univ.w, univ.h)

	for i := 0; i < univ.h; i++ {
		for j := 0; j < univ.w; j++ {
			neighbours := univ.neighbours(i, j)
			// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			// Any live cell with two or three live neighbours lives on to the next generation.
			// Any live cell with more than three live neighbours dies, as if by overpopulation.
			// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			if univ.Alive(i, j) {
				after.Set(i, j, neighbours == 2 || neighbours == 3)
			} else {
				after.Set(i, j, neighbours == 3)
			}
		}
	}
	univ.cells = after.cells
}
