package gameLogic_test

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	gl "github.com/alanzeng6181/game-of-go/gameLogic"
)

func testBoard(state string, t *testing.T, action func(*gl.Board, *testing.T)) {
	/*
		example string. we start with a black stone, and alternate, the sequence is randomized.

		asterik character is treated as '-', it is used to indicate next move after the initiazation in our test cases.

		b w - - -
		* - - - -
		- b - - -
		- - - - -
		- - - - -
	*/

	var p int16 = 0
	blackPs := make([]int16, 0)
	whitePs := make([]int16, 0)

	positionCounts := 0
	for _, c := range state {
		if c == ' ' || c == '\t' || c == '\n' {
			continue
		}

		if c == 'b' || c == 'B' {
			blackPs = append(blackPs, p)
		} else if c == 'w' || c == 'W' {
			whitePs = append(whitePs, p)
		} else if c != '-' && c != '*' {
			t.Errorf("Invalid input, unexpected character: %c, only space, tab, newline and *, b, B, w, W are allowed", c)
		}
		p++
		positionCounts++
	}

	sqrt := math.Sqrt(float64(positionCounts))
	if sqrt != float64(int16(sqrt)) {
		t.Errorf("Invalid input, position count %d is not perfect square", positionCounts)
	}

	if len(blackPs)-len(whitePs) != 0 && len(blackPs)-len(whitePs) != 1 {
		t.Error("Invalid input, there needs to be 0 or 1 more number of black stone than white stone")
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(blackPs), func(i, j int) { blackPs[i], blackPs[j] = blackPs[j], blackPs[i] })
	rand.Shuffle(len(whitePs), func(i, j int) { whitePs[i], whitePs[j] = whitePs[j], whitePs[i] })

	board := gl.NewBoard(int16(sqrt))

	if !board.IsBlackNext {
		t.Errorf("black is expected to start first")
	}

	for len(blackPs) != 0 || len(whitePs) != 0 {
		var p int16
		if board.IsBlackNext {
			p = blackPs[0]
			blackPs = blackPs[1:]
		} else {
			p = whitePs[0]
			whitePs = whitePs[1:]
		}
		if ok, _ := board.Move(p); !ok {
			t.Errorf("expected position %d to be a valid move, but was not", p)
		} else {
			fmt.Printf("%d ", p)
		}
	}

	if action != nil {
		action(board, t)
	}
}

func TestBoardCornerCapture(t *testing.T) {
	//(0,0) should be captured

	const state = `
		b w - - -
		* - b - -
		- - - - -
		- - - - -
		- - - - -
	`

	testBoard(state, t, func(board *gl.Board, t *testing.T) {
		var p int16 = 5
		if ok, _ := board.Move(p); !ok {
			t.Errorf("expected position %d to be a valid move, but was not", p)
		}

		if board.BlackCapture != 1 {
			t.Errorf("expected black captures to be 1, but was %d", board.BlackCapture)
		}
	})
}

func TestBoardInvalidMoveNoLiberty(t *testing.T) {
	// (0,0)shoudl fail, because it would have no liberty
	const state = `
		* w b - -
		w - b - -
		- - - - -
		- - - - -
		- - - - -
	`
	testBoard(state, t, func(board *gl.Board, t *testing.T) {
		message := "position is surrounded by opponent stones and would not cause any openent stones to be captured"
		if ok, err := board.Move(0); ok || !strings.Contains(err.Error(), message) {
			t.Errorf("expected invalid move with message => %s", message)
		}
	})
}

func TestBoardNoLibertyButCapture(t *testing.T) {
	// attempted move should succeed, because it would capture (2,2) and cause it to have 1 liberty

	const state = `
		b - w - -
		- w * w -
		- b w b -
		- - b - -
		- - - - -
	`

	testBoard(state, t, func(board *gl.Board, t *testing.T) {
		var p int16 = 7
		if ok, err := board.Move(p); !ok {
			t.Errorf("expected valid move but got => %s", err.Error())
		}
		if board.WhiteCapture != 1 {
			t.Errorf("expected black captures to be 1, but was %d", board.BlackCapture)
		}
	})
}

