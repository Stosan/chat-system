package chatsystem

import (
	"encoding/json"
	"log"
	"net/http"
	"net/rpc"
	"sync"
	"time"
appl "chatsystem/application"
	"github.com/gorilla/websocket"
)

// Client represents a connected client
type Client struct {
	ID   string
	Conn *websocket.Conn
}

// Server represents the central server
type Server struct {
	clients       map[string]*Client
	mutex         sync.RWMutex
	persistChan   chan appl.Message
	chatChan      chan appl.Message
	historyChan   chan string
	upgrader      websocket.Upgrader
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		clients:     make(map[string]*Client),
		persistChan: make(chan appl.Message, 100),
		chatChan:    make(chan appl.Message, 100),
		historyChan: make(chan string, 100),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// RegisterUser registers a new user
func (s *Server) RegisterUser(userID string, conn *websocket.Conn) bool {
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
func (s *Server) RemoveUser(userID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[userID]; exists {
		client.Conn.Close()
		delete(s.clients, userID)
	}
}

// GetClient returns a client by ID
func (s *Server) GetClient(userID string) (*Client, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	client, exists := s.clients[userID]
	return client, exists
}

// FetchHistoricalMessages fetches historical messages via RPC
func (s *Server) FetchHistoricalMessages(userID string) {
	client, err := rpc.DialHTTP("tcp", "localhost:8082")
	if err != nil {
		log.Printf("Error connecting to history service: %v", err)
		return
	}
	defer client.Close()

	var messages []string
	err = client.Call("HistoryService.GetMessages", appl.QueryMessage{UserID: userID}, &messages)
	if err != nil {
		log.Printf("Error fetching historical messages: %v", err)
		return
	}

	if clientConn, exists := s.GetClient(userID); exists {
		for _, msg := range messages {
			clientConn.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

// PersistMessage sends message to persistence service via RPC
func (s *Server) PersistMessage(msg appl.Message) {
	client, err := rpc.DialHTTP("tcp", "localhost:8081")
	if err != nil {
		log.Printf("Error connecting to persistence service: %v", err)
		return
	}
	defer client.Close()

	msgJSON, _ := json.Marshal(msg)
	persistMsg := appl.PersistMessage{
		Receiver: msg.Receiver,
		Message:  string(msgJSON),
	}

	var result string
	err = client.Call("PersistenceService.SaveMessage", persistMsg, &result)
	if err != nil {
		log.Printf("Error persisting message: %v", err)
	}
}

// HandleConnection handles WebSocket connections
func (s *Server) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	for {
		var msg appl.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		msg.Timestamp = time.Now()

		switch msg.Type {
		case "new_client":
			if s.RegisterUser(msg.Sender, conn) {
				response := appl.Message{
					Text:      "Registration successful",
					Sender:    "server",
					Receiver:  msg.Sender,
					Type:      "registration_success",
					Timestamp: time.Now(),
				}
				conn.WriteJSON(response)

				// Fetch historical messages in goroutine
				go s.FetchHistoricalMessages(msg.Sender)
			} else {
				response := appl.Message{
					Text:      "User already exists",
					Sender:    "server",
					Receiver:  msg.Sender,
					Type:      "registration_error",
					Timestamp: time.Now(),
				}
				conn.WriteJSON(response)
			}

		case "chat":
			// Send to chat channel for delivery
			s.chatChan <- msg
			// Send to persist channel for storage
			s.persistChan <- msg

		case "session_end":
			s.RemoveUser(msg.Sender)
			return
		}
	}
}

// ProcessChatMessages processes chat messages from channel
func (s *Server) ProcessChatMessages() {
	for msg := range s.chatChan {
		if receiver, exists := s.GetClient(msg.Receiver); exists {
			receiver.Conn.WriteJSON(msg)
		}
	}
}

// ProcessPersistMessages processes persistence messages from channel
func (s *Server) ProcessPersistMessages() {
	for msg := range s.persistChan {
		go s.PersistMessage(msg)
	}
}