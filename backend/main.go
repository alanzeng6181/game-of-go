package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	gamelogic "github.com/alanzeng6181/game-of-go/gameLogic"
	net "github.com/alanzeng6181/game-of-go/network"
	middlewares "github.com/alanzeng6181/game-of-go/network/middlewares"
	"github.com/alanzeng6181/game-of-go/security"
	"github.com/gorilla/websocket"
)

func login(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Unable to read request body"))
	}

	var credential = struct {
		Username string
		Password string
	}{}
	err = json.Unmarshal(bytes, &credential)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("request body need to be in json format"))
		return
	}

	if userId, err := profileManager.GetUserId(credential.Username, credential.Password); err == nil && userId != "" {
		tokenStr, err := security.GetToken(userId)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Unable to get asigned jwt token"))
		} else {
			w.WriteHeader(200)
			w.Write([]byte(tokenStr))
		}
		return
	}

	w.WriteHeader(401)
	w.Write([]byte("invalid credential"))
}

func loadProfile(w http.ResponseWriter, r *http.Request) {
	if player := getPlayer(r); player != nil {
		playingGame, obervingGames := gameManager.GetGames(player.PlayerId)
		bytes, err := json.Marshal(struct {
			PlayerId       int64
			UserId         string
			Level          int8
			CurrentGameId  *int64
			Wins           int32
			Losses         int32
			ObservingGames []int64
		}{
			PlayerId:       player.PlayerId,
			UserId:         player.UserId,
			Level:          player.Level,
			CurrentGameId:  playingGame,
			Wins:           player.Wins,
			Losses:         player.Losses,
			ObservingGames: obervingGames,
		})

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("error occurred when writing profile data"))
		} else {
			w.WriteHeader(200)
			w.Write(bytes)
		}
		return
	}
	w.WriteHeader(400)
	w.Write([]byte("unable to locate user"))
}

func getPlayer(r *http.Request) *gamelogic.Player {
	if user := r.Context().Value("user"); user != nil {
		return profileManager.GetPlayerByUserId(user.(string))
	}
	return nil
}

func findGame(w http.ResponseWriter, r *http.Request) {
	bodyContent, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("unable to read request body"))
	}
	var gameRequest gamelogic.GameRequest
	if err := json.Unmarshal(bodyContent, &gameRequest); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("request body is expected to be GameRequest json"))
	}

	player := getPlayer(r)
	if player == nil {
		w.WriteHeader(400)
		w.Write([]byte("playerId not specified"))
	}
	gameRequest.PlayerId = player.PlayerId

	gameChan, requestId := gameManager.Match(gameRequest)
	select {
	case game := <-gameChan:
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("{gameId:%d}", game.Id)))
	case <-time.After(10 * time.Minute):
		gameManager.CancelRequest(requestId)
		w.WriteHeader(400)
		w.Write([]byte("timed out, unable to find game"))
	}
}

func wsConnect(w http.ResponseWriter, r *http.Request) {
	player := getPlayer(r)

	if player == nil {
		w.WriteHeader(401)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go handleConnection(conn, player)
}

func handleConnection(conn *websocket.Conn, player *gamelogic.Player) {
	gameManager.SetConnection(player.PlayerId, conn)
	for {
		messageType, bytes, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if messageType == websocket.TextMessage {
			var message net.RequestMessage
			err := json.Unmarshal(bytes, &message)
			if err != nil {
				//handle err
			}
			if message.Command == net.Move {
				gameId, err := strconv.Atoi(message.Arguments[0])
				invalidMoveCommandErrorMsg := fmt.Sprintf("invalid move command: %v, expected ['gameId','row','col']", message.Arguments)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, errorResponseBytes(invalidMoveCommandErrorMsg))
				}
				row, err := strconv.Atoi(message.Arguments[1])
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, errorResponseBytes(invalidMoveCommandErrorMsg))
				}
				col, err := strconv.Atoi(message.Arguments[2])
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, errorResponseBytes(invalidMoveCommandErrorMsg))
				}
				gameManager.Move(int64(gameId), int16(row), int16(col), player.PlayerId)
			} else if message.Command == net.TakeBack {

			} else if message.Command == net.Resign {

			} else if message.Command == net.MessageCommand {

			} else if message.Command == net.Observe {
				gameId, err := strconv.Atoi(message.Arguments[0])
				invalidMoveCommandErrorMsg := fmt.Sprintf("invalid move command: %v, expected ['gameId','row','col']", message.Arguments)
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, errorResponseBytes(invalidMoveCommandErrorMsg))
				}
				gameManager.AddObserver(int64(player.PlayerId), int64(gameId))
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var gameManager = gamelogic.NewGameManager()

var profileManager = gamelogic.NewProfileManager()

func main() {
	http.Handle("/api/findgame", middlewares.ApplyAuth(http.HandlerFunc(findGame)))
	http.HandleFunc("/api/login", login)
	http.Handle("/api/ws", middlewares.ApplyAuth(http.HandlerFunc(wsConnect)))
	http.Handle("/api/profile", middlewares.ApplyAuth(http.HandlerFunc(loadProfile)))
	log.Println("listening at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func errorResponseBytes(errorMessage string) []byte {
	bytes, _ := json.Marshal(net.ResponseMessage{
		ResponseType: net.Error,
		Content:      errorMessage,
	})
	return bytes
}
