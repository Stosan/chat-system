// main.go - Central Server
package main

import (
"chatsystem"
	"fmt"
	"log"
	"net/http"

)




func main() {
	server := chatsystem.NewServer()

	// Start goroutines to process channels
	go server.ProcessChatMessages()
	go server.ProcessPersistMessages()

	http.HandleFunc("/ws", server.HandleConnection)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
