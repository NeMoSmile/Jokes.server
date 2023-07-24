package data

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 8082
	user     = "postgres"
	password = "yS72_3w*90P"
	dbname   = "postgres"
)

type PData struct {
	FirstPl  string `json:"FirstPl"`
	SecondPl string `json:"SecondPl"`
	ThirdPl  string `json:"ThirdPl"`
	MyTitle  string `json:"MyTitle"`
	MyText1  string `json:"MyText1"`
	MyText2  string `json:"MyText2"`
	Id       string `json:"Id"`
}

var Db *sql.DB

func init() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error

	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	// _, err = Db.Exec("DROP TABLE IF EXISTS users CASCADE")
	// if err != nil {
	// 	log.Println(err)
	// }
	// _, err = Db.Exec("DROP TABLE IF EXISTS jokes CASCADE")
	// if err != nil {
	// 	log.Println(err)
	// }
	// _, err = Db.Exec("DROP TABLE IF EXISTS codes CASCADE")
	// if err != nil {
	// 	log.Println(err)
	// }

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  w TEXT[],
  today VARCHAR(255) NOT NULL,
  month INTEGER,
  last VARCHAR(255) NOT NULL
 )`)
	if err != nil {
		log.Println(err)
	}

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS jokes (
		id VARCHAR(255) NOT NULL,
  text TEXT NOT NULL
 )`)
	if err != nil {
		log.Println(err)
	}

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS codes (
		email VARCHAR(255) NOT NULL,
  code VARCHAR(255) NOT NULL
 )`)
	if err != nil {
		log.Println(err)
	}

	// _, err = Db.Exec(`DELETE FROM users`)
	// if err != nil {
	// 	log.Println(err)
	// }

	// _, err = Db.Exec(`DELETE FROM jokes`)
	// if err != nil {
	// 	log.Println(err)
	// }

	// _, err = Db.Exec(`DELETE FROM codes`)
	// if err != nil {
	// 	log.Println(err)
	// }

	fmt.Println("Database connected")

}

func PageData(id string) PData {

	first, second, third := best()
	m := GetName(id)
	today, month := me(id)

	return PData{
		FirstPl:  first,
		SecondPl: second,
		ThirdPl:  third,
		MyTitle:  m,
		MyText1:  today,
		MyText2:  month,
		Id:       id,
	}
}

func Check(email, pass string) int {
	var count int
	err := Db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count == 0 {
		return 3
	}

	var storedPassword string
	err = Db.QueryRow(`SELECT password FROM users WHERE email = $1`, email).Scan(&storedPassword)
	if err != nil {
		log.Println(err)
	}

	if storedPassword == pass {
		return 1
	} else {
		return 2
	}
}

func generateRandomID(db *sql.DB) (string, error) {
	for {
		randomID := fmt.Sprintf("%015d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999999999999999))

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", randomID).Scan(&count)
		if err != nil {
			return "", err
		}

		if count == 0 {
			return randomID, nil
		}
	}
}

func Append(email, pass, username string) string {
	w := []string{}

	id, err := generateRandomID(Db)
	if err != nil {
		fmt.Println(err)
	}
	_, err = Db.Exec("INSERT INTO users (id, email, password, username, w, today, month, last) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		id, email, pass, username, pq.Array(w), 0, 0, "00:00:00")
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println("New user registered: " + username)
	return id

}

func Wdata(id string) []string {
	var res []string
	err := Db.QueryRow(`SELECT w FROM users WHERE id = $1`, id).Scan(pq.Array(&res))
	if err != nil {
		log.Println(err)
	}

	return res
}

func AddJoke(id, joke string) {
	_, err := Db.Exec(`INSERT INTO jokes (id, text) VALUES ($1, $2)`, id, joke)
	if err != nil {
		log.Println(err)
		return
	}
	now := time.Now().Format("15:04:05")
	_, err = Db.Exec(`UPDATE users SET last = $1 WHERE id = $2`, now, id)
	if err != nil {
		log.Println(err)
		return
	}

}

func AddWJoke(id, joke string) {
	var w []string
	err := Db.QueryRow("SELECT w FROM users WHERE id = $1", id).Scan(pq.Array(&w))
	if err != nil {
		log.Println(err)
		return
	}

	w = append(w, joke)

	_, err = Db.Exec("UPDATE users SET w = $1 WHERE id = $2", pq.Array(w), id)
	if err != nil {
		log.Println(err)

		return
	}

	var user string
	err = Db.QueryRow("SELECT id FROM jokes WHERE text = $1", joke).Scan(&user)
	if err != nil {
		log.Println(err)

		return
	}

	var today int
	err = Db.QueryRow("SELECT today FROM users WHERE id = $1", user).Scan(&today)
	if err != nil {
		log.Println(err)
		return
	}

	today += 1

	_, err = Db.Exec("UPDATE users SET today = $1 WHERE id = $2", today, user)
	if err != nil {
		log.Println(err)
		return
	}

	var month int
	err = Db.QueryRow("SELECT month FROM users WHERE id = $1", user).Scan(&month)
	if err != nil {
		log.Println(err)
		return
	}

	month += 1

	_, err = Db.Exec("UPDATE users SET month = $1 WHERE id = $2", month, user)
	if err != nil {
		log.Println(err)
		return
	}

}

func GetName(id string) string {
	var name string
	err := Db.QueryRow(`SELECT username FROM users WHERE id = $1`, id).Scan(&name)
	if err != nil {
		log.Println(err)
	}
	return name
}

func CheckJoke(id, joke string) string {
	if len(joke) > 3000 {
		return "Write jokes no longer than 3000 characters"
	}
	var res string
	err := Db.QueryRow(`SELECT last FROM users WHERE id = $1`, id).Scan(&res)
	if err != nil {
		log.Println(err)
	}
	now := time.Now().Format("15:04:05")
	resT, err := time.Parse("15:04:05", res)
	if err != nil {
		fmt.Println("Time parsing error1:", err)
	}
	nowT, err := time.Parse("15:04:05", now)
	if err != nil {
		fmt.Println("Time parsing error2:", err)
	}

	dif := nowT.Sub(resT)

	if dif < time.Hour && res != "00:00:00" {
		return "It's been less than an hour since your last joke."
	}

	var count int
	err = Db.QueryRow(`SELECT COUNT(*) FROM jokes WHERE text = $1`, joke).Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count != 0 {
		return "This joke was already written today"
	}

	return ""
}

func best() (string, string, string) {
	var anser []string
	rows, err := Db.Query(`SELECT username, today FROM users ORDER BY today DESC LIMIT 3`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var today int

		err := rows.Scan(&username, &today)
		if err != nil {
			log.Println(err)
		}

		anser = append(anser, username+" "+strconv.Itoa(today))
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	if len(anser) > 2 {
		return anser[0], anser[1], anser[2]
	}
	if len(anser) == 2 {
		return anser[0], anser[1], "-"
	}
	if len(anser) == 1 {
		return anser[0], "-", "-"
	}
	return "-", "-", "-"

}

func me(id string) (string, string) {
	var today string = "Today: "
	var month string = "Month: "
	var todayS int
	var monthS int

	err := Db.QueryRow(`SELECT today FROM users WHERE id = $1`, id).Scan(&todayS)
	if err != nil {
		log.Println(err)
	}
	err = Db.QueryRow(`SELECT month FROM users WHERE id = $1`, id).Scan(&monthS)
	if err != nil {
		log.Println(err)
	}

	today += strconv.Itoa(todayS) + " #"
	month += strconv.Itoa(monthS) + " #"

	err = Db.QueryRow(`SELECT COUNT(*) FROM users WHERE today >= (SELECT today FROM users WHERE id = $1)`, id).Scan(&todayS)
	if err != nil {
		log.Println(err)
	}
	err = Db.QueryRow(`SELECT COUNT(*) FROM users WHERE month >= (SELECT month FROM users WHERE id = $1)`, id).Scan(&monthS)
	if err != nil {
		log.Println(err)
	}
	today += strconv.Itoa(todayS)
	month += strconv.Itoa(monthS)

	return today, month

}

func NewDay() {
	_, err := Db.Exec(`UPDATE users SET today = 0, last = '00:00:00'`)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("A new day has begun")

	_, err = Db.Exec(`DELETE FROM jokes`)
	if err != nil {
		log.Println(err)
	}

}
