package main

import (
	"fmt"
	"pathfinding_algorithms_test_runner/algorithms"
	"pathfinding_algorithms_test_runner/maze"
)

func displayMazeWithPath(grid [][]maze.Node, path []maze.Node, startNode, endNode maze.Node) {
	for y := 0; y < len(grid); y++ {
		row := ""
		for x := 0; x < len(grid[y]); x++ {
			cell := grid[y][x]
			if cell.X == startNode.X && cell.Y == startNode.Y {
				row += "S"
			} else if cell.X == endNode.X && cell.Y == endNode.Y {
				row += "E"
			} else if contains(path, cell) {
				row += "P"
			} else {
				if cell.IsWall {
					row += "#"
				} else {
					row += " "
				}
			}
		}
		fmt.Println(row)
	}
}

func contains(path []maze.Node, cell maze.Node) bool {
	for _, node := range path {
		if node.X == cell.X && node.Y == cell.Y {
			return true
		}
	}
	return false
}

func runDisplayMaze() {
	// Generate the maze
	mazeData := maze.GenerateMaze(15, 15, false)

	// Extract the grid from the maze data
	grid := mazeData["gridDijkstra"].([][]maze.Node)
	startNode := mazeData["gridDijkstraStartNode"].(maze.Node)
	endNode := mazeData["gridDijkstraEndNode"].(maze.Node)

	// Define the algorithms
	algorithms := map[string]algorithms.Algorithm{
		"Dijkstra":     algorithms.Dijkstra{},
		"AStar":        algorithms.Astar{},
		"BFS":          algorithms.BFS{},
		"DFS":          algorithms.DFS{},
		"WallFollower": algorithms.WallFollower{},
	}

	// Display the maze with paths found by each algorithm
	for name, algorithm := range algorithms {
		fmt.Printf("Path found using %s:\n", name)
		visitedNodesInOrder := algorithm.FindPath(grid, &startNode, &endNode)
		path := nodesInShortestPathOrder(visitedNodesInOrder, &endNode)
		fmt.Printf("Path length: %d\n", len(path))
		displayMazeWithPath(grid, path, startNode, endNode)
		fmt.Println()
	}
}

func nodesInShortestPathOrder(visitedNodesInOrder []maze.Node, endNode *maze.Node) []maze.Node {
	var path []maze.Node
	currentNode := endNode
	for currentNode != nil {
		path = append([]maze.Node{*currentNode}, path...)
		currentNode = currentNode.PreviousNode
	}
	return path
}

func main() {
	runDisplayMaze()
}
