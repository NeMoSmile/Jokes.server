package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn *websocket.Conn
}

type ClientsManager struct {
	clients map[*Client]bool
	mutex   sync.Mutex
}

func NewClientsManager() *ClientsManager {
	return &ClientsManager{
		clients: make(map[*Client]bool),
	}
}

func (cm *ClientsManager) AddClient(client *Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[client] = true
}

func (cm *ClientsManager) RemoveClient(client *Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.clients, client)
}

func (cm *ClientsManager) BroadcastMessage(message []byte) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for client := range cm.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error writing to client:", err)
			cm.RemoveClient(client)
		}
	}
}

func main() {
	clientsManager := NewClientsManager()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}

		client := &Client{
			conn: conn,
		}

		clientsManager.AddClient(client)

		// Создаем горуту для чтения сообщений от данного клиента
		go func() {
			defer func() {
				conn.Close()
				clientsManager.RemoveClient(client)
			}()

			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("Error reading message:", err)
					break
				}

				// Обработка полученного сообщения от клиента
				// ...

				// Пример широковещательной отправки сообщений всем клиентам
				clientsManager.BroadcastMessage(message)
			}
		}()
	})

	http.HandleFunc("/pagedata", func(w http.ResponseWriter, r *http.Request) {
		// Обработка обычного POST запроса по /pagedata
		// ...
	})

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
