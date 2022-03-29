package gamelogic

import (
	"errors"
)

type Board struct {
	Positions    []Stone
	Size         int16
	BlackCapture int16
	WhiteCapture int16
	Roots        []int16
	GroupSizes   []int16
	LibertyCount []int16
	IsBlackNext  bool
}

func NewBoard(size int16) *Board {
	return &Board{
		Positions:    make([]Stone, size*size),
		Size:         size,
		Roots:        make([]int16, size*size),
		GroupSizes:   make([]int16, size*size),
		LibertyCount: make([]int16, size*size),
		IsBlackNext:  true,
	}
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

func (board *Board) move(position int16) (bool, error) {
	if board.Positions[position] > 0 {
		return false, errors.New("position is occupied already")
	}

	board.Roots[position] = position
	if board.IsBlackNext {
		board.Positions[position] = Black
	} else {
		board.Positions[position] = White
	}

	neighbors := board.neighbors(position)
	zeroLiberty := true
	zeroLibertyOpponent := true
	allOpenents := true
	for _, neighbor := range neighbors {
		if board.Positions[neighbor] == board.Positions[position] {
			allOpenents = false
		}
		if zeroLiberty && board.Positions[neighbor] == board.Positions[position] && board.LibertyCount[board.findRoot(neighbor)] != 1 {
			zeroLiberty = false
		}
		if zeroLibertyOpponent && board.Positions[neighbor] != board.Positions[position] &&
			board.Positions[neighbor] != Empty && board.LibertyCount[board.findRoot(neighbor)] != 1 {
			zeroLibertyOpponent = false
		}
	}

	if allOpenents && !zeroLibertyOpponent {
		return false, errors.New("position is surrounded by opponent stones and would not cause any openent stones to be captured")
	}

	if zeroLiberty && !zeroLibertyOpponent {
		return false, errors.New("position is would make liberty count of own stone group 0")
	}

	finalLiberty := int16(len(neighbors))
	reducedOpponentRoots := make(map[int16]bool)
	for _, neighbor := range neighbors {
		if board.Positions[position] == board.Positions[neighbor] {
			neighborRoot := board.findRoot(neighbor)
			if merged, newRoot := board.merge(position, neighborRoot); merged {
				finalLiberty += board.LibertyCount[neighborRoot] - 2
				board.LibertyCount[newRoot] = finalLiberty
			} else {
				finalLiberty -= 2
				board.LibertyCount[newRoot] = finalLiberty
			}
		} else if _, exits := reducedOpponentRoots[neighbor]; !exits {
			board.LibertyCount[board.findRoot(neighbor)] -= 1
		}
	}
	return true, nil
}

func (board *Board) capture(position int16) int16 {
	stack := make([]int16, 0)
	stack = append(stack, position)
	var captureCount int16 = 0
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		board.Positions[node] = Empty
		captureCount++

		neighbors := board.neighbors(position)
		for _, n := range neighbors {
			if board.Positions[n] == board.Positions[node] {
				stack = append(stack, n)
			}
		}
	}
	if board.IsBlackNext {
		board.BlackCapture += captureCount
	} else {
		board.WhiteCapture += captureCount
	}
	return captureCount
}

func (s *Board) findRoot(position int16) int16 {
	if s.Roots[position] < 0 {
		return -1
	}

	root := position
	for s.Roots[root] != root {
		root = s.Roots[root]
	}

	for s.Roots[position] != position {
		position, s.Roots[position] = s.Roots[position], root
	}

	return root
}

func (s *Board) merge(positionA int16, positionB int16) (bool, int16) {
	rootA := s.findRoot(positionA)
	rootB := s.findRoot(positionB)

	if rootA == rootB {
		return false, rootA
	}

	if s.GroupSizes[rootA] > s.GroupSizes[rootB] {
		rootA, rootB = rootB, rootA
		positionA, positionB = positionB, positionA
	}

	for s.Roots[positionA] != positionA {
		positionA, s.Roots[positionA] = s.Roots[positionA], rootB
	}
	s.GroupSizes[rootB] += s.GroupSizes[rootA]
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
