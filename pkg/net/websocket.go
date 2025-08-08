package net

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type MessageChannel interface {

	// ID returns the unique identifier for the WebSocket connection
	ID() string

	// ChangeID changes the client's ID to a new value.
	// If the provided ID is empty, it returns an error and keeps the current ID.
	ChangeID(id string) error

	// SendMessage sends a message to the WebSocket connection
	SendMessage(message []byte) error

	// SendBinaryMessage sends a binary message to the WebSocket connection
	SendBinaryMessage(message []byte) error

	// ReceiveMessage returns the next message received from the WebSocket connection.
	//
	// Usage:
	// for {
	// 	message, err := handler.ReceiveMessage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	fmt.Printf("Received message: %s\n", message)
	// 	// Process the message
	// 	// ...
	// }
	//
	ReceiveMessage() ([]byte, error)

	Close() error
}

type HubChannel interface {

	// SendTo sends a message to a specific client identified by receiveId.
	// It supports both text (1) and binary messages (2).
	SendTo(receiveId string, messageType int, message []byte) error

	// BroadcastMessage sends a message to all connected clients
	BroadcastMessage(message []byte) error

	BroadcastBinary(b []byte) error
}

// UpgradeToWebSocket upgrades the HTTP connection to a WebSocket connection.
// It uses the global hub connection to register the client.
func UpgradeToWebSocket(w http.ResponseWriter, r *http.Request) (HubChannel, MessageChannel, error) {
	// init globalHubConnection
	hub, err := globalHubConnection()
	if err != nil {
		log.Printf("Failed to get hub connection: %v", err)
		return nil, nil, err
	}
	// Upgrade HTTP connection to WebSocket with global hub
	return UpgradeToWebSocketCustom(hub, w, r)
}

// UpgradeToWebSocketCustom upgrades the HTTP connection to a WebSocket connection with a specific hub
//
// Example usage:
//
//	hub := NewHub()
//
//	// Upgrade the HTTP connection to a WebSocket connection with the hub in a handler function
//	func handler(w http.ResponseWriter, r *http.Request) {
//		ch, mh, err := UpgradeToWebSocketCustom(hub, w, r)
//		if err != nil {
//		    log.Printf("Failed to upgrade to WebSocket: %v", err)
//		    return
//		}
//		// Now you can use mh to send and receive messages
//		...
//	}
func UpgradeToWebSocketCustom(hub *Hub, w http.ResponseWriter, r *http.Request) (HubChannel, MessageChannel, error) {
	// utilize remoteAddr for client identification
	remoteAddr := strings.Map(func(r rune) rune {
		if r > 'a' && r < 'z' || r > 'A' && r < 'Z' || r > '0' && r < '9' || r == '-' || r == '_' {
			return r
		}
		return '0' // Replace non-alphanumeric characters with '0'
	}, r.RemoteAddr)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return nil, nil, err
	}
	// Register client with hub
	client := hub.registerClient(conn, fmt.Sprintf("CID-%s-%d", remoteAddr, time.Now().Unix()))
	return hub, client, nil
}
