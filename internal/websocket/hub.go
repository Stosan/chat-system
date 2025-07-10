package websocket

import (
	"chatsystem/internal/models"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Server represents the central server
type Hub struct {
	clients     map[string]*Client
	mutex       sync.RWMutex
	persistChan chan models.Message
	chatChan    chan models.Message
	historyChan chan string
	upgrader    websocket.Upgrader
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		persistChan: make(chan models.Message, 100),
		chatChan:    make(chan models.Message, 100),
		historyChan: make(chan string, 100),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// RegisterUser registers a new user
func (s *Hub) RegisterUser(userID string, conn *websocket.Conn) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.clients[userID]; exists {
		return false
	}

	s.clients[userID] = &Client{
		ID:   userID,
		Conn: conn,
	}
	return true
}

// RemoveUser removes a user from the server
func (s *Hub) RemoveUser(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[userID]; exists {
		client.Conn.Close()
		delete(s.clients, userID)
	}
}

// GetClient returns a client by ID
func (s *Hub) GetClient(userID string) (*Client, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	client, exists := s.clients[userID]
	return client, exists
}

func (s *Hub) SendToChat(msg models.Message) {
	s.chatChan <- msg
}

func (s *Hub) SendToPersist(msg models.Message) {
	s.persistChan <- msg
}
