package main

import (
	// "fmt"
	"log"
	// "os"

	"github.com/gdamore/tcell"
)

type Cell struct {
	x     int
	y     int
	alive bool
}

type Matrix struct {
	cells      []Cell
	w, h int
}

func newMatrix(w, h int) Matrix {
	m := Matrix{cells: []Cell{}, w: w, h: h}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			m.cells = append(m.cells, Cell{i, j, false})
		}
	}
	return m
}

func (m *Matrix) get(x, y int) (cell Cell, ok bool) {
	if x >= m.w || x < 0 || y >= m.h || y < 0 {
		return
	}
	return m.cells[x*m.h+y], true
}

func (m *Matrix) getNeighbours(x, y int) (cells []Cell) {
	if _, exists := m.get(x, y); !exists {
		log.Panicf("should exist: (%d, %d)", x, y)
	}
	for a := range []int{-1, 0, +1} {
		for b := range []int{-1, 0, +1} {
			if a == 0 && b == 0 {
				continue
			}
			cell, exists := m.get(x+a, y+b)
			if !exists {
				continue
			}
			cells = append(cells, cell)
		}
	}
	return
}

func (m *Matrix) tick(x, y int) {
	matrixAfter := newMatrix(m.w, m.h)
	for _, cell := range m.cells {
		neighbours := m.getNeighbours(cell.x, cell.y)
		aliveNeighbours := 0
		deadNeighbours := 0
		for _, n := range neighbours {
			if n.alive {
				aliveNeighbours++
			} else {
				deadNeighbours++
			}
		}
		// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
		// Any live cell with two or three live neighbours lives on to the next generation.
		// Any live cell with more than three live neighbours dies, as if by overpopulation.
		// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

		if aliveNeighbours < 2 {
			cell.alive = false
		}
		if aliveNeighbours == 2 || aliveNeighbours == 3 {
			cell.alive = true
		}
		matrixAfter.cells = append(matrixAfter.cells, cell)
	}
}

var TermH, TermW int

func main() {
	screen, _ := tcell.NewScreen()
	screen.Init()
	screen.Clear()

	TermW, TermH = screen.Size()

loop:
	for {
		screen.SetContent(0, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(1, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(2, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(3, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(4, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(5, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(6, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(7, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(8, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(9, 0, '1', nil, tcell.StyleDefault)
		screen.SetContent(10, 0, '1', nil, tcell.StyleDefault)

		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				break loop
			}
		}
	}
	screen.Fini()
}
