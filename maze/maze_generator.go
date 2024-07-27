package maze

import (
	"math"
	"math/rand"
	"time"
)

type Cell struct {
	X, Y    int
	IsWall  bool
	Visited bool
}

type Maze struct {
	Width, Height int
	Grid          [][]Cell
	Stack         []Cell
	CurrentCell   *Cell
	Start         *Cell
	End           *Cell
}

func NewMaze(width, height int) *Maze {
	m := &Maze{
		Width:  width*2 + 1,
		Height: height*2 + 1,
		Grid:   make([][]Cell, height*2+1),
		Stack:  []Cell{},
	}

	for y := 0; y < m.Height; y++ {
		m.Grid[y] = make([]Cell, m.Width)
		for x := 0; x < m.Width; x++ {
			isWall := x%2 == 0 || y%2 == 0
			m.Grid[y][x] = Cell{X: x, Y: y, IsWall: isWall}
		}
	}

	m.CurrentCell = &m.Grid[1][1]
	m.Start = &m.Grid[1][1]
	m.End = &m.Grid[m.Height-2][m.Width-2]

	return m
}

func (m *Maze) getCell(x, y int) *Cell {
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return nil
	}
	return &m.Grid[y][x]
}

func (m *Maze) getNeighbors(cell *Cell) *Cell {
	var neighbors []*Cell

	top := m.getCell(cell.X, cell.Y-2)
	right := m.getCell(cell.X+2, cell.Y)
	bottom := m.getCell(cell.X, cell.Y+2)
	left := m.getCell(cell.X-2, cell.Y)

	if top != nil && !top.Visited {
		neighbors = append(neighbors, top)
	}
	if right != nil && !right.Visited {
		neighbors = append(neighbors, right)
	}
	if bottom != nil && !bottom.Visited {
		neighbors = append(neighbors, bottom)
	}
	if left != nil && !left.Visited {
		neighbors = append(neighbors, left)
	}

	if len(neighbors) > 0 {
		if rand.Float64() < 0.75 {
			return neighbors[rand.Intn(len(neighbors))]
		}

		var maxDistance float64
		var farthestCell *Cell

		for _, neighbor := range neighbors {
			distance := math.Hypot(float64(neighbor.X-m.Start.X), float64(neighbor.Y-m.Start.Y))
			if distance > maxDistance {
				maxDistance = distance
				farthestCell = neighbor
			}
		}

		return farthestCell
	}

	return nil
}

func (m *Maze) generateMazeNotGlobal() {
	m.CurrentCell.Visited = true
	nextCell := m.getNeighbors(m.CurrentCell)

	if nextCell != nil {
		nextCell.Visited = true

		m.Stack = append(m.Stack, *m.CurrentCell)

		wallX := (m.CurrentCell.X + nextCell.X) / 2
		wallY := (m.CurrentCell.Y + nextCell.Y) / 2
		m.Grid[wallY][wallX].IsWall = false

		m.CurrentCell = nextCell
	} else if len(m.Stack) > 0 {
		if rand.Float64() < 0.4 {
			backtrackCell := m.Stack[rand.Intn(len(m.Stack))]
			m.CurrentCell = &backtrackCell
		} else {
			m.CurrentCell = &m.Stack[len(m.Stack)-1]
			m.Stack = m.Stack[:len(m.Stack)-1]
		}
	}
}

func GenerateMaze(numRows, numCols int, singlePath bool) *Maze {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	maze := NewMaze(numRows, numCols)

	for len(maze.Stack) > 0 || !maze.CurrentCell.Visited {
		maze.generateMazeNotGlobal()
	}

	if !singlePath {
		for i := 0; i < numRows*numCols/10; i++ {
			x := r.Intn(numRows-2) + 1
			y := r.Intn(numCols-2) + 1
			maze.Grid[x][y].IsWall = false
		}
	}

	return maze
}
