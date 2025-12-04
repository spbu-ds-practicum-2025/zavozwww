package main

import (
	rout "authServ/internal/handler/http/router"
	config "authServ/internal/handler/http/server"
	repositories "authServ/internal/repositories/postgres"
	"authServ/internal/services"
	email "authServ/pkg/emailSender"
	"authServ/pkg/logger"
	"context"
	"fmt"
	http "net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Helper function для чтения env (чтобы не загромождать main)
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("FATAL PANIC DETECTED BEFORE LOGGER INIT: %v\n", r)
			os.Exit(1)
		}
	}()
	fmt.Println("DEBUG: Starting main function")

	fmt.Printf("DEBUG: APP_CONFIG_PATH=%s\n", os.Getenv("APP_CONFIG_PATH"))

	ctx := context.Background()
	fmt.Println("DEBUG: Context created")

	fmt.Println("DEBUG: Before logger initialization")
	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Println("Error initializing logger:", "error", err)
		return 
	}
	logger.Info("Logger initialized successfully.")

	var cfgPg repositories.PgConfig
	err = cfgPg.GetConfig()
	if err != nil {
		logger.Error("Error getting PG config:", "error", err)
		return
	}

	dsnMigrate := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfgPg.User, cfgPg.Password, cfgPg.Host, cfgPg.Port, cfgPg.DBName, cfgPg.SSLMode)

	const migrationsPath = "file:///app/migration"

	m, err := migrate.New(
		migrationsPath,
		dsnMigrate,
	)
	if err != nil {
		logger.Error("Error creating migrator", "err", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Failed to apply migrations", "err", err)
		return
	}
	logger.Info("Migrations applied successfully.")

	pgxPool, err := repositories.NewPostgresStorage(ctx, cfgPg)
	if err != nil {
		logger.Error("Error connecting to Postgres:", "error", err)
		return
	}
	defer pgxPool.Close()
	logger.Info("Successfully connected to Postgres")
	userRepo := repositories.NewUserRepository(pgxPool)

	refreshTokenRepo := repositories.NewRefreshTokenRepository(pgxPool)
	profileRepo := repositories.NewUserProfileRepository(pgxPool)

	smtpPortStr := getEnv("SMTP_PORT", "1025") // Используем 1025 для Mailpit
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		logger.Error("Invalid SMTP port:", "error", err)
		return
	}

	smtpCfg := email.SMTPConfig{
		Host:     getEnv("SMTP_HOST", "mailpit"), // Используем 'mailpit' для имени сервиса
		Port:     smtpPort,
		Username: getEnv("SMTP_USER", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", "MyService <no-reply@example.com>"),
	}

	templatesDir := getEnv("EMAIL_TEMPLATES_DIR", "./templates")

	emailSender, err := email.NewGomailSender(smtpCfg, templatesDir)
	if err != nil {
		logger.Error("FAILED to initialize Email Sender", "error", err)
		return
	}

	logger.Info("Email sender initialized successfully")
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		logger.Error("SECRET_KEY environment variable is not set")
		return
	}

	accessTTL := os.Getenv("ACCESS_TOKEN_TTL")
	if accessTTL == "" {
		logger.Error("ACCESS_TOKEN_TTL environment variable is not set")
		return
	}
	accessTTLInt, err := strconv.Atoi(accessTTL)
	if err != nil {
		logger.Error("Invalid ACCESS_TOKEN_TTL:", "error", err)
		return
	}
	userservice := services.NewUserService(userRepo, refreshTokenRepo, emailSender, logger, secretKey, time.Duration(accessTTLInt)*time.Second)
	profileservice := services.NewProfileService(profileRepo)

	logger.Info("Services initialized successfully")

	handler := rout.NewHandler(userservice, profileservice, logger)
	router := handler.InitRoutes()

	serverConfig := config.MustLoad()
	serverHost := serverConfig.Server.Host
	serverPort := serverConfig.Server.Port

	logger.Info(fmt.Sprintf("Starting server on %s:%s", serverHost, serverPort))
	if err := http.ListenAndServe(serverHost+":"+serverPort, router); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}
