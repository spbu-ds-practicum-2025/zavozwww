package email_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	email "authServ/pkg/emailSender"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestTemplates создает временные HTML шаблоны для тестирования
func createTestTemplates(t *testing.T) string {
	tempDir := t.TempDir()

	// Создаем тестовый шаблон для верификации email
	verificationTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Subject}}</title>
</head>
<body>
    <h1>Email Verification</h1>
    <p>Hello {{.Username}},</p>
    <p>Your verification code is: <strong>{{.Code}}</strong></p>
    <p>Year: {{.Year}}</p>
</body>
</html>`

	// Создаем тестовый шаблон для сброса пароля
	resetTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Subject}}</title>
</head>
<body>
    <h1>Password Reset</h1>
    <p>Hello {{.Username}},</p>
    <p>Click the link to reset your password: <a href="{{.ResetLink}}">Reset Password</a></p>
    <p>Year: {{.Year}}</p>
</body>
</html>`

	// Сохраняем шаблоны в файлы
	err := os.WriteFile(filepath.Join(tempDir, "verification.html"), []byte(verificationTemplate), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "reset.html"), []byte(resetTemplate), 0644)
	require.NoError(t, err)

	return tempDir
}

// === ТЕСТЫ ДЛЯ NewGomailSender ===

func TestNewGomailSender(t *testing.T) {
	t.Run("Успешная инициализация с валидными шаблонами", func(t *testing.T) {
		tempDir := createTestTemplates(t)

		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, tempDir)

		require.NoError(t, err)
		assert.NotNil(t, sender)
	})

	t.Run("Ошибка при отсутствии директории с шаблонами", func(t *testing.T) {
		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, "/nonexistent/directory")

		require.Error(t, err)
		assert.Nil(t, sender)
		assert.ErrorIs(t, err, email.ErrTemplateParsing)
	})

	t.Run("Ошибка при пустой директории с шаблонами", func(t *testing.T) {
		tempDir := t.TempDir() // Пустая директория

		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, tempDir)

		// ParseGlob возвращает ошибку, если не найдено ни одного файла по маске
		require.Error(t, err)
		assert.Nil(t, sender)
		assert.ErrorIs(t, err, email.ErrTemplateParsing)
	})
}

// === ТЕСТЫ ДЛЯ Send (Unit тесты без реального SMTP) ===

func TestGomailSender_Send_TemplateExecution(t *testing.T) {
	t.Run("Успешное выполнение шаблона с данными", func(t *testing.T) {
		tempDir := createTestTemplates(t)

		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, tempDir)
		require.NoError(t, err)

		ctx := context.Background()
		msg := email.Message{
			To:           []string{"user@example.com"},
			Subject:      "Email Verification",
			TemplateName: "verification.html",
			TemplateData: map[string]interface{}{
				"Username": "John Doe",
				"Code":     "123456",
			},
		}

		// Примечание: этот тест попытается отправить email, но упадет из-за недоступности SMTP
		// В реальных условиях нужно мокировать dialer
		err = sender.Send(ctx, msg)

		// Ожидаем ошибку подключения к SMTP, но не ошибку шаблона
		if err != nil {
			assert.NotErrorIs(t, err, email.ErrTemplateExecute)
			assert.Contains(t, err.Error(), "failed to send email")
		}
	})

	t.Run("Ошибка при несуществующем шаблоне", func(t *testing.T) {
		tempDir := createTestTemplates(t)

		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, tempDir)
		require.NoError(t, err)

		ctx := context.Background()
		msg := email.Message{
			To:           []string{"user@example.com"},
			Subject:      "Test Email",
			TemplateName: "nonexistent.html",
			TemplateData: map[string]interface{}{
				"Data": "test",
			},
		}

		err = sender.Send(ctx, msg)

		require.Error(t, err)
		assert.ErrorIs(t, err, email.ErrTemplateExecute)
	})

	t.Run("Успешное выполнение шаблона без дополнительных данных", func(t *testing.T) {
		tempDir := createTestTemplates(t)

		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, tempDir)
		require.NoError(t, err)

		ctx := context.Background()
		msg := email.Message{
			To:           []string{"user@example.com"},
			Subject:      "Password Reset",
			TemplateName: "reset.html",
			TemplateData: map[string]interface{}{
				"Username":  "Jane Doe",
				"ResetLink": "https://example.com/reset",
			},
		}

		err = sender.Send(ctx, msg)

		// Ожидаем ошибку SMTP, но не ошибку шаблона
		if err != nil {
			assert.NotErrorIs(t, err, email.ErrTemplateExecute)
		}
	})

	t.Run("Автоматическое добавление Year и Subject в TemplateData", func(t *testing.T) {
		tempDir := createTestTemplates(t)

		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, tempDir)
		require.NoError(t, err)

		ctx := context.Background()
		msg := email.Message{
			To:           []string{"user@example.com"},
			Subject:      "Test Subject",
			TemplateName: "verification.html",
			TemplateData: map[string]interface{}{
				"Username": "Test User",
				"Code":     "999999",
			},
		}

		// Вызываем Send - Year и Subject должны добавиться автоматически
		err = sender.Send(ctx, msg)

		// Проверяем, что нет ошибки выполнения шаблона (Year присутствует)
		if err != nil {
			assert.NotErrorIs(t, err, email.ErrTemplateExecute)
		}
	})

	t.Run("Обработка nil TemplateData", func(t *testing.T) {
		tempDir := createTestTemplates(t)

		config := email.SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		}

		sender, err := email.NewGomailSender(config, tempDir)
		require.NoError(t, err)

		ctx := context.Background()
		msg := email.Message{
			To:           []string{"user@example.com"},
			Subject:      "Test",
			TemplateName: "verification.html",
			TemplateData: nil, // nil данные
		}

		err = sender.Send(ctx, msg)

		if err != nil {
			assert.NotErrorIs(t, err, email.ErrTemplateExecute)
		}
	})
}

func TestMessage_Structure(t *testing.T) {
	t.Run("Создание сообщения с полными данными", func(t *testing.T) {
		msg := email.Message{
			To:           []string{"user1@example.com", "user2@example.com"},
			Subject:      "Test Subject",
			TemplateName: "verification.html",
			TemplateData: map[string]interface{}{
				"Key1": "Value1",
				"Key2": 123,
			},
		}

		assert.Equal(t, 2, len(msg.To))
		assert.Equal(t, "Test Subject", msg.Subject)
		assert.Equal(t, "verification.html", msg.TemplateName)
		assert.Equal(t, "Value1", msg.TemplateData["Key1"])
		assert.Equal(t, 123, msg.TemplateData["Key2"])
	})

	t.Run("Создание сообщения с минимальными данными", func(t *testing.T) {
		msg := email.Message{
			To:           []string{"user@example.com"},
			Subject:      "Minimal",
			TemplateName: "template.html",
		}

		assert.NotNil(t, msg.To)
		assert.Nil(t, msg.TemplateData)
	})
}

// === ТЕСТЫ ДЛЯ Errors ===

func TestErrors(t *testing.T) {
	t.Run("ErrTemplateParsing определен", func(t *testing.T) {
		assert.NotNil(t, email.ErrTemplateParsing)
		assert.Equal(t, "failed to parse email templates", email.ErrTemplateParsing.Error())
	})

	t.Run("ErrTemplateExecute определен", func(t *testing.T) {
		assert.NotNil(t, email.ErrTemplateExecute)
		assert.Equal(t, "failed to execute email template", email.ErrTemplateExecute.Error())
	})
}
