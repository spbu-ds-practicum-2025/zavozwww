package email

import (
	"context"
	"errors"
)

var (
	ErrTemplateParsing = errors.New("failed to parse email templates")
	ErrTemplateExecute = errors.New("failed to execute email template")
)

// Message инкапсулирует все данные для отправки одного письма.
type Message struct {
	To           []string
	Subject      string
	TemplateName string
	TemplateData map[string]interface{}
}

// EmailSender определяет контракт для отправки электронных писем.
type EmailSender interface {
	Send(ctx context.Context, msg Message) error
}
