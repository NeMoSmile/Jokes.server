package start

import (
	"fmt"
	"log"
	"net/http"

	commands "github.com/NeMoSmile/Jokes.server.git/internal/commands"
	d "github.com/NeMoSmile/Jokes.server.git/internal/data"
	h "github.com/NeMoSmile/Jokes.server.git/internal/handlers"
	"github.com/robfig/cron"
)

func Start(host string) {

	go commands.Listen()

	c := cron.New()

	// The server clears the data every day
	c.AddFunc("0 0 0 * * *", func() {
		d.NewDay()
	})
	c.Start()

	http.HandleFunc("/", h.No)

	// websocket connection
	http.HandleFunc("/ws", h.MessageHandler)

	// Loading information on the page
	http.HandleFunc("/pagedata", h.PagedataHandler)

	// Checking if the user has entered their data correctly. If password and email are correct sent 1 if password is not correct 2 if user does not exist 3
	http.HandleFunc("/check", h.CheckHandler)

	// Adding a new user
	http.HandleFunc("/append", h.AppendHandler)

	// tagged user jokes
	http.HandleFunc("/wdata", h.WdataHandler)

	// sending confirmation code
	http.HandleFunc("/send", h.SendHandler)

	// check email verification code
	http.HandleFunc("/confirm", h.ConfirmHandler)

	// user existence check
	http.HandleFunc("/checkuser", h.CheckUserHandler)

	// getting user id
	http.HandleFunc("/getid", h.GetIdHandler)

	http.HandleFunc("/getmess", h.GetMess)

	fmt.Println("Server listening on " + host)

	log.Fatal(http.ListenAndServe(host, nil))
	select {}
}
