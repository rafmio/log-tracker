// https://mimepost.com/blog/sending-an-email-with-golang/
package main

import (
	"log"
	"net/smtp"
)

func main() {

	// Setup host information
	host := "smtp.yandex.ru"
	port := "587"

	// Setup headers
	to := []string{"rafmio@yandex.ru"}
	msg := []byte("To: rafmio@yandex.ru\r\n" +
		"Subject: hello Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	from := "rafmio@yandex.ru"

	// Set up authentication information.
	username := "rafmio"
	password := "qwerty"

	auth := smtp.PlainAuth("", username, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
