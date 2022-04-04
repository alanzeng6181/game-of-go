package gamelogic

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	connections map[int64]*websocket.Conn
	mu          sync.RWMutex
}

func (connectionManager *ConnectionManager) AddConnection(playerId int64, conn *websocket.Conn) {
	connectionManager.mu.Lock()
	defer connectionManager.mu.Unlock()
	if oldVal, ok := connectionManager.connections[playerId]; ok {
		oldVal.Close()
	}
	connectionManager.connections[playerId] = conn
}

func (connectionManager *ConnectionManager) GetConnection(playerId int64) *websocket.Conn {
	return connectionManager.connections[playerId]
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[int64]*websocket.Conn),
	}
}
