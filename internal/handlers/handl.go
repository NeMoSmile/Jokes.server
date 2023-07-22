package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	c "github.com/NeMoSmile/Jokes.server.git/internal/clientsManager"
	d "github.com/NeMoSmile/Jokes.server.git/internal/data"
	"github.com/gorilla/websocket"
)

type Message struct {
	Text     string `json:"message"`
	Email    string `json:"email"`
	Outgoing bool   `json:"outgoing"`
}

var clientsManager = c.NewClientsManager()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func PagedataHandler(w http.ResponseWriter, r *http.Request) {
	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(d.PageData(message["email"]))
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	fmt.Println("user connected")

	client := &c.Client{
		Conn: conn,
	}

	clientsManager.AddClient(client)

	// Создаем горуту для чтения сообщений от данного клиента
	go func() {
		defer func() {
			conn.Close()
			clientsManager.RemoveClient(client)
		}()

		for {
			_, mess, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}

			var message Message

			err = json.Unmarshal(mess, &message)

			if err != nil {
				log.Println("Error unmarshaling message:", err)
				break
			}

			if message.Outgoing {
				erro := d.CheckJoke(message.Email, message.Text)
				if erro == "" {
					clientsManager.BroadcastMessage(conn, message.Text, message.Email, false)
					d.AddJoke(message.Email, message.Text)
				} else {
					clientsManager.BroadcastMessage(conn, erro, "", true)
				}
			} else {
				d.AddWJoke(message.Email, message.Text)
			}
		}
	}()
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {

	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(d.Check(message["email"], message["pass"]))
}

func AppendHandler(w http.ResponseWriter, r *http.Request) {
	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	d.Append(message["email"], message["pass"], message["name"])
}

func WdataHandler(w http.ResponseWriter, r *http.Request) {
	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(d.Wdata(message["email"]))
}
