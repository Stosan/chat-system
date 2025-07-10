
// ========================================
// client.go - Client Service
package main

import (
	"bufio"

	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a chat message
type Message struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

// Client represents the chat client
type Client struct {
	conn   *websocket.Conn
	userID string
}

// NewClient creates a new client instance
func NewClient() *Client {
	return &Client{}
}

// Connect connects to the server
func (c *Client) Connect() error {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// Register registers the user with the server
func (c *Client) Register(userID string) error {
	c.userID = userID
	msg := Message{
		Text:      "",
		Sender:    userID,
		Receiver:  "",
		Type:      "new_client",
		Timestamp: time.Now(),
	}
	return c.conn.WriteJSON(msg)
}

// SendMessage sends a message to another user
func (c *Client) SendMessage(receiver, text string) error {
	msg := Message{
		Text:      text,
		Sender:    c.userID,
		Receiver:  receiver,
		Type:      "chat",
		Timestamp: time.Now(),
	}
	return c.conn.WriteJSON(msg)
}

// ListenForMessages listens for incoming messages
func (c *Client) ListenForMessages() {
	for {
		var msg Message
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
func (c *Client) Disconnect() {
	msg := Message{
		Text:      "",
		Sender:    c.userID,
		Receiver:  "",
		Type:      "session_end",
		Timestamp: time.Now(),
	}
	c.conn.WriteJSON(msg)
	c.conn.Close()
}

func main() {
	client := NewClient()

	err := client.Connect()
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your user ID: ")
	scanner.Scan()
	userID := scanner.Text()

	err = client.Register(userID)
	if err != nil {
		log.Fatal("Failed to register:", err)
	}

	// Start listening for messages in goroutine
	go client.ListenForMessages()

	fmt.Println("Chat started! Use format: to:<user_id> <message>")
	fmt.Println("Type 'quit' to exit")

	for {
		scanner.Scan()
		input := scanner.Text()

		if input == "quit" {
			client.Disconnect()
			break
		}

		if strings.HasPrefix(input, "to:") {
			parts := strings.SplitN(input[3:], " ", 2)
			if len(parts) >= 2 {
				receiver := parts[0]
				message := parts[1]
				err := client.SendMessage(receiver, message)
				if err != nil {
					fmt.Printf("Error sending message: %v\n", err)
				}
			} else {
				fmt.Println("Invalid format. Use: to:<user_id> <message>")
			}
		} else {
			fmt.Println("Invalid format. Use: to:<user_id> <message>")
		}
	}
}

