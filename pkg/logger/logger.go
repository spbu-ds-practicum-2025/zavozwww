package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Logger представляет собой настраиваемый логгер с уровнем логирования и путем к файлу.
type Logger struct {
	level string `yaml:"level"`
	slog.Logger
}

// configFile представляет собой структуру конфигурационного файла для логгера.
type configFile struct {
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
}

// NewLogger инициализирует и возвращает новый экземпляр Logger на основе конфигурационного файла.
func NewLogger() (*Logger, error) {
	const op = "pkg.logger.NewLogger"
	var cfgFile configFile

	cfgPath := os.Getenv("APP_CONFIG_PATH")
	if cfgPath == "" {
		return nil, fmt.Errorf("%s: environment variable APP_CONFIG_PATH is not set", op)
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfgFile); err != nil {
		return nil, fmt.Errorf("%s: read config %q: %w", op, cfgPath, err)
	}

	cfg := Logger{
		level: cfgFile.Logging.Level,
	}

	out := os.Stdout

	var handler slog.Handler
	if cfgFile.Logging.Format == "JSON" {
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level:     getLevel(cfg.level),
			AddSource: false,
		})
	} else {
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{
			Level:     getLevel(cfg.level),
			AddSource: false,
		})
	}

	slogLogger := slog.New(handler)

	l := &Logger{
		level:  cfg.level,
		Logger: *slogLogger,
	}

	return l, nil
}

// getLevel преобразует строковое представление уровня логирования в slog.Level.
func getLevel(level string) slog.Level {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Level возвращает уровень логирования.
func (l *Logger) Level() string {
	return l.level
}
