export function astar(grid, startNode, endNode) {
    const openList = [startNode];
    const closedList = [];
    const visitedNodesInOrder = [];
    startNode.distance = 0;
    startNode.h = heuristic(startNode, endNode);
    startNode.f = startNode.h;
    while (openList.length > 0) {
        const currentNode = getLowestFScore(openList);
        openList.splice(openList.indexOf(currentNode), 1);
        closedList.push(currentNode);
        visitedNodesInOrder.push(currentNode);
        if (currentNode === endNode) {
            return visitedNodesInOrder;
        }
        const neighbors = getUnvisitedNeighbors(currentNode, grid);
        for (const neighbor of neighbors) {
            if (closedList.includes(neighbor) || neighbor.isWall) continue;
            const gScore = currentNode.distance + 1;
            const hScore = heuristic(neighbor, endNode);
            if (!openList.includes(neighbor)) {
                openList.push(neighbor);
                neighbor.distance = gScore;
                neighbor.h = hScore;
                neighbor.f = gScore + hScore;
                neighbor.previousNode = currentNode;
            } else if (gScore < neighbor.distance) {
                neighbor.distance = gScore;
                neighbor.f = gScore + hScore;
                neighbor.previousNode = currentNode;
            }
        }
    }
    return visitedNodesInOrder;
}

// f = g + h
function getLowestFScore(nodes){
    return nodes.reduce((lowest, node) => node.f < lowest.f ? node : lowest, nodes[0]);
}

function heuristic(node, endNode){
    return Math.abs(node.row - endNode.row) + Math.abs(node.col - endNode.col);
}

function getUnvisitedNeighbors(node, grid){
    const neighbors = [];
    const {col, row} = node;
    if(row > 0) neighbors.push(grid[row - 1][col]);
    if(row < grid.length - 1) neighbors.push(grid[row + 1][col]);
    if(col > 0) neighbors.push(grid[row][col - 1]);
    if(col < grid[0].length - 1) neighbors.push(grid[row][col + 1]);
    return neighbors.filter(neighbor => !neighbor.isVisited);
}