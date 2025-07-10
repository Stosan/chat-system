package main
// ========================================


import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

// PersistMessage represents message for persistence
type PersistMessage struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

// PersistenceService handles message persistence
type PersistenceService struct{}

// SaveMessage saves a message to file
func (ps *PersistenceService) SaveMessage(msg PersistMessage, result *string) error {
	// Create directory if it doesn't exist
	dir := "messages"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create file path based on receiver ID
	filePath := filepath.Join(dir, fmt.Sprintf("%s.json", msg.Receiver))

	// Read existing messages
	var messages []string
	if data, err := os.ReadFile(filePath); err == nil {
		json.Unmarshal(data, &messages)
	}

	// Append new message
	messages = append(messages, msg.Message)

	// Write back to file
	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	*result = "Message saved successfully"
	return nil
}

func main() {
	ps := &PersistenceService{}
	rpc.Register(ps)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error starting persistence service:", err)
	}

	fmt.Println("Persistence service starting on :8081")
	log.Fatal(http.Serve(listener, nil))
}
