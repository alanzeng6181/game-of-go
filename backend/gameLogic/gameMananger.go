package gamelogic

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type GameManager struct {
	gameRequests []struct {
		gameReqeust GameRequest
		otherChan   chan GameMatchResult
	}
	lock              sync.RWMutex
	games             map[int64]*Game
	requestIdCounter  int64
	gameId            int64
	connectionManager *ConnectionManager
}

func (gameManager *GameManager) AddObserver(playerId int64, gameId int64) {
	game := gameManager.games[gameId]
	if game.BlackPlayer.PlayerId != playerId && game.WhitePlayer.PlayerId != playerId {
		game.Observers = append(game.Observers, &Player{PlayerId: playerId})
	}
}

func (gameManager *GameManager) SetConnection(playerId int64, conn *websocket.Conn) {
	gameManager.connectionManager.AddConnection(playerId, conn)
}

func (gameManager *GameManager) Move(gameId int64, row int16, col int16, playerId int64) error {
	if game, ok := gameManager.games[gameId]; ok {
		color := White
		if game.BlackPlayer.PlayerId == playerId {
			color = Black
		}
		gameState, err := game.Move1(row, col, color)
		if err != nil {
			return err
		}
		gameManager.broadcast(game, gameState)
		return nil
	} else {
		return errors.New("unable to make a move")
	}
}

func (gameManager *GameManager) broadcast(game *Game, gameState *GameState) {
	allPlayers := make([]*Player, 0)
	allPlayers = append(allPlayers, game.BlackPlayer)
	allPlayers = append(allPlayers, game.WhitePlayer)
	allPlayers = append(allPlayers, game.Observers...)

	for _, player := range allPlayers {
		conn := gameManager.connectionManager.GetConnection(player.PlayerId)
		if conn != nil {
			bytes, _ := json.Marshal(gameState)
			conn.WriteMessage(websocket.TextMessage, bytes)
		}
	}
}

func (gameManager *GameManager) Comment(gameId int64, playerId int64, comment string) error {
	if game, ok := gameManager.games[gameId]; ok {

		game.Comments = append(game.Comments, Comment{Timestamp: time.Now().Unix(), Content: comment})

		gameState := game.GetGameState()
		gameManager.broadcast(game, &gameState)

		return nil
	} else {
		return errors.New("unable to make a move")
	}
}

func (gameManager *GameManager) GetGames(playerId int64) (*int64, []int64) {
	var playingGame *int64
	observingGames := make([]int64, 0)
	for gameId, game := range gameManager.games {
		if game.BlackPlayer.PlayerId == playerId || game.WhitePlayer.PlayerId == playerId {
			playingGame = &gameId
		} else {
			for _, o := range game.Observers {
				if o.PlayerId == playerId {
					observingGames = append(observingGames, gameId)
					break
				}
			}
		}
	}
	return playingGame, observingGames
}

func NewGameManager() *GameManager {
	return &GameManager{
		gameRequests: make([]struct {
			gameReqeust GameRequest
			otherChan   chan GameMatchResult
		}, 0),
		games:             make(map[int64]*Game),
		connectionManager: NewConnectionManager(),
	}
}

type GameMatchResult struct {
	*Game
	Stone
}

func (gameManager *GameManager) GetGame(gameId int64) *Game {
	return gameManager.games[gameId]
}

func (gameManager *GameManager) Match(gameRequest GameRequest) (<-chan GameMatchResult, int64) {
	requestId := atomic.AddInt64(&gameManager.requestIdCounter, 1)
	var gameChan chan GameMatchResult = make(chan GameMatchResult, 1)
	for _, gr := range gameManager.gameRequests {
		if gr.gameReqeust.IsMatch(gameRequest) {
			game, err := NewGame(ClockConfig{TimeLimit: time.Duration(gameRequest.TimeoutInSeconds) * time.Second,
				OverTime:      time.Duration(gameRequest.OverTimeInSeconds) * time.Second,
				OverTimeCount: gameRequest.OverTimeCount}, int16(gameRequest.BoardSize), &Player{UserId: "abc", PlayerId: gameRequest.PlayerId}, &Player{
				UserId:   "edf",
				PlayerId: gr.gameReqeust.PlayerId,
			})
			if err != nil {
				log.Println(err)
				break
			}
			gr.otherChan <- struct {
				*Game
				Stone
			}{game, White}
			gameChan <- GameMatchResult{game, Black}
			newGameId := atomic.AddInt64(&gameManager.gameId, 1)
			game.Id = newGameId
			gameManager.games[newGameId] = game
			return (<-chan GameMatchResult)(gameChan), requestId
		}
	}
	gameRequest.RequestId = requestId
	gameManager.gameRequests = append(gameManager.gameRequests, struct {
		gameReqeust GameRequest
		otherChan   chan GameMatchResult
	}{gameReqeust: gameRequest, otherChan: gameChan})
	return (<-chan GameMatchResult)(gameChan), requestId
}

func (gameManager *GameManager) CancelRequest(requestId int64) bool {
	gameManager.lock.RLock()
	defer gameManager.lock.Unlock()

	delete := -1
	for i, gr := range gameManager.gameRequests {
		if gr.gameReqeust.RequestId == requestId {
			delete = i
		}
	}

	if delete >= 0 {
		gameManager.gameRequests = append(gameManager.gameRequests[:delete], gameManager.gameRequests[delete+1:]...)
		return true
	}
	return false
}
