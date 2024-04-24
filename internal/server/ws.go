package server

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/kulikvl/sharea/internal/storage"
	"log"
)

func (s *Server) setupWebSockets() {
	var clients = make(map[*websocket.Conn]bool)
	var storageChange = make(chan string)

	s.App.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	s.App.Get("/ws/files", websocket.New(func(c *websocket.Conn) {
		defer func() {
			log.Println("Bye Bye client: ", c.IP())
			err := c.Close()
			if err != nil {
				log.Println("Failed to close connection with the client: ", c.IP())
			}
		}()

		log.Println("New WS client: ", c.IP())
		clients[c] = true
		storageChange <- "new client!"

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
					log.Println("WebSocket closed normally for client: ", c.IP())
				} else {
					log.Println("Failed to read message from the client:", err)
				}
				delete(clients, c)
				break
			}
		}
	}))

	go s.Storage.Watch(storageChange)
	go broadcast(storageChange, clients, &s.Storage)
}

// broadcast broadcasts storage changes to all active clients.
func broadcast(changes <-chan string, clients map[*websocket.Conn]bool, storage *storage.Storage) {
	for {
		<-changes

		msg, err := createMessage(storage)
		if err != nil {
			log.Println("Failed to create broadcast message, but will try again next storage change")
			continue
		}

		for client := range clients {
			if err := sendMessage(msg, client); err != nil {
				log.Println("Failed to send broadcast message to one of the clients, delete that client")
				delete(clients, client)
			}
		}
	}
}

// sendMessage sends storage files info message to the client.
func sendMessage(msg []byte, conn *websocket.Conn) error {
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
		_ = conn.Close()
		return fmt.Errorf("failed to send storage files info message to client: %w", err)
	}

	return nil
}

// createMessage creates storage files info message in JSON format.
func createMessage(storage *storage.Storage) ([]byte, error) {
	files, err := storage.GetFilesInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage files info: %w", err)
	}

	msg, err := json.Marshal(files)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal storage files info to json: %w", err)
	}

	return msg, nil
}
