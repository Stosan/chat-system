package websocket

import (
	"chatsystem/internal/models"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ChatClient struct {
	conn      *websocket.Conn
	userID    string
	sendCh    chan models.Message
	receiveCh chan models.Message
	closeCh   chan struct{}
	mu        sync.RWMutex
	closed    bool
}

// NewClient creates a new client instance
func NewChatClient(userID string) *ChatClient {
	return &ChatClient{
		userID:    userID,
		sendCh:    make(chan models.Message, 100),
		receiveCh: make(chan models.Message, 100),
		closeCh:   make(chan struct{}),
	}
}

func (c *ChatClient) ChatConnect(ctx context.Context) error {
	wsURL := "ws://localhost:5100/ws/server"

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	// Start goroutines for read/write
	go c.readLoop()
	go c.writeLoop()

	return nil
}

func (c *ChatClient) ChatRegister() error {
	msg := models.Message{
		Sender:    c.userID,
		Type:      "new_client",
		Timestamp: time.Now(),
	}

	select {
	case c.sendCh <- msg:
		return nil
	case <-c.closeCh:
		return fmt.Errorf("client is closed")
	case <-time.After(5 * time.Second):
		return fmt.Errorf("registration timeout")
	}
}

func (c *ChatClient) SendChatMessage(receiver, text string) error {
	msg := models.Message{
		Sender:    c.userID,
		Receiver:  receiver,
		Text:      text,
		Type:      "chat",
		Timestamp: time.Now(),
	}

	select {
	case c.sendCh <- msg:
		return nil
	case <-c.closeCh:
		return fmt.Errorf("client is closed")
	case <-time.After(5 * time.Second):
		return fmt.Errorf("send timeout")
	}
}

func (c *ChatClient) GetMessages() <-chan models.Message {
	return c.receiveCh
}

func (c *ChatClient) readLoop() {
	defer c.cleanup()

	for {
		select {
		case <-c.closeCh:
			return
		default:
			c.mu.RLock()
			conn := c.conn
			c.mu.RUnlock()

			if conn == nil {
				return
			}

			conn.SetReadDeadline(time.Now().Add(60 * time.Second))

			var msg models.Message
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("Read error for user %s: %v", c.userID, err)
				return
			}

			select {
			case c.receiveCh <- msg:
			case <-c.closeCh:
				return
			default:
				// Drop message if channel is full
				log.Printf("Dropping message for user %s - channel full", c.userID)
			}
		}
	}
}

func (c *ChatClient) writeLoop() {
	defer c.cleanup()

	for {
		select {
		case <-c.closeCh:
			return
		case msg := <-c.sendCh:
			c.mu.RLock()
			conn := c.conn
			c.mu.RUnlock()

			if conn == nil {
				return
			}

			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Write error for user %s: %v", c.userID, err)
				return
			}
		}
	}
}

func (c *ChatClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true
	close(c.closeCh)

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *ChatClient) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

// === Client Manager ===
type ClientManager struct {
	clients map[string]*ChatClient
	mu      sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*ChatClient),
	}
}

func (cm *ClientManager) GetOrCreateClient(userID string) *ChatClient {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, exists := cm.clients[userID]; exists {
		return client
	}

	client := NewChatClient(userID)
	cm.clients[userID] = client
	return client
}

func (cm *ClientManager) RemoveClient(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, exists := cm.clients[userID]; exists {
		client.Close()
		delete(cm.clients, userID)
	}
}

func (cm *ClientManager) GetClient(userID string) (*ChatClient, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	client, exists := cm.clients[userID]
	return client, exists
}
