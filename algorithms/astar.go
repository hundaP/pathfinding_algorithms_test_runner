package algorithms

import (
	"container/heap"
	"math"

	"pathfinding_algorithms_test_runner/maze"
)

func heuristic(node, endNode *maze.Node) float32 {
	manhattanDistance := float32(math.Abs(float64(node.X-endNode.X)) + math.Abs(float64(node.Y-endNode.Y)))
	//euclideanDistance := float32(math.Sqrt(float64(node.X-endNode.X)*float64(node.X-endNode.X) + float64(node.Y-endNode.Y)*float64(node.Y-endNode.Y)))
	//chebyshevDistance := float32(math.Max(float64(node.X-endNode.X), float64(node.Y-endNode.Y)))
	return manhattanDistance
}

func AstarAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	rows, cols := uint16(len(grid)), uint16(len(grid[0]))
	totalNodes := uint32(rows) * uint32(cols)

	openSet := &PriorityQueue{useAstar: true}
	heap.Init(openSet)

	gScore := make([]float32, totalNodes)
	fScore := make([]float32, totalNodes)
	for i := range gScore {
		gScore[i] = math.MaxFloat32
		fScore[i] = math.MaxFloat32
	}

	startIndex := uint32(startNode.Y)*uint32(cols) + uint32(startNode.X)
	gScore[startIndex] = 0
	fScore[startIndex] = heuristic(startNode, endNode)

	startNode.G = 0
	startNode.F = fScore[startIndex]
	heap.Push(openSet, startNode)

	visitedNodesInOrder := []maze.Node{}
	closedSet := make([]bool, totalNodes)

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*maze.Node)
		currentIndex := uint32(current.Y)*uint32(cols) + uint32(current.X)

		visitedNodesInOrder = append(visitedNodesInOrder, *current)

		if current == endNode {
			return visitedNodesInOrder
		}

		closedSet[currentIndex] = true

		for _, neighbor := range getUnvisitedNeighbors(current, grid) {
			if neighbor.Y >= rows || neighbor.X >= cols {
				continue
			}

			neighborIndex := uint32(neighbor.Y)*uint32(cols) + uint32(neighbor.X)
			if neighbor.IsWall || closedSet[neighborIndex] {
				continue
			}

			tentativeGScore := gScore[currentIndex] + 1

			if tentativeGScore < gScore[neighborIndex] {
				neighbor.PreviousNode = current
				gScore[neighborIndex] = tentativeGScore
				fScore[neighborIndex] = gScore[neighborIndex] + heuristic(neighbor, endNode)
				neighbor.G = tentativeGScore
				neighbor.F = fScore[neighborIndex]

				if !contains(openSet, neighbor) {
					heap.Push(openSet, neighbor)
				} else {
					heap.Fix(openSet, openSet.IndexOf(neighbor))
				}
			}
		}
	}

	return visitedNodesInOrder
}
