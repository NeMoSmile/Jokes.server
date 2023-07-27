package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	c "github.com/NeMoSmile/Jokes.server.git/internal/clientsManager"
	d "github.com/NeMoSmile/Jokes.server.git/internal/data"

	confirm "github.com/NeMoSmile/Jokes.server.git/internal/confirm"
	"github.com/gorilla/websocket"
)

type Message struct {
	Text     string `json:"message"`
	Id       string `json:"id"`
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
	json.NewEncoder(w).Encode(d.PageData(message["id"]))
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
				erro := d.CheckJoke(message.Id, message.Text)
				if erro == "" {
					clientsManager.BroadcastMessage(conn, message.Text, d.GetName(message.Id), false)
					d.AddJoke(message.Id, message.Text)
				} else {
					clientsManager.BroadcastMessage(conn, erro, "", true)
				}
			} else {
				d.AddWJoke(message.Id, message.Text)
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
	json.NewEncoder(w).Encode(d.Append(message["email"], message["pass"], message["name"]))
}

func WdataHandler(w http.ResponseWriter, r *http.Request) {
	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(d.Wdata(message["id"]))
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	confirm.Send(message["email"])
}

func ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(confirm.Check(message["email"], message["code"]))
}

func CheckUserHandler(w http.ResponseWriter, r *http.Request) {

	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(confirm.CheckUser(message["id"]))
}

func GetIdHandler(w http.ResponseWriter, r *http.Request) {

	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(confirm.GetId(message["email"]))
}

func GetMess(w http.ResponseWriter, r *http.Request) {

	var message map[string]string

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(d.GetMess(message["id"]))
}

func No(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Привет!")
}
