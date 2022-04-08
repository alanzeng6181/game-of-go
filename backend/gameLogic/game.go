package gamelogic

import (
	"errors"
	"math/rand"
)

type Game struct {
	Id             int64
	Board          *Board
	BlackPlayer    *Player
	WhitePlayer    *Player
	Observers      []*Player
	Winner         *Player
	Comments       []Comment
	Clock          *Clock
	clockInputChan chan<- Stone
	Status         Status
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
		Board:          NewBoard(boardSize),
		BlackPlayer:    playerA,
		WhitePlayer:    playerB,
		Comments:       make([]Comment, 0),
		Clock:          clock,
		Status:         NotStarted,
		clockInputChan: input,
	}, err
}

func (game *Game) Move(position int16, color Stone) (*GameState, error) {

	if game.Board.IsBlackNext && color == White {
		return nil, errors.New("it's Black's turn")
	}
	if !game.Board.IsBlackNext && color == Black {
		return nil, errors.New("it's White's turn")
	}

	if game.Status != NotStarted && game.Status != InProgress {
		return nil, errors.New("game is over, no more moves allowed")
	}

	if game.Status == NotStarted {
		game.Status = InProgress
		game.Clock.Start()
	}

	game.clockInputChan <- color

	if _, err := game.Board.Move(position); err != nil {
		return nil, err
	}

	return &GameState{
		Positions:     game.Board.Positions,
		WhiteCaptures: int8(game.Board.WhiteCapture),
		BlackCaptures: int8(game.Board.WhiteCapture),
		Status:        game.Status,
		ClockState:    game.Clock.ToClockState(),
		Comments:      game.Comments,
	}, nil
}

func (game *Game) GetGameState() GameState {
	return GameState{
		Positions:     game.Board.Positions,
		WhiteCaptures: int8(game.Board.WhiteCapture),
		BlackCaptures: int8(game.Board.WhiteCapture),
		Status:        game.Status,
		ClockState:    game.Clock.ToClockState(),
		Comments:      game.Comments,
	}
}

func (game *Game) Move1(row int16, col int16, color Stone) (*GameState, error) {
	return game.Move(row*game.Board.Size+col, color)
}

func (game *Game) Resign(color Stone) {

}

type Status string

const (
	NotStarted   Status = "NotStarted"
	InProgress   Status = "InProgress"
	Over         Status = "Over"
	BlackWon     Status = "BlackWon"
	WhiteWon     Status = "WhiteWon"
	BlackResign  Status = "BlackResign"
	WhiteResign  Status = "WhiteResign"
	BlackTimeout Status = "BlackTimeout"
	WhiteTimeout Status = "WhiteTimeout"
)
