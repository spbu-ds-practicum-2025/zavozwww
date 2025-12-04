package email

import (
	"errors"
)

var (
	ErrTemplateParsing = errors.New("failed to parse email templates")
	ErrTemplateExecute = errors.New("failed to execute email template")
)

type Message struct {
	To           []string
	Subject      string
	TemplateName string
	TemplateData map[string]interface{}
}
