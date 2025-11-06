package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"gopkg.in/gomail.v2"
)

// gomailSender теперь хранит кэш распарсенных шаблонов.
type gomailSender struct {
	dialer    *gomail.Dialer
	from      string
	templates *template.Template
}

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

	// ParseGlob найдет все файлы *.html в директории и распарсит их.
	// Это делается один раз при старте приложения для эффективности.
	templates, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTemplateParsing, err)
	}

	return &gomailSender{
		dialer:    dialer,
		from:      config.From,
		templates: templates,
	}, nil
}

// Send теперь находит нужный шаблон, исполняет его и отправляет результат.
func (s *gomailSender) Send(ctx context.Context, msg Message) error {
	// Добавляем в данные для шаблона текущий год и тему.
	if msg.TemplateData == nil {
		msg.TemplateData = make(map[string]interface{})
	}
	msg.TemplateData["Year"] = time.Now().Year()
	msg.TemplateData["Subject"] = msg.Subject

	// Используем bytes.Buffer для "рендеринга" шаблона в память.
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

	// Отправляем письмо.
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
