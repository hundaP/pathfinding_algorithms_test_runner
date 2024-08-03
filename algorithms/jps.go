package algorithms

/*import (
	"container/heap"
	"math"
	"pathfinding_algorithms_test_runner/maze"
)

type PriorityQueue []*maze.Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*maze.Node)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

func heuristic(a, b *maze.Node) uint32 {
	return uint32(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

const MaxRecursionDepth = 100

func AstarAlgorithm(grid [][]maze.Node, startNode, endNode *maze.Node) []maze.Node {
	openSet := &PriorityQueue{}
	heap.Init(openSet)
	heap.Push(openSet, startNode)

	cameFrom := make(map[*maze.Node]*maze.Node)
	gScore := make(map[*maze.Node]uint32)
	gScore[startNode] = 0

	fScore := make(map[*maze.Node]uint32)
	fScore[startNode] = heuristic(startNode, endNode)

	visitedNodesInOrder := []maze.Node{}
	closedSet := make(map[*maze.Node]bool)

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*maze.Node)
		visitedNodesInOrder = append(visitedNodesInOrder, *current)

		if current.X == endNode.X && current.Y == endNode.Y {
			endNode.PreviousNode = current.PreviousNode // Ensure endNode has its PreviousNode set
			return visitedNodesInOrder
		}

		closedSet[current] = true

		for _, neighbor := range identifySuccessors(grid, current, endNode) {
			if neighbor.IsWall || closedSet[neighbor] {
				continue
			}

			tentativeGScore := gScore[current] + uint32(math.Sqrt(float64((neighbor.X-current.X)*(neighbor.X-current.X)+(neighbor.Y-current.Y)*(neighbor.Y-current.Y))))

			if g, ok := gScore[neighbor]; !ok || tentativeGScore < g {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + heuristic(neighbor, endNode)
				neighbor.PreviousNode = current // Set the previous node

				// Only push the neighbor to the open set if it's not already in it
				if !ok {
					heap.Push(openSet, neighbor)
				} else {
					// Update the position of the neighbor in the priority queue
					heap.Fix(openSet, findIndex(openSet, neighbor))
				}
			}
		}
	}

	return visitedNodesInOrder
}

func identifySuccessors(grid [][]maze.Node, node, endNode *maze.Node) []*maze.Node {
	var successors []*maze.Node
	for _, direction := range getNeighbors(grid, node) {
		jumpNode := jump(grid, node, direction, endNode, 0)
		if jumpNode != nil {
			successors = append(successors, jumpNode)
		}
	}
	return successors
}

func jump(grid [][]maze.Node, current, direction, endNode *maze.Node, depth int) *maze.Node {
	if depth > MaxRecursionDepth {
		return nil
	}

	x, y := int(current.X)+int(direction.X), int(current.Y)+int(direction.Y)
	if x < 0 || y < 0 || x >= len(grid[0]) || y >= len(grid) || grid[y][x].IsWall {
		return nil
	}
	jumpNode := &grid[y][x]
	if jumpNode == endNode {
		return jumpNode
	}

	// Check for forced neighbors
	if direction.X != 0 {
		// Horizontal movement
		if (y > 0 && grid[y-1][x].IsWall && !grid[y-1][int(current.X)].IsWall) ||
			(y < len(grid)-1 && grid[y+1][x].IsWall && !grid[y+1][int(current.X)].IsWall) {
			return jumpNode
		}
	} else if direction.Y != 0 {
		// Vertical movement
		if (x > 0 && grid[y][x-1].IsWall && !grid[int(current.Y)][x-1].IsWall) ||
			(x < len(grid[0])-1 && grid[y][x+1].IsWall && !grid[int(current.Y)][x+1].IsWall) {
			return jumpNode
		}
	}

	return jump(grid, jumpNode, direction, endNode, depth+1)
}

func getNeighbors(grid [][]maze.Node, node *maze.Node) []*maze.Node {
	var neighbors []*maze.Node

	directions := []struct{ dx, dy int }{
		{dx: -1, dy: 0}, {dx: 1, dy: 0},
		{dx: 0, dy: -1}, {dx: 0, dy: 1},
	}

	for _, dir := range directions {
		x, y := int(node.X)+dir.dx, int(node.Y)+dir.dy
		if x >= 0 && y >= 0 && x < len(grid[0]) && y < len(grid) {
			neighbors = append(neighbors, &grid[y][x])
		}
	}

	return neighbors
}

func findIndex(pq *PriorityQueue, node *maze.Node) int {
	for i, n := range *pq {
		if n == node {
			return i
		}
	}
	return -1
}
*/
