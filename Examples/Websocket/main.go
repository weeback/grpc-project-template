package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github/weeback/grpc-project-template/pkg"
	"github/weeback/grpc-project-template/pkg/net"
)

func main() {

	pkg.Import()

	// Set up HTTP routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWS)

	// Start the server
	port := ":8080"
	fmt.Printf("WebSocket server starting on http://localhost%s\n", port)
	fmt.Println("Open your browser and navigate to http://localhost:8080 to test the WebSocket connection")

	log.Fatal(http.ListenAndServe(port, nil))
}

// serveWS handles WebSocket requests from clients
func serveWS(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket

	ch, mh, err := net.UpgradeToWebSocket(w, r)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	// Send welcome message to the new client
	welcomeMessage := fmt.Sprintf("Welcome to the WebSocket server! You are connected as %s", mh.ID())
	if err := mh.SendMessage([]byte(welcomeMessage)); err != nil {
		log.Printf("Failed to send welcome message: %v", err)
		return
	}

	log.Printf("New client connected: %s", mh.ID())

	// Start reading messages from the client
	for {
		message, err := mh.ReceiveMessage()
		if err != nil {
			log.Printf("Error receiving message from client %s: %v", mh.ID(), err)
			break
		}

		// Process the received message
		// | | |
		// v v v

		if strings.HasPrefix(string(message), "/ALL") || strings.HasPrefix(string(message), "/all") {
			parts := bytes.SplitN(message, []byte(" "), 2) // Split the message into command and content
			// Broadcast the received message to all clients
			if err := ch.BroadcastMessage(fmt.Appendf([]byte("From "), "%s: %s", mh.ID(), parts[1])); err != nil {
				log.Printf("Error broadcasting message from client %s: %v", mh.ID(), err)
				break
			}
			// Skip echoing the message back to the client
			continue
		}

		if strings.HasPrefix(string(message), "/CID-") {
			parts := bytes.SplitN(message, []byte(" "), 2) // Split the message into command and content
			receive := strings.Replace(string(parts[0]), "/", "", 1)
			content := bytes.Replace(message, fmt.Appendf(parts[0], " "), []byte{}, 1)
			log.Printf("From %s to %s: %s", mh.ID(), receive, string(content))
			// Send the message to the specified client
			if err := ch.SendTo(receive, fmt.Appendf([]byte("From "), "%s: %s", mh.ID(), content)); err != nil {
				log.Printf("Error sending message to client %s: %v", receive, err)
				break
			}
		}

		// Echo the message back to the client
		if err := mh.SendMessage(fmt.Appendf([]byte("You: "), "%s", message)); err != nil {
			log.Printf("Error sending message to client %s: %v", mh.ID(), err)
			break
		}
	}
	log.Printf("Client %s disconnected", mh.ID())
	// Clean up resources if needed
	if err := mh.Close(); err != nil {
		log.Printf("Error closing connection for client %s: %v", mh.ID(), err)
	}
}

// serveHome serves the HTML page with WebSocket client
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Test</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        #messages { border: 1px solid #ccc; height: 300px; overflow-y: scroll; padding: 10px; margin: 10px 0; }
        #messageInput { width: 300px; padding: 5px; }
        button { padding: 5px 10px; margin-left: 5px; }
        .message { margin: 5px 0; }
        .welcome { color: green; }
        .system { color: blue; }
    </style>
</head>
<body>
    <h1>WebSocket Test</h1>
    <div id="messages"></div>
    <input type="text" id="messageInput" placeholder="Type a message..." />
    <button onclick="sendMessage()">Send</button>
    <button onclick="disconnect()">Disconnect</button>

    <script>
        let ws;
        const messages = document.getElementById('messages');
        const messageInput = document.getElementById('messageInput');

        function connect() {
            ws = new WebSocket('ws://localhost:8080/ws');
            
            ws.onopen = function(event) {
                addMessage('Connected to server', 'system');
            };
            
            ws.onmessage = function(event) {
                try {
                    const data = JSON.parse(event.data);
                    let className = data.type === 'welcome' ? 'welcome' : 'message';
                    let displayMsg = '';
                    
                    if (data.type === 'welcome') {
                        displayMsg = data.message;
                    } else if (data.type === 'message') {
                        displayMsg = '[' + data.timestamp + '] ' + data.from + ': ' + data.content;
                    }
                    
                    addMessage(displayMsg, className);
                } catch (e) {
                    addMessage(event.data, 'message');
                }
            };
            
            ws.onclose = function(event) {
                addMessage('Disconnected from server', 'system');
            };
            
            ws.onerror = function(error) {
                addMessage('WebSocket error: ' + error, 'system');
            };
        }

        function addMessage(message, className) {
            const div = document.createElement('div');
            div.className = 'message ' + className;
            div.textContent = message;
            messages.appendChild(div);
            messages.scrollTop = messages.scrollHeight;
        }

        function sendMessage() {
            if (ws && ws.readyState === WebSocket.OPEN) {
                const message = messageInput.value.trim();
                if (message) {
                    ws.send(message);
                    messageInput.value = '';
                }
            } else {
                addMessage('Not connected to server', 'system');
            }
        }

        function disconnect() {
            if (ws) {
                ws.close();
            }
        }

        // Connect when page loads
        connect();

        // Send message on Enter key
        messageInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                sendMessage();
            }
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
