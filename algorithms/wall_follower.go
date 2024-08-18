package algorithms

import "pathfinding_algorithms_test_runner/maze"

// WallFollowerAlgorithm performs the wall follower algorithm on the grid
func WallFollowerAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	visitedNodesInOrder := []maze.Node{}
	startNode.Distance = 0
	currentNode := startNode
	var previousNode *maze.Node

	for currentNode != endNode {
		currentNode.IsVisited = true
		currentNode.NoOfVisits++
		visitedNodesInOrder = append(visitedNodesInOrder, *currentNode)

		neighbors := getUnvisitedNeighbors(currentNode, grid)
		var nextNode *maze.Node

		for _, neighbor := range neighbors {
			if neighbor != previousNode && !neighbor.IsWall {
				nextNode = neighbor
				break
			}
		}

		if nextNode != nil {
			nextNode.Distance = currentNode.Distance + 1
			nextNode.PreviousNode = currentNode
			nextNode.PreviousID = currentNode.ID
			previousNode = currentNode
			currentNode = nextNode
		} else {
			if previousNode != nil {
				currentNode = previousNode
				previousNode = currentNode.PreviousNode
			} else {
				if currentNode.PreviousID != 0 {
					currentNode = findNodeByID(grid, currentNode.PreviousID)
					previousNode = findNodeByID(grid, currentNode.PreviousID).PreviousNode
				} else {
					break
				}
			}
		}
	}

	return visitedNodesInOrder
}

// Helper function to find a node by its ID
func findNodeByID(grid [][]maze.Node, id uint32) *maze.Node {
	for _, row := range grid {
		for _, node := range row {
			if node.ID == id {
				return &node
			}
		}
	}
	return nil
}
