package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"

	d "github.com/NeMoSmile/Jokes.server.git/internal/data"
	"github.com/lib/pq"
)

func Listen() {
	for {
		reader := bufio.NewReader(os.Stdin)

		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		str = str[:len(str)-1]
		switch str {
		case "delete":
			_, err = d.Db.Exec("DROP TABLE IF EXISTS users CASCADE")
			if err != nil {
				log.Println(err)
			}
			_, err = d.Db.Exec("DROP TABLE IF EXISTS jokes CASCADE")
			if err != nil {
				log.Println(err)
			}
			_, err = d.Db.Exec("DROP TABLE IF EXISTS codes CASCADE")
			if err != nil {
				log.Println(err)
			}

		case "clear all":
			_, err = d.Db.Exec(`DELETE FROM users`)
			if err != nil {
				log.Println(err)
			}

			_, err = d.Db.Exec(`DELETE FROM jokes`)
			if err != nil {
				log.Println(err)
			}

			_, err = d.Db.Exec(`DELETE FROM codes`)
			if err != nil {
				log.Println(err)
			}
		case "clear":
			_, err = d.Db.Exec(`DELETE FROM jokes`)
			if err != nil {
				log.Println(err)
			}

			_, err = d.Db.Exec("UPDATE users SET w = $1", pq.Array([]string{}))
			if err != nil {
				fmt.Println("Ошибка выполнения запроса:", err)
				return
			}

			_, err = d.Db.Exec(`DELETE FROM codes`)
			if err != nil {
				log.Println(err)
			}
		case "new month":
			_, err := d.Db.Exec(`UPDATE users SET month = 0`)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = d.Db.Exec("UPDATE users SET w = $1", pq.Array([]string{}))
			if err != nil {
				fmt.Println("Ошибка выполнения запроса:", err)
				return
			}

		}
	}

}
