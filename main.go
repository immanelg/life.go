package main

import (
	"fmt"
)

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

type Cell struct {
	x int
	y int
	alive bool
}

type Matrix struct {
	cells []Cell
	rows, cols int
}

func newMatrix(rows, cols int) {
	m := Matrix{cells: []Cell{}, rows: rows, cols: cols}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			m.cells = append(m.cells, Cell{i, j, false})
		}
	}
}

func (m *Matrix) get(x, y int) (cell Cell, ok bool) {
	if x > m.rows || x < 0 || y > m.cols || y < 0 {
		return 
	}
	return m.cells[x * m.cols + y], true
}

func (m *Matrix) tick(x, y int) {

}

var TermX, TermY int

func main() {
	fmt.Println("hello, world!")
}
