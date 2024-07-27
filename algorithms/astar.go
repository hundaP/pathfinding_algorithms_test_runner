package algorithms

import (
	"math"
	"pathfinding_algorithms_test_runner/maze"
)

func heuristic(node, endNode *maze.Node) float64 {
	return math.Abs(float64(node.X-endNode.X)) + math.Abs(float64(node.Y-endNode.Y))
}

func AstarAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	openList := NewPriorityQueue()
	closedList := make(map[*maze.Node]bool)
	visitedNodesInOrder := []maze.Node{}

	startNode.Distance = 0
	startNode.H = heuristic(startNode, endNode)
	startNode.F = startNode.H
	openList.Enqueue(startNode)

	for !openList.IsEmpty() {
		currentNode := openList.Dequeue()
		closedList[currentNode] = true
		visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

		if currentNode == endNode {
			return visitedNodesInOrder
		}

		neighbors := getUnvisitedNeighbors(currentNode, grid)
		for _, neighbor := range neighbors {
			if closedList[neighbor] || neighbor.IsWall {
				continue
			}

			gScore := currentNode.Distance + 1
			hScore := neighbor.H
			if hScore == 0 {
				hScore = heuristic(neighbor, endNode)
				neighbor.H = hScore
			}

			if !openList.Contains(neighbor) {
				neighbor.Distance = gScore
				neighbor.F = gScore + hScore
				neighbor.PreviousNode = currentNode
				openList.Enqueue(neighbor)
			} else if gScore < neighbor.Distance {
				neighbor.Distance = gScore
				neighbor.F = gScore + hScore
				neighbor.PreviousNode = currentNode
				openList.Update(neighbor, neighbor.F)
			}
		}
	}

	return visitedNodesInOrder
}
