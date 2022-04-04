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
	"github.com/gorilla/websocket"
)

func login(w http.ResponseWriter, r *http.Request) {
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
	query := r.URL.Query()
	gameIdStr := query.Get("gameId")
	playerIdStr := query.Get("playerId")
	if gameIdStr == "" || playerIdStr == "" {
		w.WriteHeader(400)
		w.Write([]byte("gameId and playerId must be specified"))
	}
	//Replace
	gameId, _ := strconv.Atoi(gameIdStr)
	playerId, _ := strconv.Atoi(playerIdStr)

	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if err != nil {
		log.Println(err)
		return
	}
	gameManager.AddConnection(int64(playerId), conn)

	for {
		messageType, bytes, err := conn.ReadMessage()
		if err != nil {
			//handle err
		}
		if messageType == websocket.TextMessage {
			var message net.RequestMessage
			err := json.Unmarshal(bytes, &message)
			if err != nil {
				//handle err
			}
			if message.Command == net.Move {
				row, err := strconv.Atoi(message.Arguments[0])
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, parseErrorResponse(fmt.Sprintf("invalid move, %v", message.Arguments)))
				}
				col, err := strconv.Atoi(message.Arguments[1])
				if err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("invalid move, %v", message.Arguments)))
				}

				gameManager.Move(int64(gameId), int16(row), int16(col), int64(playerId))
			} else if message.Command == net.TakeBack {

			} else if message.Command == net.Resign {

			} else if message.Command == net.MessageCommand {

			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var gameManager = gamelogic.NewGameManager()

func main() {
	http.HandleFunc("/findgame", findGame)
	http.HandleFunc("/login", login)
	http.HandleFunc("/ws", wsConnect)
	log.Println("listening at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseErrorResponse(errorMessage string) []byte {
	return []byte(errorMessage)
}
