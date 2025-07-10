package websocket

import (
	"chatsystem/internal/models"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected client
type Client struct {
	ID   string
	Conn *websocket.Conn
}

// Connect connects to the server
func (c *ChatClient) Connect() error {
	u := url.URL{Scheme: "ws", Host: "localhost:5100", Path: "/ws/chat"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// Register registers the user with the server
func (c *ChatClient) Register(userID string) error {
	c.userID = userID
	msg := models.Message{
		Text:      "",
		Sender:    userID,
		Receiver:  "",
		Type:      "new_client",
		Timestamp: time.Now(),
	}
	return c.conn.WriteJSON(msg)
}

// SendMessage sends a message to another user
func (c *ChatClient) SendMessage(receiver, text string) error {
	msg := models.Message{
		Text:      text,
		Sender:    c.userID,
		Receiver:  receiver,
		Type:      "chat",
		Timestamp: time.Now(),
	}
	return c.conn.WriteJSON(msg)
}

// ListenForMessages listens for incoming messages
func (c *ChatClient) ListenForMessages() {
	for {
		var msg models.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch msg.Type {
		case "registration_success":
			fmt.Println("✓ Successfully registered!")
		case "registration_error":
			fmt.Println("✗ Registration failed:", msg.Text)
		case "chat":
			fmt.Printf("[%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), msg.Sender, msg.Text)
		default:
			fmt.Printf("Historical: %s\n", msg.Text)
		}
	}
}

// Disconnect disconnects from the server
func (c *ChatClient) Disconnect() {
	msg := models.Message{
		Text:      "",
		Sender:    c.userID,
		Receiver:  "",
		Type:      "session_end",
		Timestamp: time.Now(),
	}
	c.conn.WriteJSON(msg)
	c.conn.Close()
}
