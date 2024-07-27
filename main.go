package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"pathfinding_algorithms_test_runner/algorithms"
	"pathfinding_algorithms_test_runner/maze"
)

type Metrics struct {
	Time              []float64
	VisitedNodes      []int
	VisitedPercentage []float64
	PathLength        []int
	MemoryUsed        []float64
}

var algorithmsMap = map[string]algorithms.Algorithm{
	"dijkstra":     algorithms.Dijkstra{},
	"astar":        algorithms.Astar{},
	"bfs":          algorithms.BFS{},
	"dfs":          algorithms.DFS{},
	"wallFollower": algorithms.WallFollower{},
}

func main() {
	if len(os.Args) < 3 {
		runTestsWithIncreasingSize()
	} else {
		mazeSize, _ := strconv.Atoi(os.Args[1])
		numTests, _ := strconv.Atoi(os.Args[2])
		runTest(mazeSize, numTests)
	}
}

func runTestsWithIncreasingSize() {
	size := 725
	for {
		fmt.Printf("Running tests with maze size %d\n", size)
		err := runTest(size, 10)
		if err != nil {
			fmt.Printf("Test failed for maze size %d: %s\n", size, err.Error())
			break
		}
		size += 25
	}
}

func runTest(mazeSize, numTests int) error {
	numRows := mazeSize
	numCols := mazeSize
	metricsSPOn := initializeMetrics()
	metricsSPOff := initializeMetrics()

	// Test mazes with a single path
	for i := 0; i < numTests; i++ {
		grids, startNodes, endNodes := getInitialGrid(numRows, numCols, true)
		var wg sync.WaitGroup
		for algorithm := range algorithmsMap {
			wg.Add(1)
			go func(algorithm string) {
				defer wg.Done()
				runAlgorithm(algorithm, grids[algorithm], startNodes[algorithm], endNodes[algorithm], metricsSPOn)
			}(algorithm)
		}
		wg.Wait()
		fmt.Printf("Completed test %d of %d for mazes with a single path, for size: %d\n", i+1, numTests, mazeSize)
	}

	// Test mazes with multiple paths
	for i := 0; i < numTests; i++ {
		grids, startNodes, endNodes := getInitialGrid(numRows, numCols, false)
		var wg sync.WaitGroup
		for algorithm := range algorithmsMap {
			wg.Add(1)
			go func(algorithm string) {
				defer wg.Done()
				runAlgorithm(algorithm, grids[algorithm], startNodes[algorithm], endNodes[algorithm], metricsSPOff)
			}(algorithm)
		}
		wg.Wait()
		fmt.Printf("Completed test %d of %d for mazes with multiple paths, for size: %d\n", i+1, numTests, mazeSize)
		runtime.GC()
	}

	averagesSPOn := calculateAverages(metricsSPOn)
	averagesSPOff := calculateAverages(metricsSPOff)
	writeResultsToCsv(fmt.Sprintf("./data/averages%dx%dx%d.csv", numRows, numCols, numTests), averagesSPOn, averagesSPOff)

	return nil
}

func initializeMetrics() map[string]*Metrics {
	return map[string]*Metrics{
		"dijkstra":     {},
		"astar":        {},
		"bfs":          {},
		"dfs":          {},
		"wallFollower": {},
	}
}

func runAlgorithm(algorithm string, grid [][]maze.Node, startNode maze.Node, endNode maze.Node, metrics map[string]*Metrics) {
	startTime := time.Now()
	initialMemoryUsage := runtime.MemStats{}
	runtime.ReadMemStats(&initialMemoryUsage)
	visitedNodesInOrder := algorithmsMap[algorithm].FindPath(grid, &startNode, &endNode)

	finalMemoryUsage := runtime.MemStats{}
	runtime.ReadMemStats(&finalMemoryUsage)
	endTime := time.Now()
	nodesInShortestPathOrder := getNodesInShortestPathOrder(&endNode)

	fmt.Println(len(nodesInShortestPathOrder))

	timeTaken := endTime.Sub(startTime).Milliseconds() // Convert to milliseconds

	memoryUsed := float64(finalMemoryUsage.HeapAlloc-initialMemoryUsage.HeapAlloc) / (1024 * 1024) // Convert to MB

	totalNodes := len(grid) * len(grid[0])
	wallNodes := countWallNodes(grid)
	nonWallNodes := totalNodes - wallNodes
	visitedPercentage := (float64(len(visitedNodesInOrder)) / float64(nonWallNodes)) * 100

	metrics[algorithm].Time = append(metrics[algorithm].Time, float64(timeTaken))
	metrics[algorithm].VisitedNodes = append(metrics[algorithm].VisitedNodes, len(visitedNodesInOrder))
	metrics[algorithm].VisitedPercentage = append(metrics[algorithm].VisitedPercentage, visitedPercentage)
	metrics[algorithm].PathLength = append(metrics[algorithm].PathLength, len(nodesInShortestPathOrder))
	metrics[algorithm].MemoryUsed = append(metrics[algorithm].MemoryUsed, memoryUsed)
}