func TestBoardLibertyCount(t *testing.T) {
	// attempted move should succeed, because it would capture (2,2) and cause it to have 1 liberty

	const state = `
		b - w - -
		- w * w -
		- b w b -
		- - b - -
		- - b - -
	`
	testBoard(state, t, func(board *gl.Board, t *testing.T) {
		assertEqual(2, board.GetLiberty(0, 0), t)
		assertEqual(3, board.GetLiberty(0, 2), t)
		assertEqual(3, board.GetLiberty(1, 1), t)
		assertEqual(3, board.GetLiberty(1, 3), t)
		assertEqual(2, board.GetLiberty(2, 1), t)
		assertEqual(1, board.GetLiberty(2, 2), t)
		assertEqual(2, board.GetLiberty(2, 3), t)
		assertEqual(4, board.GetLiberty(3, 2), t)
		assertEqual(4, board.GetLiberty(4, 2), t)

		var p int16 = 7
		if ok, err := board.Move(p); !ok {
			t.Errorf("expected valid move but got => %s", err.Error())
		}

		assertEqual(2, board.GetLiberty(0, 0), t)
		assertEqual(4, board.GetLiberty(0, 2), t)
		assertEqual(4, board.GetLiberty(1, 1), t)
		assertEqual(4, board.GetLiberty(1, 3), t)
		assertEqual(2, board.GetLiberty(2, 1), t)
		assertEqual(4, board.GetLiberty(2, 2), t)
		assertEqual(2, board.GetLiberty(2, 3), t)
		assertEqual(4, board.GetLiberty(3, 2), t)
		assertEqual(4, board.GetLiberty(4, 2), t)
	})
}

func assertEqual[T comparable](expected T, actual T, t *testing.T) {
	if expected != actual {
		t.Errorf("expected %v, but got %v", expected, actual)
	}
}

func TestBoardLibertyCountSpecficOrder(t *testing.T) {
	// attempted move should succeed, because it would capture (2,2) and cause it to have 1 liberty

	/*
		const state = `
			b - w - -
			- w * w -
			- b w b -
			- - b - -
			- - b - -
		`*/
	board := gl.NewBoard(5)
	moves := []int16{17, 2, 13, 8, 0, 6, 22, 12, 11}
	for _, p := range moves {
		board.Move(p)
	}
	assertEqual(2, board.GetLiberty(0, 0), t)
	assertEqual(3, board.GetLiberty(0, 2), t)
	assertEqual(3, board.GetLiberty(1, 1), t)
	assertEqual(3, board.GetLiberty(1, 3), t)
	assertEqual(2, board.GetLiberty(2, 1), t)
	assertEqual(1, board.GetLiberty(2, 2), t)
	assertEqual(2, board.GetLiberty(2, 3), t)
	assertEqual(4, board.GetLiberty(3, 2), t)
	assertEqual(4, board.GetLiberty(4, 2), t)
}

func TestBoardLibertyCountSpecficOrder2(t *testing.T) {
	// attempted move should succeed, because it would capture (2,2) and cause it to have 1 liberty

	/*
		const state = `
			b - w - -
			- w * w -
			- b w b -
			- - b - -
			- - b - -
		`*/
	board := gl.NewBoard(5)
	moves := []int16{11, 8, 0, 6, 22, 12, 13, 2, 17}
	for _, p := range moves {
		board.Move(p)
	}
	assertEqual(2, board.GetLiberty(0, 0), t)
	assertEqual(3, board.GetLiberty(0, 2), t)
	assertEqual(3, board.GetLiberty(1, 1), t)
	assertEqual(3, board.GetLiberty(1, 3), t)
	assertEqual(2, board.GetLiberty(2, 1), t)
	assertEqual(1, board.GetLiberty(2, 2), t)
	assertEqual(2, board.GetLiberty(2, 3), t)
	assertEqual(4, board.GetLiberty(3, 2), t)
	assertEqual(4, board.GetLiberty(4, 2), t)
}
