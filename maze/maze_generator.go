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

type Node struct {
	X, Y         int
	IsStart      bool
	IsEnd        bool
	Distance     float64
	IsVisited    bool
	IsWall       bool
	PreviousNode *Node
	GridId       int
	NoOfVisits   int
	H, F         float64
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

func createNode(x, y int, isWall bool, start, end *Cell, gridId int) Node {
	return Node{
		X:            x,
		Y:            y,
		IsStart:      start != nil && x == start.X && y == start.Y,
		IsEnd:        end != nil && x == end.X && y == end.Y,
		Distance:     math.Inf(1),
		IsVisited:    false,
		IsWall:       isWall,
		PreviousNode: nil,
		GridId:       gridId,
	}
}

func GenerateMaze(numRows, numCols int, singlePath bool) map[string]interface{} {
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

	grids := make(map[int][][]Node)
	for i := 1; i <= 5; i++ {
		grid := make([][]Node, len(maze.Grid))
		for y, row := range maze.Grid {
			grid[y] = make([]Node, len(row))
			for x, cell := range row {
				grid[y][x] = createNode(cell.X, cell.Y, cell.IsWall, maze.Start, maze.End, i)
			}
		}
		grids[i] = grid
	}

	return map[string]interface{}{
		"gridDijkstra":              grids[1],
		"gridAstar":                 grids[2],
		"gridBFS":                   grids[3],
		"gridDFS":                   grids[4],
		"gridWallFollower":          grids[5],
		"gridDijkstraStartNode":     grids[1][maze.Start.Y][maze.Start.X],
		"gridDijkstraEndNode":       grids[1][maze.End.Y][maze.End.X],
		"gridAstarStartNode":        grids[2][maze.Start.Y][maze.Start.X],
		"gridAstarEndNode":          grids[2][maze.End.Y][maze.End.X],
		"gridBFSStartNode":          grids[3][maze.Start.Y][maze.Start.X],
		"gridBFSEndNode":            grids[3][maze.End.Y][maze.End.X],
		"gridDFSStartNode":          grids[4][maze.Start.Y][maze.Start.X],
		"gridDFSEndNode":            grids[4][maze.End.Y][maze.End.X],
		"gridWallFollowerStartNode": grids[5][maze.Start.Y][maze.Start.X],
		"gridWallFollowerEndNode":   grids[5][maze.End.Y][maze.End.X],
	}
}
