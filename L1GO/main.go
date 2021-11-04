package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"
)

const BOARD_SIZE = 4

const shuffleshuffle = 30

var HEURISTIC int
var lastNode *Node

type Node struct {
	data      [BOARD_SIZE * BOARD_SIZE]int8
	heuristic int8
	moves     int8
	Index     int8
	prevNode  *Node
}

func node_new(data [BOARD_SIZE * BOARD_SIZE]int8, moves int8, prevNode *Node) *Node {
	node := Node{data: data, moves: moves, prevNode: prevNode}
	node.heuristic = node_heuristic(data)
	return &node
}

func node_heuristic(data [BOARD_SIZE * BOARD_SIZE]int8) int8 {
	switch HEURISTIC {
	case 0:
		return heuristic_manhattan(data)
	case 1:
		return heuristic_hammingPlus(data)
	}
	return 0
}

func heuristic_manhattan(data [BOARD_SIZE * BOARD_SIZE]int8) int8 {
	heur := int8(0)

	for i, n := range data {

		if n != 0 {
			A, B := board_get2Dfrom1D(int8(i))
			Acor, Bcor := board_get2Dfrom1D(n - 1)
			heur += Abs(A - Acor)
			heur += Abs(B - Bcor)
		}
	}

	mask := []int{-1, -BOARD_SIZE, +1, BOARD_SIZE, -BOARD_SIZE - 1, -BOARD_SIZE + 1, BOARD_SIZE + 1, BOARD_SIZE - 1}
	for i, n := range data {
		x, y := board_get2Dfrom1D(int8(i))
		guide := calcIfNeigh(x, y)

		for j, move := range guide {
			if move == 1 && data[i+mask[j]] != n+int8(mask[j]) {
				heur++
			}
		}
	}
	return heur
}

func heuristic_hammingPlus(data [BOARD_SIZE * BOARD_SIZE]int8) int8 {
	heur := int8(0)

	for i, n := range data {

		if n != 0 {
			if n-1 != int8(i) {
				heur += 1
			}
		}
	}

	mask := []int{-1, -BOARD_SIZE, +1, BOARD_SIZE, -BOARD_SIZE - 1, -BOARD_SIZE + 1, BOARD_SIZE + 1, BOARD_SIZE - 1}
	for i, n := range data {
		x, y := board_get2Dfrom1D(int8(i))
		guide := calcIfNeigh(x, y)

		for j, move := range guide {
			if move == 1 && data[i+mask[j]] != n+int8(mask[j]) {
				heur++
			}
		}
	}

	return heur
}

func calcIfNeigh(a int8, b int8) [8]int8 {
	var ret [8]int8
	for i := range ret {
		ret[i] = 0
	}

	if a > 0 {
		ret[0] = 1
	}

	if a < BOARD_SIZE-1 {
		ret[2] = 1
	}

	if b > 0 {
		ret[1] = 1
	}

	if b < BOARD_SIZE-1 {
		ret[3] = 1
	}

	if ret[0] == 1 && ret[1] == 1 {
		ret[4] = 1
	}

	if ret[2] == 1 && ret[1] == 1 {
		ret[5] = 1
	}

	if ret[2] == 1 && ret[3] == 1 {
		ret[6] = 1
	}

	if ret[0] == 1 && ret[3] == 1 {
		ret[7] = 1
	}

	return ret
}

func board_get1Dfrom2D(i int8, j int8) int8 {
	return (i*BOARD_SIZE + j)
}

