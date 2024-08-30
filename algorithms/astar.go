package algorithms

import (
	"container/heap"
	"math"

	"pathfinding_algorithms_test_runner/maze"
)

func heuristic(node, endNode *maze.Node) float32 {
	//manhattanDistance := float32(math.Abs(float64(node.X-endNode.X)) + math.Abs(float64(node.Y-endNode.Y)))
	//euclideanDistance := float32(math.Sqrt(float64(node.X-endNode.X)*float64(node.X-endNode.X) + float64(node.Y-endNode.Y)*float64(node.Y-endNode.Y)))
	//chebyshevDistance := float32(math.Max(float64(node.X-endNode.X), float64(node.Y-endNode.Y)))
	canberraDistance := float32(math.Abs(float64(node.X-endNode.X))/(float64(node.X)+float64(endNode.X)) + math.Abs(float64(node.Y-endNode.Y))/(float64(node.Y)+float64(endNode.Y)))

	return canberraDistance
}

func AstarAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	openList := &PriorityQueue{useAstar: true}
	heap.Init(openList)

	closedSet := make(map[*maze.Node]bool)
	inOpenSet := make(map[*maze.Node]bool)
	visitedNodesInOrder := []maze.Node{}

	startNode.Distance = 0
	startNode.G = 0
	startNode.F = heuristic(startNode, endNode)
	heap.Push(openList, startNode)
	inOpenSet[startNode] = true

	for openList.Len() > 0 {
		currentNode := heap.Pop(openList).(*maze.Node)
		delete(inOpenSet, currentNode)

		if currentNode == endNode {
			return visitedNodesInOrder
		}

		closedSet[currentNode] = true
		visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

		neighbors := getUnvisitedNeighbors(currentNode, grid)
		for _, neighbor := range neighbors {
			if closedSet[neighbor] || neighbor.IsWall {
				continue
			}

			gScore := currentNode.G + 1
			hScore := heuristic(neighbor, endNode)

			if !inOpenSet[neighbor] {
				neighbor.Distance = uint32(gScore)
				neighbor.G = gScore
				neighbor.F = gScore + hScore
				neighbor.PreviousNode = currentNode
				neighbor.IsVisited = true
				heap.Push(openList, neighbor)
				inOpenSet[neighbor] = true
			} else if gScore < neighbor.G {
				neighbor.Distance = uint32(gScore)
				neighbor.G = gScore
				neighbor.F = gScore + hScore
				neighbor.PreviousNode = currentNode
				heap.Fix(openList, openList.IndexOf(neighbor))
			}
		}
	}

	return visitedNodesInOrder
}

