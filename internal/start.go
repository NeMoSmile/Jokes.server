package start

import (
	"fmt"
	"log"
	"net/http"

	d "github.com/NeMoSmile/Jokes.server.git/internal/data"
	h "github.com/NeMoSmile/Jokes.server.git/internal/handlers"
	"github.com/robfig/cron"
)

func Start(host string) {

	c := cron.New()

	// Запуск задачи каждый день в 00:00:00
	c.AddFunc("0 0 0 * * *", func() {
		d.NewDay()
	})
	c.Start()

	http.HandleFunc("/ws", h.MessageHandler)

	http.HandleFunc("/pagedata", h.PagedataHandler)

	http.HandleFunc("/check", h.CheckHandler)

	http.HandleFunc("/append", h.AppendHandler)

	http.HandleFunc("/wdata", h.WdataHandler)

	fmt.Println("Server listening on " + host)

	log.Fatal(http.ListenAndServe(host, nil))
	select {}
}
