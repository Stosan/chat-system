package main

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

// QueryMessage represents query for historical messages
type QueryMessage struct {
	UserID string `json:"user_id"`
}

// HistoryService handles historical message retrieval
type HistoryService struct{}

// GetMessages retrieves historical messages for a user
func (hs *HistoryService) GetMessages(query QueryMessage, result *[]string) error {
	// Create file path based on user ID
	filePath := filepath.Join("messages", fmt.Sprintf("%s.json", query.UserID))

	// Read messages from file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			*result = []string{}
			return nil
		}
		return err
	}

	var messages []string
	err = json.Unmarshal(data, &messages)
	if err != nil {
		return err
	}

	*result = messages
	return nil
}

func main() {
	hs := &HistoryService{}
	rpc.Register(hs)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("Error starting history service:", err)
	}

	fmt.Println("History service starting on :8082")
	log.Fatal(http.Serve(listener, nil))
}

