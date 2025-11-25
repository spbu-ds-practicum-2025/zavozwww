package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"gopkg.in/gomail.v2"
)

type gomailSender struct {
	dialer    *gomail.Dialer
	from      string
	templates *template.Template
}

// SMTPConfig содержит конфигурацию для подключения к SMTP-серверу.
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// NewGomailSender теперь также парсит все шаблоны из указанной директории.
func NewGomailSender(config SMTPConfig, templatesDir string) (EmailSender, error) {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	templatesBody, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTemplateParsing, err)
	}

	return &gomailSender{
		dialer:    dialer,
		from:      config.From,
		templates: templatesBody,
	}, nil
}

// Send отправляет письмо, рендеря шаблон с переданными данными.
func (s *gomailSender) Send(ctx context.Context, msg Message) error {
	if msg.TemplateData == nil {
		msg.TemplateData = make(map[string]interface{})
	}
	msg.TemplateData["Year"] = time.Now().Year()
	msg.TemplateData["Subject"] = msg.Subject

	var body bytes.Buffer
	err := s.templates.ExecuteTemplate(&body, msg.TemplateName, msg.TemplateData)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTemplateExecute, err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/html", body.String())

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// EmailSender определяет контракт для отправки электронных писем.
type EmailSender interface {
	Send(ctx context.Context, msg Message) error
}
