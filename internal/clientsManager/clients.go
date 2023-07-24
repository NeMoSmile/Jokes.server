package clientsmanager

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
}

type SendingMessage struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Outgoing bool   `json:"outgoing"`
	Err      bool   `json:"err"`
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

func (cm *ClientsManager) BroadcastMessage(conn *websocket.Conn, text, username string, erro bool) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	sendingAll := SendingMessage{
		Text:     text,
		Username: username,
		Outgoing: false,
		Err:      false,
	}

	sending := SendingMessage{
		Text:     text,
		Username: username,
		Outgoing: true,
		Err:      erro,
	}
	if erro {
		err := conn.WriteJSON(sending)
		if err != nil {
			log.Println("Error writing to client:", err)
		}
	} else {
		for client := range cm.clients {
			if client.Conn == conn {
				err := client.Conn.WriteJSON(sending)
				if err != nil {
					log.Println("Error writing to client:", err)
					cm.RemoveClient(client)
				}
			} else {
				err := client.Conn.WriteJSON(sendingAll)
				if err != nil {
					log.Println("Error writing to client:", err)
					cm.RemoveClient(client)
				}
			}
		}
	}
}
