export function wallFollower(grid, startNode, endNode) {
    const visitedNodesInOrder = [];
    startNode.distance = 0;
    let currentNode = startNode;
    let previousNode = null;
    while (currentNode !== endNode) {
        currentNode.isVisited = true;
        currentNode.noOfVisits = (currentNode.noOfVisits || 0) + 1; // increment visit count
        visitedNodesInOrder.push(currentNode);
        const neighbors = getUnvisitedNeighbors(currentNode, grid);
        let nextNode = null;
        for (const neighbor of neighbors) {
            if (neighbor !== previousNode && !neighbor.isWall) {
                nextNode = neighbor;
                break;
            }
        }
        if (nextNode) {
            nextNode.distance = currentNode.distance + 1;
            nextNode.previousNode = currentNode;
            previousNode = currentNode;
            currentNode = nextNode;
        } else {
            if (previousNode) {
                currentNode = previousNode;
                previousNode = currentNode.previousNode;
            } else {
                break;
            }
        }
    }
    return visitedNodesInOrder;
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
