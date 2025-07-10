# chat-system

// ========================================
// go.mod
module real-time-chat

go 1.21

require github.com/gorilla/websocket v1.5.0

// ========================================
// README.md
# Real-Time Chat Application in Go

A distributed real-time chat application built with Go, featuring WebSocket connections, RPC communication, and concurrent message processing.

## Architecture

The application consists of four main components:

1. **Central Server** (`main.go`) - Handles WebSocket connections and coordinates between services
2. **Client** (`client.go`) - Terminal-based chat client for users
3. **Persistence Service** (`persistence.go`) - Stores messages to local files via RPC
4. **History Service** (`history.go`) - Retrieves historical messages via RPC

## Features

- Real-time messaging using WebSockets
- User registration and session management
- Message persistence and historical retrieval
- Concurrent processing using goroutines and channels
- RPC communication between services
- Graceful handling of service failures

## Installation

```bash
go mod init chatsystem
go get github.com/gorilla/websocket
```

## Usage

1. Start the persistence service:
```bash
go run persistence.go
```

2. Start the history service:
```bash
go run history.go
```

3. Start the central server:
```bash
go run main.go
```

4. Run multiple clients:
```bash
go run client.go
```

## Message Format

Messages use the format: `to:<user_id> <message>`

Example: `to:albusdd Hello there!`

## Key Concepts

- **Goroutines**: Lightweight threads for concurrent execution
- **Channels**: Communication mechanism between goroutines
- **WebSockets**: Real-time bidirectional communication
- **RPC**: Remote procedure calls for inter-service communication
- **JSON**: Message serialization format