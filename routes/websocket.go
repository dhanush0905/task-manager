package routes

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// A map to track connected clients
var clients = make(map[*websocket.Conn]bool)
var mutex = sync.Mutex{} // To manage concurrent writes

// Setup WebSocket Routes
func SetupWebSocketRoutes(app *fiber.App) {
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Add client to map
		mutex.Lock()
		clients[c] = true
		mutex.Unlock()

		defer func() {
			mutex.Lock()
			delete(clients, c)
			mutex.Unlock()
			c.Close()
		}()

		for {
			// Read message (optional, we mainly send updates)
			_, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			fmt.Println("Received message:", string(msg))
		}
	}))
}

// Send updates to all connected WebSocket clients
func BroadcastUpdate(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("Error sending WebSocket message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
