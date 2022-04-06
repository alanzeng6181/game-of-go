package gamelogic

import (
	"errors"
)

type Board struct {
	Positions    []Stone
	Size         int16
	BlackCapture int16
	WhiteCapture int16
	roots        []int16
	groupSizes   []int16
	libertyCount []int16
	IsBlackNext  bool
}

func NewBoard(size int16) *Board {
	return &Board{
		Positions:    make([]Stone, size*size),
		Size:         size,
		roots:        make([]int16, size*size),
		groupSizes:   make([]int16, size*size),
		libertyCount: make([]int16, size*size),
		IsBlackNext:  true,
	}
}

func (board *Board) GetLiberty(row int16, col int16) int16 {
	return board.libertyCount[board.findRoot(row*board.Size+col)]
}

func (board *Board) neighbors(position int16) []int16 {
	neighbors := make([]int16, 0)

	if position%board.Size != 0 {
		neighbors = append(neighbors, position-1)
	}

	if position%board.Size != board.Size-1 {
		neighbors = append(neighbors, position+1)
	}

	if position/board.Size != 0 {
		neighbors = append(neighbors, position-board.Size)
	}

	if position/board.Size != board.Size-1 {
		neighbors = append(neighbors, position+board.Size)
	}
	return neighbors
}

func (board *Board) Move1(row int16, col int16) (bool, error) {
	return board.Move(row*board.Size + col)
}

func (board *Board) Move(position int16) (bool, error) {
	if board.Positions[position] > 0 {
		return false, errors.New("position is occupied already")
	}

	color := White
	if board.IsBlackNext {
		color = Black
	}

	neighbors := board.neighbors(position)
	zeroLiberty := true
	zeroLibertyOpponent := false
	allOpenents := true
	for _, neighbor := range neighbors {
		if board.Positions[neighbor] == color || board.Positions[neighbor] == Empty {
			allOpenents = false
		}
		if zeroLiberty && (board.Positions[neighbor] == Empty || (board.Positions[neighbor] == color && board.libertyCount[board.findRoot(neighbor)] != 1)) {
			zeroLiberty = false
		}
		if board.Positions[neighbor] != color &&
			board.Positions[neighbor] != Empty && board.libertyCount[board.findRoot(neighbor)] == 1 {
			zeroLibertyOpponent = true
		}
	}

	if allOpenents && !zeroLibertyOpponent {
		return false, errors.New("position is surrounded by opponent stones and would not cause any openent stones to be captured")
	}

	if zeroLiberty && !zeroLibertyOpponent {
		return false, errors.New("position would make liberty count of own stone group 0")
	}

	board.IsBlackNext = !board.IsBlackNext

	board.roots[position] = position
	board.groupSizes[position] = 1
	board.Positions[position] = color

	reducedOpponentRoots := make(map[int16]bool)
	captureList := make([]int16, 0)
	for _, neighbor := range neighbors {
		if board.Positions[position] == board.Positions[neighbor] {
			neighborRoot := board.findRoot(neighbor)
			board.merge(position, neighborRoot)
		} else if board.Positions[neighbor] == Empty {
			continue
		} else if _, exits := reducedOpponentRoots[neighbor]; !exits {
			opponentRoot := board.findRoot(neighbor)
			board.libertyCount[opponentRoot] -= 1
			if board.libertyCount[opponentRoot] == 0 {
				captureList = append(captureList, opponentRoot)
			}
		}
	}

	board.libertyCount[board.findRoot(position)] = board.calculateLiberty(position)
	for _, captureRoot := range captureList {
		board.capture(captureRoot)
	}

	return true, nil
}

func (board *Board) calculateLiberty(position int16) int16 {
	var liberty int16 = 0
	queue := make([]int16, 0)
	queue = append(queue, position)
	added := make(map[int16]bool)
	counted := make(map[int16]bool)
	added[position] = true
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, n := range board.neighbors(node) {
			if board.Positions[n] == Empty {
				if _, ok := counted[n]; !ok {
					liberty++
					counted[n] = true
				}
			} else if board.Positions[n] == board.Positions[node] {
				if _, ok := added[n]; !ok {
					added[n] = true
					queue = append(queue, n)
				}
			}
		}
	}
	return liberty
}

func (board *Board) capture(position int16) int16 {
	color := board.Positions[position]

	stack := make([]int16, 0)
	stack = append(stack, position)
	var captureCount int16 = 0

	processed := make(map[int16]bool)
	processed[position] = true

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		board.Positions[node] = Empty
		captureCount++

		neighbors := board.neighbors(position)
		for _, n := range neighbors {

			if board.Positions[n] == color {
				if _, done := processed[n]; done {
					continue
				}
				stack = append(stack, n)
				processed[n] = true
			} else if board.Positions[n] != Empty {
				nRoot := board.findRoot(n)
				board.libertyCount[nRoot] += 1
			}
		}
	}
	if color == Black {
		board.BlackCapture += captureCount
	} else {
		board.WhiteCapture += captureCount
	}
	return captureCount
}

func (s *Board) findRoot(position int16) int16 {
	if s.roots[position] < 0 {
		return -1
	}

	root := position
	for s.roots[root] != root {
		root = s.roots[root]
	}

	for s.roots[position] != position {
		position, s.roots[position] = s.roots[position], root
	}

	return root
}

func (s *Board) merge(positionA int16, positionB int16) (bool, int16) {
	rootA := s.findRoot(positionA)
	rootB := s.findRoot(positionB)

	if rootA == rootB {
		return false, rootA
	}

	if s.groupSizes[rootA] > s.groupSizes[rootB] {
		rootA, rootB = rootB, rootA
		positionA, positionB = positionB, positionA
	}

	for s.roots[positionA] != positionA {
		positionA, s.roots[positionA] = s.roots[positionA], rootB
	}
	s.roots[positionA] = rootB
	s.groupSizes[rootB] += s.groupSizes[rootA]
	return true, rootB
}

type Stone int8

const (
	Empty Stone = iota
	Black
	White
)

func (s Stone) OppositeColor() Stone {
	switch s {
	case Black:
		return White
	case White:
		return Black
	default:
		return s
	}
}
