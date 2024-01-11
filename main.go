package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

type Matrix struct {
	cells [][]bool
	w, h  int
}

func NewMatrix(w, h int) Matrix {
	m := Matrix{cells: make([][]bool, h), w: w, h: h}
	for i := 0; i < m.h; i++ {
		m.cells[i] = make([]bool, w)
	}
	return m
}

func (m *Matrix) neighbours(x, y int) (alive int) {
	for a := -1; a <= 1; a++ {
		for b := -1; b <= 1; b++ {
			if a == 0 && b == 0 {
				continue
			}
			if m.Exists(x+a, y+b) && m.Alive(x+a, y+b) {
				alive++
			}
		}
	}
	return
}

func (m *Matrix) Exists(x, y int) bool {
	return x < m.h && x >= 0 && y < m.w && y >= 0
}

func (m *Matrix) Alive(x, y int) bool {
	return m.cells[x][y]
}

func (m *Matrix) Set(x, y int, alive bool) {
	m.cells[x][y] = alive
}

func (m *Matrix) Randomize() {
	for i := 0; i < m.h; i++ {
		for j := 0; j < m.w; j++ {
			m.cells[i][j] = rand.Intn(5) == 1
		}
	}
}

func (m *Matrix) Tick() {
	after := NewMatrix(m.w, m.h)

	for i := 0; i < m.h; i++ {
		for j := 0; j < m.w; j++ {
			neighbours := m.neighbours(i, j)
			// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
			// Any live cell with two or three live neighbours lives on to the next generation.
			// Any live cell with more than three live neighbours dies, as if by overpopulation.
			// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			if m.Alive(i, j) {
				after.Set(i, j, neighbours == 2 || neighbours == 3)
			} else {
				after.Set(i, j, neighbours == 3)
			}
		}
	}
	for i := 0; i < m.h; i++ {
		for j := 0; j < m.w; j++ {
			m.cells[i][j] = after.cells[i][j]
		}
	}
}

func main() {
	screen, _ := tcell.NewScreen()
	screen.Init()
	screen.Clear()

	TermW, TermH := screen.Size()

	matrix := NewMatrix(TermW, TermH)
	matrix.Randomize()

	events := make(chan tcell.Event)

	go func() {
		for {
			events <- screen.PollEvent()
		}
	}()

eventloop:
	for {
		select {
		case ev := <-events:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyESC {
					break eventloop
				}
			}
		case <-time.After(time.Duration(time.Second / 3)):
			for i := 0; i < matrix.h; i++ {
				for j := 0; j < matrix.w; j++ {
					if matrix.Alive(i, j) {
						screen.SetContent(j, i, tcell.RuneBlock, nil, tcell.StyleDefault)
					} else {
						screen.SetContent(j, i, ' ', nil, tcell.StyleDefault)
					}
				}
			}
		}
		matrix.Tick()
		screen.Show()
	}
	screen.Fini()
}
