package main

import (
	// "fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

type Cell struct {
	x     int
	y     int
	alive bool
}

type Matrix struct {
	cells []Cell
	w, h  int
}

func NewMatrix(w, h int) Matrix {
	m := Matrix{cells: []Cell{}, w: w, h: h}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			m.cells = append(m.cells, Cell{i, j, false})
		}
	}
	return m
}

func (m *Matrix) exists(x, y int) bool {
	return x < m.w && x >= 0 && y < m.h && y >= 0
}

func (m *Matrix) get(x, y int) *Cell {
	if !m.exists(x, y) {
		log.Panicf("get: key does not exist: (%d, %d)", x, y)
	}
	return &m.cells[x*m.h+y]
}

func (m *Matrix) Live(x, y int) {
	m.get(x, y).alive = true
}

func (m *Matrix) neighbours(x, y int) (alive int) {
	if !m.exists(x, y) {
		log.Panicf("neighbours: key does not exist: (%d, %d)", x, y)
	}
	for a := -1; a <= 1; a++ {
		for b := -1; b <= 1; b++ {
			if a == 0 && b == 0 {
				continue
			}
			if m.exists(x+a, y+b) && m.get(x+a, y+b).alive {
				alive++
			}
		}
	}
	return
}

func (m *Matrix) Randomize() {
	for i := range m.cells {
		m.cells[i].alive = rand.Intn(8) == 1 
	}
}

func (m *Matrix) Tick() {
	newCells := make([]Cell, len(m.cells))

	for i, cell := range m.cells {
		c := cell
		neighbours := m.neighbours(cell.x, cell.y)

		// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
		// Any live cell with two or three live neighbours lives on to the next generation.
		// Any live cell with more than three live neighbours dies, as if by overpopulation.
		// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
		if cell.alive {
			c.alive = neighbours == 2 || neighbours == 3
		} else {
			c.alive = neighbours == 3
		}
		newCells[i] = c
	}
	for i, c := range newCells {
		m.cells[i] = c
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
				if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEsc {
					break eventloop
				}
			}
		case <-time.After(time.Duration(time.Second/2)):
			for _, c := range matrix.cells {
				if c.alive {
					screen.SetContent(c.x, c.y, '@', nil, tcell.StyleDefault)
				} else {
					screen.SetContent(c.x, c.y, ' ', nil, tcell.StyleDefault)
				}
			}
			matrix.Tick()
			screen.Show()
		}
	}
	screen.Fini()
}
