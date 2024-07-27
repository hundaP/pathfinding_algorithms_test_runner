export function astar(grid, startNode, endNode) {
    class PriorityQueue {
        constructor(comparator = (a, b) => a.f - b.f) {
            this.heap = [null];
            this.comparator = comparator;
            this.nodeIndices = new Map();
        }

        enqueue(node) {
            this.heap.push(node);
            this.nodeIndices.set(node, this.heap.length - 1);
            this.bubbleUp();
        }

        dequeue() {
            const min = this.heap[1];
            this.nodeIndices.delete(min);
            this.heap[1] = this.heap[this.heap.length - 1];
            this.nodeIndices.set(this.heap[1], 1);
            this.heap.pop();
            this.sinkDown();
            return min;
        }

        bubbleUp() {
            let index = this.heap.length - 1;
            while (index > 1 && this.comparator(this.heap[Math.floor(index / 2)], this.heap[index]) > 0) {
                this.nodeIndices.set(this.heap[Math.floor(index / 2)], index);
                this.nodeIndices.set(this.heap[index], Math.floor(index / 2));
                [this.heap[Math.floor(index / 2)], this.heap[index]] = [this.heap[index], this.heap[Math.floor(index / 2)]];
                index = Math.floor(index / 2);
            }
        }

        sinkDown() {
            let index = 1;
            while (index * 2 < this.heap.length) {
                let smallerIndex = index * 2;
                if (index * 2 + 1 < this.heap.length && this.comparator(this.heap[index * 2 + 1], this.heap[index * 2]) < 0) {
                    smallerIndex = index * 2 + 1;
                }
                if (this.comparator(this.heap[smallerIndex], this.heap[index]) < 0) {
                    this.nodeIndices.set(this.heap[smallerIndex], index);
                    this.nodeIndices.set(this.heap[index], smallerIndex);
                    [this.heap[smallerIndex], this.heap[index]] = [this.heap[index], this.heap[smallerIndex]];
                    index = smallerIndex;
                } else {
                    break;
                }
            }
        }

        decreaseKey(node, newDistance) {
            let index = this.nodeIndices.get(node);
            if (index === undefined) return;
            this.heap[index].distance = newDistance;
            this.bubbleUp();
        }

        contains(node) {
            return this.nodeIndices.has(node);
        }

        isEmpty() {
            return this.heap.length === 1;
        }
    }

    const openList = new PriorityQueue((a, b) => a.f - b.f);
    const closedList = new Set();
    const visitedNodesInOrder = [];
    startNode.distance = 0;
    startNode.h = heuristic(startNode, endNode);
    startNode.f = startNode.h;
    openList.enqueue(startNode);
    while (!openList.isEmpty()) {
        const currentNode = openList.dequeue();
        closedList.add(currentNode);
        visitedNodesInOrder.push(currentNode);
        if (currentNode === endNode) {
            return visitedNodesInOrder;
        }
        const neighbors = getUnvisitedNeighbors(currentNode, grid);
        for (const neighbor of neighbors) {
            if (closedList.has(neighbor) || neighbor.isWall) continue;
            const gScore = currentNode.distance + 1;
            const hScore = neighbor.h || heuristic(neighbor, endNode);
            neighbor.h = hScore;
            if (!openList.contains(neighbor)) {
                openList.enqueue(neighbor);
                neighbor.distance = gScore;
                neighbor.f = gScore + hScore;
                neighbor.previousNode = currentNode;
            } else if (gScore < neighbor.distance) {
                neighbor.distance = gScore;
                neighbor.f = gScore + hScore;
                neighbor.previousNode = currentNode;
                openList.decreaseKey(neighbor, neighbor.f);
            }
        }
    }
    return visitedNodesInOrder;
}

function heuristic(node, endNode) {
    return Math.abs(node.row - endNode.row) + Math.abs(node.col - endNode.col);
}

function getUnvisitedNeighbors(node, grid) {
    const neighbors = [];
    const { col, row } = node;
    if (row > 0) neighbors.push(grid[row - 1][col]);
    if (row < grid.length - 1) neighbors.push(grid[row + 1][col]);
    if (col > 0) neighbors.push(grid[row][col - 1]);
    if (col < grid[0].length - 1) neighbors.push(grid[row][col + 1]);
    return neighbors.filter(neighbor => !neighbor.isVisited);
}
