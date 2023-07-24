package confirm

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	d "github.com/NeMoSmile/Jokes.server.git/internal/data"
)

var db = d.Db

func Send(email string) {
	code := fmt.Sprintf("%06d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(999999))

	sendMail(email, code)

	err := insertOrUpdateUser(email, code)
	if err != nil {
		fmt.Println(err)
	}
}

func Check(email, code string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM codes WHERE email = $1 AND code = $2", email, code).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if count > 0 {
		return true
	}
	return false
}

func CheckUser(id string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	if count > 0 {
		return true
	}
	return false
}

func GetId(email string) string {
	var id string
	err := db.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&id)
	if err != nil {
		fmt.Println(err)
	}
	return id
}

func insertOrUpdateUser(email string, code string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM codes WHERE email = $1", email).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		_, err = db.Exec("UPDATE codes SET code = $1 WHERE email = $2", code, email)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec("INSERT INTO codes (email, code) VALUES ($1, $2)", email, code)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendMail(email, code string) {
	// Настройте параметры аутентификации SMTP
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := "jokes.com342@gmail.com"
	senderPassword := "jumgugznztymsdas"

	// Настройте получателя, тему и содержимое письма
	recipientEmail := email
	subject := "confirmation code"
	body := "<h1>Your code: " + code + "</h1>"

	// Форматируем письмо в соответствии с требованиями HTML
	message := fmt.Sprintf("From: %s\r\n", senderEmail)
	message += fmt.Sprintf("To: %s\r\n", recipientEmail)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message += fmt.Sprintf("%s\r\n", body)

	// Настроить аутентификацию SMTP
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Отправить письмо через клиента SMTP
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, senderEmail, []string{recipientEmail}, []byte(message))
	if err != nil {
		fmt.Println("error sending mail: ", err)
		return
	}

}
