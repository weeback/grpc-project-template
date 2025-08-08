package net

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		clients:    make(map[*Client]bool),
		once:       sync.Once{}, // No pending clients initially
	}

	// Start the hub in a goroutine
	go hub.run()
	return hub
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	once       sync.Once
}

// Run starts the hub and handles client registration/unregistration and broadcasting
func (h *Hub) run() {
	h.once.Do(func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
				log.Printf("Client %s connected. Total clients: %d", client.id, len(h.clients))

				// Send welcome message to new client
				// welcome := fmt.Sprintf(`{"type":"welcome","message":"Welcome client %s!","timestamp":"%s"}`,
				// 	client.id, time.Now().Format(time.RFC3339))
				// select {
				// case client.send <- []byte(welcome):
				// default:
				// 	close(client.send)
				// 	delete(h.clients, client)
				// }

			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
					close(client.recv)
					log.Printf("Client %s disconnected. Total clients: %d", client.id, len(h.clients))
				} else {
					log.Printf("Client %s not found in hub or already disconnected", client.id)
				}

			case message := <-h.broadcast:
				// Broadcast message to all connected clients
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
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

func (h *Hub) SendTo(receiveId string, message []byte) error {
	for client := range h.clients {
		if client.id == receiveId {
			return client.SendMessage(message)
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

func (c *Client) ReceiveMessage() ([]byte, error) {
	message, ok := <-c.recv
	if !ok {
		return nil, fmt.Errorf("client %s disconnected", c.id)
	}
	return message, nil
}

// func (c *Client) SendTo(receiveId string, message []byte) error {
// 	// Find the target client by ID
// 	for receive := range c.hub.clients {
// 		if receive.id == receiveId {
// 			select {
// 			case receive.send <- message:
// 				return nil
// 			case <-time.After(writeTimeout):
// 				return fmt.Errorf("failed to send message to client %s, connection timeout", receive.id)
// 			default:
// 				// Notify hub to unregister the client
// 				receive.hub.unregister <- receive
// 				return fmt.Errorf("failed to send message to client %s, closing connection", receive.id)
// 			}
// 		}
// 	}
// 	return fmt.Errorf("client %s not found", receiveId)
// }

// func (c *Client) BroadcastMessage(message []byte) error {
// 	select {
// 	case c.hub.broadcast <- message:
// 		return nil
// 	case <-time.After(writeTimeout):
// 		return fmt.Errorf("failed to broadcast message from client %s, closing connection", c.id)
// 	default:
// 		// Notify hub to unregister the client
// 		c.hub.unregister <- c
// 		return fmt.Errorf("failed to broadcast message from client %s, closing connection", c.id)
// 	}
// }

// readPump handles reading messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

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
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Process received message
		log.Printf("Received from client %s: %s", c.id, string(message))

		// Process the message (e.g., broadcast it to other clients)

		select {
		case c.recv <- message:
		case <-time.After(idleTimeout):
			log.Printf("Client %s idle timeout, skip message: %s\n", c.id, string(message))
		default:
			close(c.recv)
			delete(c.hub.clients, c)
		}

		// // Echo the message back to all clients with sender info
		// broadcastMsg := fmt.Sprintf(`{"type":"message","from":"%s","content":"%s","timestamp":"%s"}`,
		// 	c.id, string(message), time.Now().Format(time.RFC3339))
		// c.hub.broadcast <- []byte(broadcastMsg)
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingCycle)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if len(message) == 0 {
				continue // Skip empty messages
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// This block is commented out to avoid sending queued messages.
			// It only available if you want to implement message queuing.
			// On application level, you can manage message queuing if needed.
			// This can be done by buffering messages in the `send` channel.
			//
			// Add queued messages to the current message
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write([]byte{'\n'})
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}

			// Process write message
			log.Printf("Sent to client %s: %s", c.id, string(message))

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
