package gamelogic

import (
	"errors"
	"math/rand"
)

type Game struct {
	Id          int64
	Board       *Board
	BlackPlayer *Player
	WhitePlayer *Player
	Winner      *Player
	Messages    []string
	Clock       *Clock
	Status      Status
}

func NewGame(clockConfig ClockConfig, boardSize int16, playerA *Player, playerB *Player) (*Game, error) {
	num := rand.Int31n(1000)
	if num < 500 {
		playerA, playerB = playerB, playerA
	}
	input := make(chan Stone)
	clock, err := NewClock(clockConfig, input)
	if err != nil {
		return nil, err
	}
	return &Game{
		Board:       NewBoard(boardSize),
		BlackPlayer: playerA,
		WhitePlayer: playerB,
		Messages:    make([]string, 0),
		Clock:       clock,
		Status:      NotStarted,
	}, err
}

func (game *Game) Move(position int16, color Stone) ([]Stone, error) {
	if game.Board.IsBlackNext && color == White {
		return nil, errors.New("it's Black's turn")
	}
	if !game.Board.IsBlackNext && color == White {
		return nil, errors.New("it's White's turn")
	}
	if _, err := game.Board.move(position); err != nil {
		return nil, err
	}

	return game.Board.Positions, nil
}

func (game *Game) Resign(color Stone) {

}

type Status int8

const (
	NotStarted Status = iota
	InProgress
	Over
	BlackWon
	WhiteWon
	BlackResign
	WhiteResign
	BlackTimeout
	WhiteTimeout
)