func board_get2Dfrom1D(i int8) (int8, int8) {
	return (i % BOARD_SIZE), (i / BOARD_SIZE)
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {

	return pq[i].heuristic+pq[i].moves < pq[j].heuristic+pq[j].moves
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.Index = -1
	*pq = old[0 : n-1]
	return node
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := int8(len(*pq))
	item := x.(*Node)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq PriorityQueue) Swap(i int, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = int8(i)
	pq[j].Index = int8(j)
}

func Solver(display bool, board [BOARD_SIZE * BOARD_SIZE]int8) {
	if !isSolvable(board) {
		fmt.Println("BAD")
		return
	}
	added := 0
	visited := make([]*Node, 0)
	pq := &PriorityQueue{}
	heap.Push(pq, node_new(board, 0, nil))

	for {
		removed := heap.Pop(pq).(*Node)

		if !searchIfExists(visited, removed) {
			visited = append(visited, removed)
		}

		if removed == nil {
			break
		}

		if isSolved(removed.data) {
			lastNode = removed
			break
		}

		outcomes := getOutcomes(removed.data)

		for _, n := range outcomes {
			if removed.prevNode != nil && removed.prevNode.data == n {
				continue
			}
			new_node := node_new(n, removed.moves+1, removed)

			heap.Push(pq, new_node)
			added++
		}
	}
	if display {
		displaySolution()
	}
	fmt.Println(len(visited), " ", getStepCount(), " ", added)
	//fmt.Println("solved")
}

func searchIfExists(visited []*Node, node *Node) bool {
	for _, n := range visited {
		if node.data == n.data {
			return true
		}
	}
	return false
}
func getInvCount(board [BOARD_SIZE * BOARD_SIZE]int8) int {
	inv_count := 0

	for i, n := range board {
		for j := i + 1; j < BOARD_SIZE*BOARD_SIZE; j++ {
			if n > 0 && board[j] > 0 && board[j] < n {
				inv_count++
			}
		}
	}

	return inv_count
}

func isSolvable(board [BOARD_SIZE * BOARD_SIZE]int8) bool {
	invCount := getInvCount(board)

	if BOARD_SIZE%2 == 0 {
		return (invCount%2 == 0)
	}

	xPos := 0

	for i, n := range board {
		if n == 0 {
			xPos = i
			break
		}
	}

	if xPos%2 == 0 {
		return (invCount%2 == 0)
	}

	return (invCount%2 == 1)
}

func displaySolution() {
	node := lastNode
	for i := 0; node != nil; i++ {
		fmt.Println(i)
		printBoard(node.data)
		node = node.prevNode
	}
}

func getStepCount() int {
	node := lastNode
	i := 0
	for ; node != nil; i++ {
		node = node.prevNode
	}
	return i
}

func printBoard(board [BOARD_SIZE * BOARD_SIZE]int8) {
	for i := 0; i < BOARD_SIZE; i++ {
		for j := 0; j < BOARD_SIZE; j++ {
			if board[board_get1Dfrom2D(int8(i), int8(j))] == 0 {
				fmt.Print(" \t")
				continue
			}
			fmt.Print(" \t", board[board_get1Dfrom2D(int8(i), int8(j))], " ")
		}
		fmt.Println()
	}
	fmt.Println()
}
func isSolved(board [BOARD_SIZE * BOARD_SIZE]int8) bool {

	for i, n := range board {
		if n != 0 && int8(i+1) != n {
			return false
		}
	}
	return board[len(board)-1] == 0
}

func Abs(a int8) int8 {
	if a < 0 {
		return -a
	}
	return a
}

func getOutcomes(board [BOARD_SIZE * BOARD_SIZE]int8) [][BOARD_SIZE * BOARD_SIZE]int8 {
	ret := make([][BOARD_SIZE * BOARD_SIZE]int8, 0)

	zeroCords := []int8{0, 0}
	zeroPos := 0

	for i, n := range board {
		if n == 0 {
			zeroCords[0], zeroCords[1] = board_get2Dfrom1D(int8(i))
			zeroPos = i
			break
		}
	}

	swapPos := 0

	if zeroCords[1] > 0 {
		swapPos = zeroPos - BOARD_SIZE
		board_clone := board_clone(board)
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[swapPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		ret = append(ret, board_clone)
	}
	if zeroCords[1] < BOARD_SIZE-1 {
		swapPos = zeroPos + BOARD_SIZE
		board_clone := board_clone(board)
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[swapPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		ret = append(ret, board_clone)
	}
	if zeroCords[0] > 0 {
		swapPos = zeroPos - 1
		board_clone := board_clone(board)
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[swapPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		ret = append(ret, board_clone)
	}
	if zeroCords[0] < BOARD_SIZE-1 {
		swapPos = zeroPos + 1
		board_clone := board_clone(board)
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[swapPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		board_clone[zeroPos] = board_clone[zeroPos] ^ board_clone[swapPos]
		ret = append(ret, board_clone)
	}

	return ret
}

func board_clone(board [BOARD_SIZE * BOARD_SIZE]int8) [BOARD_SIZE * BOARD_SIZE]int8 {
	var ret [BOARD_SIZE * BOARD_SIZE]int8

	for i, n := range board {
		ret[i] = n
	}

	return ret
}

func board_create() [BOARD_SIZE * BOARD_SIZE]int8 {
	var board [BOARD_SIZE * BOARD_SIZE]int8

	for i := range board {
		board[i] = int8(i + 1)
	}
	board[len(board)-1] = 0
	board = shuffle(board)

	printBoard(board)
	return board
}

func shuffle(board [BOARD_SIZE * BOARD_SIZE]int8) [BOARD_SIZE * BOARD_SIZE]int8 {
	possible := [4]int{-1, -BOARD_SIZE, 1, BOARD_SIZE}
	allowed := [4]int{1, 1, 1, 1}
	previous := -1
	var next int
	rand.Seed(time.Now().UnixNano())
	zeroCords := []int8{0, 0}
	zeroPos := 0

	for j, n := range board {
		if n == 0 {
			zeroCords[0], zeroCords[1] = board_get2Dfrom1D(int8(j))
			zeroPos = j
			break
		}
	}

	for i := 0; i < shuffleshuffle; i++ {
		if zeroCords[0] == 0 {
			allowed[0] = 0
		}

		if zeroCords[0] == BOARD_SIZE-1 {
			allowed[2] = 0
		}

		if zeroCords[1] == 0 {
			allowed[1] = 0
		}

		if zeroCords[1] == BOARD_SIZE-1 {
			allowed[3] = 0
		}
		next = rand.Intn(4)
		for allowed[next] == 0 || next == previous {
			next = rand.Intn(4)
		}

		board[zeroPos] = board[zeroPos] ^ board[zeroPos+possible[next]]
		board[zeroPos+possible[next]] = board[zeroPos] ^ board[zeroPos+possible[next]]
		board[zeroPos] = board[zeroPos] ^ board[zeroPos+possible[next]]
		previous = (next + 2) % 4
		zeroPos = zeroPos + possible[next]
		zeroCords[0], zeroCords[1] = board_get2Dfrom1D(int8(zeroPos))
		allowed = [4]int{1, 1, 1, 1}
		// fmt.Println("debug")
		// printBoard(board)
	}

	board[zeroPos] = board[zeroPos] ^ board[BOARD_SIZE*BOARD_SIZE-1]
	board[BOARD_SIZE*BOARD_SIZE-1] = board[zeroPos] ^ board[BOARD_SIZE*BOARD_SIZE-1]
	board[zeroPos] = board[zeroPos] ^ board[BOARD_SIZE*BOARD_SIZE-1]

	if !isSolvable(board) {
		board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE] = board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE] ^ board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE-1]
		board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE-1] = board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE] ^ board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE-1]
		board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE] = board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE] ^ board[BOARD_SIZE*BOARD_SIZE-BOARD_SIZE-1]
	}

	return board
}

func main() {
	board := board_create()
	//board := [BOARD_SIZE * BOARD_SIZE]int8{7, 13, 11, 5, 14, 1, 15, 12, 6, 9, 4, 8, 10, 3, 2}
	HEURISTIC = 0
	Solver(true, board)
	//HEURISTIC = 1
	//Solver(false, board)

}
