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
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

func login(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Unable to read request body"))
	}

	var credential = struct {
		username string
		password string
	}{}
	err = json.Unmarshal(bytes, &credential)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("request body need to be in json format"))
	}

	if userId, err := profileManager.GetPlayerId(credential.username, credential.password); err == nil && playerId >= 0 {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": userId,
			"nbf":  time.Now().Add(24 * time.Hour).Unix(),
		})
		tokenString, err := token.SignedString(middlewares.JWTSigningKey)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Unable to get asigned jwt token"))
		}
		w.WriteHeader(200)
		w.Write([]byte(tokenString))
	}

	w.WriteHeader(401)
	w.Write([]byte("invalid credential"))
}

func loadProfile(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("user")
}

func (r *http.Request) getPlayer() *Player {
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

	player := r.getPlayer()
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
	if err != nil {
		log.Println(err)
		return
	}
	gameManager.AddPlayer(int64(playerId), int64(gameId), conn)
	go handleConnection(conn, gameId, playerId)
}

func handleConnection(conn *websocket.Conn, gameId int, playerId int) {
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
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var gameManager = gamelogic.NewGameManager()

var profileManager = gamelogic.NewProfileManager()

func main() {
	http.Handle("/api/findgame", middlewares.ApplyAuth(findGame))
	http.Handle("/api/login", login)
	http.Handle("/api/ws", middlewares.ApplyAuth(wsConnect))
	http.Handle("/api/profile", middlewares.ApplyAuth(loadProfile))
	log.Println("listening at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func parseErrorResponse(errorMessage string) []byte {
	return []byte(errorMessage)
}
