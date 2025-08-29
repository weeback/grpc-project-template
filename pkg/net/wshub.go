package net

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// WebSocket upgrader with basic configuration
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from any origin in development
			// In production, you should implement proper origin checking
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	hubInstance *Hub
	hubOnce     = sync.Once{}

	writeTimeout = 30 * time.Second
	readTimeout  = 60 * time.Second
	idleTimeout  = 120 * time.Second
	pingCycle    = 54 * time.Second // Ping interval to keep the connection alive
)

func globalHubConnection() (*Hub, error) {
	hubOnce.Do(func() {
		hubInstance = NewHub()
	})
	if hubInstance == nil {
		return nil, fmt.Errorf("failed to create hub instance")
	}
	return hubInstance, nil
}

// NewHub creates a new Hub
func NewHub() *Hub {

	// Initialize the hub with channels and client map
	hub := &Hub{
		broadcast:    make(chan []byte, 256),
		broadcastBin: make(chan []byte, 256),
		register:     make(chan *Client, 256),
		unregister:   make(chan *Client, 256),
		clients:      make(map[*Client]bool),
		once:         sync.Once{}, // No pending clients initially
	}

	// Start the hub in a goroutine
	go hub.run()
	return hub
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	clients      map[*Client]bool
	broadcast    chan []byte
	broadcastBin chan []byte
	register     chan *Client
	unregister   chan *Client
	once         sync.Once
}

// Run starts the hub and handles client registration/unregistration and broadcasting
func (h *Hub) run() {
	h.once.Do(func() {
		entry := getLogEntry()
		// Start the hub's main loop
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
				entry.Debug("Client connected",
					zap.String("client_id", client.id), zap.Int("total_clients", len(h.clients)))

			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
					close(client.recv)
					entry.Debug("Client disconnected",
						zap.String("client_id", client.id),
						zap.Int("total_clients", len(h.clients)))
				} else {
					entry.Warn("Client not found in hub or already disconnected",
						zap.String("client_id", client.id))
				}

			case message := <-h.broadcast:
				// Broadcast message to all connected clients
				for client := range h.clients {
					if err := client.write(websocket.TextMessage, message); err != nil {
						entry.Error("Error broadcast sending text message to client",
							zap.String("client_id", client.id),
							zap.Error(err))
					}
				}
			case bin := <-h.broadcastBin:
				// Broadcast binary message to all connected clients
				for client := range h.clients {
					if err := client.write(websocket.BinaryMessage, bin); err != nil {
						entry.Error("Error broadcast sending binary message to client",
							zap.String("client_id", client.id),
							zap.Error(err))
					}
				}
			}
		}
	})
}

func (h *Hub) registerClient(conn *websocket.Conn, id string) *Client {

	// Ensure the hub is initialized
	// Create new client
	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
		recv: make(chan []byte, 256),
		id:   id,
	}
	// If id is empty, generate a unique ID
	if id == "" {
		client.id = fmt.Sprintf("client-%s-%d", conn.RemoteAddr().String(), time.Now().UnixNano())
	}

	// Register client with hub
	h.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()

	return client
}

func (h *Hub) SendTo(receiveId string, messageType int, message []byte) error {
	for client := range h.clients {
		if client.id == receiveId {
			switch messageType {
			case websocket.TextMessage:
				return client.SendMessage(message)
			case websocket.BinaryMessage:
				return client.SendBinaryMessage(message)
			default:
				return fmt.Errorf("unsupported message type: %d, please use websocket.TextMessage (1) or websocket.BinaryMessage (2)", messageType)
			}
		}
	}
	return fmt.Errorf("client %s not found", receiveId)
}

func (h *Hub) BroadcastMessage(message []byte) error {
	select {
	case h.broadcast <- message:
		return nil
	case <-time.After(writeTimeout):
		return fmt.Errorf("failed to broadcast message: write timeout")
	default:
		// Notify hub to unregister the client
		return fmt.Errorf("failed to broadcast message: closing connection")
	}
}

func (h *Hub) BroadcastBinary(b []byte) error {
	select {
	case h.broadcastBin <- b:
		return nil
	case <-time.After(writeTimeout):
		return fmt.Errorf("failed to broadcast message: write timeout")
	default:
		// Notify hub to unregister the client
		return fmt.Errorf("failed to broadcast message: closing connection")
	}
}

