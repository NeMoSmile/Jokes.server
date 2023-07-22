package data

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
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
	Email    string `json:"Email"`
}

var db *sql.DB

func init() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	// Подключаемся к базе данных
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	// Создаем таблицу "users"
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  w TEXT[],
  today INTEGER,
  month INTEGER,
  last TIME
 )`)
	if err != nil {
		log.Println(err)
	}

	// Создаем таблицу "jokes"
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS jokes (
		email VARCHAR(255) NOT NULL,
  text TEXT NOT NULL
 )`)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Таблицы успешно созданы!")

}

func PageData(email string) PData {

	first, second, third := best()
	m := getName(email)
	today, month := me(email)

	return PData{
		FirstPl:  first,
		SecondPl: second,
		ThirdPl:  third,
		MyTitle:  m,
		MyText1:  today,
		MyText2:  month,
		Email:    email,
	}
}

func Check(email, pass string) int {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count == 0 {
		return 3
	}

	var storedPassword string
	err = db.QueryRow(`SELECT password FROM users WHERE email = $1`, email).Scan(&storedPassword)
	if err != nil {
		log.Println(err)
	}

	if storedPassword == pass {
		return 1
	} else {
		return 2
	}
}

func Append(email, pass, username string) {
	_, err := db.Exec("INSERT INTO users (email, password, username, w, today, month, last) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		email, pass, username, []string{}, 0, 0, "00:00:00")
	if err != nil {
		log.Println(err)
		return
	}

}

func Wdata(email string) []string {
	var res []string
	err := db.QueryRow(`SELECT w FROM users WHERE email = $1`, email).Scan(&res)
	if err != nil {
		log.Println(err)
	}
	return res
}

func AddJoke(email, joke string) {
	_, err := db.Exec(`INSERT INTO jokes (email, text) VALUES ($1, $2)`, email, joke)
	if err != nil {
		log.Println(err)
		return
	}
	now := time.Now().Format("15:04:05")
	_, err = db.Exec(`UPDATE users SET last = $1 WHERE email = $2`, now, email)
	if err != nil {
		log.Println(err)
		return
	}

}

func AddWJoke(email, joke string) {
	var w []string
	err := db.QueryRow("SELECT w FROM users WHERE email = $1", email).Scan(&w)
	if err != nil {
		log.Println(err)
		return
	}

	// Добавить текст в массив "w"
	w = append(w, joke)

	// Обновить запись в базе данных с новым массивом "w"
	_, err = db.Exec("UPDATE users SET w = $1 WHERE email = $2", w, email)
	if err != nil {
		log.Println(err)

		return
	}

	var user string
	err = db.QueryRow("SELECT email FROM jokes WHERE text = $1", joke).Scan(&user)
	if err != nil {
		log.Println(err)

		return
	}

	var today int
	err = db.QueryRow("SELECT today FROM users WHERE email = $1", user).Scan(&today)
	if err != nil {
		log.Println(err)
		return
	}

	today += 1

	_, err = db.Exec("UPDATE users SET today = $1 WHERE email = $2", today, email)
	if err != nil {
		log.Println(err)
		return
	}

	var month int
	err = db.QueryRow("SELECT month FROM users WHERE email = $1", user).Scan(&month)
	if err != nil {
		log.Println(err)
		return
	}

	month += 1

	_, err = db.Exec("UPDATE users SET month = $1 WHERE email = $2", month, email)
	if err != nil {
		log.Println(err)
		return
	}

}

func getName(email string) string {
	var name string
	err := db.QueryRow(`SELECT username FROM users WHERE email = $1`, email).Scan(&name)
	if err != nil {
		log.Println(err)
	}
	return name
}

func CheckJoke(email, joke string) string {
	if len(joke) > 3000 {
		return "Write jokes no longer than 3000 characters"
	}
	var res string
	err := db.QueryRow(`SELECT last FROM users WHERE email = $1`, email).Scan(&res)
	if err != nil {
		log.Println(err)
	}
	timeFormat := "15:04:05"
	time2Str := time.Now().Format("15:04:05")

	// Преобразование строк в значения time.Time
	time1, _ := time.Parse(timeFormat, res)
	time2, _ := time.Parse(timeFormat, time2Str)

	// Вычисление разницы между временами
	diff := time2.Sub(time1)

	if int(diff.Hours()) < 1 && res != "00:00:00" {
		return "It's been less than an hour since your last joke."
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM jokes WHERE joke = $1`, joke).Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count < 1 {
		return "This joke was already written today"
	}

	return ""
}

func best() (string, string, string) {
	var anser []string
	rows, err := db.Query(`SELECT username FROM users ORDER BY today DESC LIMIT 3`)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var username string

		err := rows.Scan(&username)
		if err != nil {
			log.Println(err)
		}

		anser = append(anser, username)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	return anser[0], anser[1], anser[2]

}

func me(email string) (string, string) {
	var today string = "Today: "
	var month string = "Month: "
	var todayS int
	var monthS int

	err := db.QueryRow(`SELECT today FROM users WHERE email = $1`, email).Scan(&todayS)
	if err != nil {
		log.Println(err)
	}
	err = db.QueryRow(`SELECT month FROM users WHERE email = $1`, email).Scan(&monthS)
	if err != nil {
		log.Println(err)
	}

	today += strconv.Itoa(todayS) + " #"
	month += strconv.Itoa(monthS) + " #"

	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE today >= (SELECT today FROM users WHERE email = $1)`, email).Scan(&todayS)
	if err != nil {
		log.Println(err)
	}
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE month >= (SELECT month FROM users WHERE email = $1)`, email).Scan(&monthS)
	if err != nil {
		log.Println(err)
	}
	today += strconv.Itoa(todayS)
	month += strconv.Itoa(monthS)

	return today, month

}

func NewDay() {
	_, err := db.Exec(`UPDATE users SET today = 0, last = '00:00:00'`)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("A new day has begun")

	_, err = db.Exec(`DELETE FROM jokes`)
	if err != nil {
		log.Println(err)
	}

}