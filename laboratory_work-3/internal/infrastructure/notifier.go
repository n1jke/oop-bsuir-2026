package infrastructure

import "fmt"

// SmtpMailer - имитация почтового сервиса.
type SMTPMailer struct {
	server string
}

func NewSMTPMailer(svr string) *SMTPMailer {
	return &SMTPMailer{server: svr}
}

func (s *SMTPMailer) Notify(to, subject, body string) {
	fmt.Printf(">> Connecting to SMTP server %s...\n", s.server)
	fmt.Printf(">> Sending EMAIL to %s\n   Subject: %s\n   Body: %s\n", to, subject, body)
}

// TelegramMailer - имитация бота в телеграмм дл отправки сообщений менеджеру.
type TelegramMailer struct {
	connString string
}

func NewTelegramMailer(conn string) *TelegramMailer {
	return &TelegramMailer{connString: conn}
}

func (t *TelegramMailer) Notify(to, subject, body string) {
	fmt.Printf(">> Connecting to Telegram bot %s...\n", t.connString)
	fmt.Printf(">> Sending MESSAGE to %s\n   Subject: %s\n   Body: %s\n", to, subject, body)
}