func getNodesInShortestPathOrder(endNode *maze.Node) []*maze.Node {
	var nodesInShortestPathOrder []*maze.Node
	currentNode := endNode
	for currentNode != nil {
		nodesInShortestPathOrder = append([]*maze.Node{currentNode}, nodesInShortestPathOrder...)
		currentNode = currentNode.PreviousNode
	}
	return nodesInShortestPathOrder
}

func countWallNodes(grid [][]maze.Node) int {
	count := 0
	for _, row := range grid {
		for _, node := range row {
			if node.IsWall {
				count++
			}
		}
	}
	return count
}

func calculateAverages(metrics map[string]*Metrics) map[string]map[string]float64 {
	averages := make(map[string]map[string]float64)
	for algorithm, metric := range metrics {
		averages[algorithm] = make(map[string]float64)
		numTests := len(metric.Time)
		for _, time := range metric.Time {
			averages[algorithm]["time"] += time
		}
		for _, visitedNodes := range metric.VisitedNodes {
			averages[algorithm]["visitedNodes"] += float64(visitedNodes)
		}
		for _, visitedPercentage := range metric.VisitedPercentage {
			averages[algorithm]["visitedPercentage"] += visitedPercentage
		}
		for _, pathLength := range metric.PathLength {
			averages[algorithm]["pathLength"] += float64(pathLength)
		}
		for _, memoryUsed := range metric.MemoryUsed {
			averages[algorithm]["memoryUsed"] += memoryUsed
		}
		for key := range averages[algorithm] {
			averages[algorithm][key] /= float64(numTests)
		}
	}
	return averages
}

func writeResultsToCsv(filename string, averagesSPOn, averagesSPOff map[string]map[string]float64) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Algorithm", "SinglePath", "Time", "VisitedNodes", "VisitedPercentage", "PathLength", "MemoryUsed"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %s", err)
	}

	for algorithm, metrics := range averagesSPOn {
		row := []string{
			algorithm,
			"true",
			fmt.Sprintf("%.4f", metrics["time"]),
			fmt.Sprintf("%.0f", metrics["visitedNodes"]),
			fmt.Sprintf("%.2f", metrics["visitedPercentage"]),
			fmt.Sprintf("%.0f", metrics["pathLength"]),
			fmt.Sprintf("%.2f", metrics["memoryUsed"]),
		}
		if err := writer.Write(row); err != nil {
			log.Fatalf("Failed to write row for %s: %s", algorithm, err)
		}
	}

	for algorithm, metrics := range averagesSPOff {
		row := []string{
			algorithm,
			"false",
			fmt.Sprintf("%.4f", metrics["time"]),
			fmt.Sprintf("%.0f", metrics["visitedNodes"]),
			fmt.Sprintf("%.2f", metrics["visitedPercentage"]),
			fmt.Sprintf("%.0f", metrics["pathLength"]),
			fmt.Sprintf("%.2f", metrics["memoryUsed"]),
		}
		if err := writer.Write(row); err != nil {
			log.Fatalf("Failed to write row for %s: %s", algorithm, err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatalf("Error flushing writer: %s", err)
	}
}

func getInitialGrid(numRows, numCols int, singlePath bool) (map[string][][]maze.Node, map[string]maze.Node, map[string]maze.Node) {
	grids := make(map[string][][]maze.Node)
	startNodes := make(map[string]maze.Node)
	endNodes := make(map[string]maze.Node)

	mazeData := maze.GenerateMaze(numRows, numCols, singlePath)

	grids["dijkstra"] = mazeData["gridDijkstra"].([][]maze.Node)
	grids["astar"] = mazeData["gridAstar"].([][]maze.Node)
	grids["bfs"] = mazeData["gridBFS"].([][]maze.Node)
	grids["dfs"] = mazeData["gridDFS"].([][]maze.Node)
	grids["wallFollower"] = mazeData["gridWallFollower"].([][]maze.Node)

	startNodes["dijkstra"] = mazeData["gridDijkstraStartNode"].(maze.Node)
	startNodes["astar"] = mazeData["gridAstarStartNode"].(maze.Node)
	startNodes["bfs"] = mazeData["gridBFSStartNode"].(maze.Node)
	startNodes["dfs"] = mazeData["gridDFSStartNode"].(maze.Node)
	startNodes["wallFollower"] = mazeData["gridWallFollowerStartNode"].(maze.Node)

	endNodes["dijkstra"] = mazeData["gridDijkstraEndNode"].(maze.Node)
	endNodes["astar"] = mazeData["gridAstarEndNode"].(maze.Node)
	endNodes["bfs"] = mazeData["gridBFSEndNode"].(maze.Node)
	endNodes["dfs"] = mazeData["gridDFSEndNode"].(maze.Node)
	endNodes["wallFollower"] = mazeData["gridWallFollowerEndNode"].(maze.Node)

	return grids, startNodes, endNodes
}