// Client represents a WebSocket client connection
type Client struct {
	conn       *websocket.Conn
	send, recv chan []byte
	hub        *Hub
	id         string
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) ChangeID(id string) error {
	// Change the client's ID
	if id == "" {
		return fmt.Errorf("invalid ID provided for client %s, keeping the current ID", c.id)
	}
	c.id = id // Log the ID change
	getLogEntry().Info("Client ID changed", zap.String("client_id", c.id))
	return nil
}

func (c *Client) Close() error {
	// Close the WebSocket connection and notify the hub
	c.hub.unregister <- c
	// Close the channels to prevent further sends/receives
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) SendMessage(message []byte) error {
	return c.write(websocket.TextMessage, message)
}

func (c *Client) SendBinaryMessage(message []byte) error {
	return c.write(websocket.BinaryMessage, message)
}

func (c *Client) ReceiveMessage() ([]byte, error) {
	message, ok := <-c.recv
	if !ok {
		return nil, fmt.Errorf("client %s disconnected", c.id)
	}
	return message, nil
}

func (c *Client) write(messageType int, data []byte) error {
	var message []byte
	switch messageType {
	case websocket.TextMessage:
		if !utf8.Valid(data) {
			return fmt.Errorf("invalid UTF-8 data for text message")
		}
		message = append([]byte{websocket.TextMessage, 0xFF}, data...)
	case websocket.BinaryMessage:
		// Binary messages can contain any data, no validation needed
		message = append([]byte{websocket.BinaryMessage, 0xFF}, data...)
	default:
		return fmt.Errorf("unsupported message type: %d", messageType)
	}
	select {
	case c.send <- message:
		return nil
	case <-time.After(writeTimeout):
		return fmt.Errorf("failed to send message to client %s, connection timeout", c.id)
	default:
		// Notify hub to unregister the client
		c.hub.unregister <- c
		return fmt.Errorf("failed to send message to client %s, closing connection", c.id)
	}
}

// readPump handles reading messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	entry := getLogEntry()

	// Set read deadline and pong handler for keepalive
	c.conn.SetReadDeadline(time.Now().Add(readTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				entry.Error("WebSocket read error",
					zap.String("client_id", c.id),
					zap.Error(err))
			}
			break
		}

		// Process received message
		entry.Info("Received message from client",
			zap.String("client_id", c.id),
			zap.ByteString("message", message))

		// Process the message (e.g., broadcast it to other clients)

		select {
		case c.recv <- message:
		case <-time.After(idleTimeout):
			entry.Warn("Client idle timeout, skip message",
				zap.String("client_id", c.id),
				zap.ByteString("message", message))
		default:
			close(c.recv)
			delete(c.hub.clients, c)
		}
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *Client) writePump() {

	entry := getLogEntry()
	ticker := time.NewTicker(pingCycle)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		entry.Info("Writing to client",
			zap.String("client_id", c.id))
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if len(message) == 0 {
				entry.Warn("No message to send to client",
					zap.String("client_id", c.id))
				continue
			}

			// Determine message type: binary if contains non-UTF8, text otherwise
			actualMessage := []byte{}
			messageType := websocket.TextMessage
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))

			if len(message) > 2 && message[1] == 0xFF {
				switch message[0] {
				case websocket.TextMessage:
					messageType = websocket.TextMessage
					actualMessage = message[2:] // Remove the prefix
				case websocket.BinaryMessage:
					messageType = websocket.BinaryMessage
					actualMessage = message[2:] // Remove the prefix
				default:
					entry.Warn("Unknown message type, treating as text",
						zap.String("client_id", c.id),
						zap.ByteString("message", message))
					messageType = websocket.TextMessage
					actualMessage = message[1:] // Remove the prefix
				}
			}

			if err := c.conn.WriteMessage(messageType, actualMessage); err != nil {
				entry.Error("Failed to write message to client",
					zap.String("client_id", c.id),
					zap.Error(err))
			}

			// Log the message sent to the client
			entry = entry.With(zap.String("client_id", c.id), zap.ByteString("message", actualMessage))
			if messageType == websocket.BinaryMessage {
				entry = entry.With(zap.String("message",
					fmt.Sprintf("Binary message of length %d: %s", len(actualMessage), base64.StdEncoding.EncodeToString(actualMessage))))
			}
			// Process write message
			entry.Info("Sent to client")

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				entry.Error("Failed to send ping to client",
					zap.String("client_id", c.id),
					zap.Error(err))
				return
			}
		}
	}
}
