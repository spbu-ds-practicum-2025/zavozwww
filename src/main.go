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
)

// Helper function для чтения env (чтобы не загромождать main)
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	ctx := context.Background()
	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	fmt.Println("Logger level:", logger.Level())

	var cfgPg repositories.PgConfig
	err = cfgPg.GetConfig()
	if err != nil {
		fmt.Println("Error getting PG config:", err)
		return
	}

	pgxPool, err := repositories.NewPostgresStorage(ctx, cfgPg)
	if err != nil {
		fmt.Println("Error connecting to Postgres:", err)
		return
	}
	defer pgxPool.Close()
	fmt.Println("Successfully connected to Postgres")

	userRepo := repositories.NewUserRepository(pgxPool)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(pgxPool) // Пока закомментировал, чтобы линтер не ругался если не используется
	profileRepo := repositories.NewUserProfileRepository(pgxPool)

	smtpPortStr := getEnv("SMTP_PORT", "587")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		fmt.Println("Invalid SMTP port:", err)
		return
	}

	smtpCfg := email.SMTPConfig{
		Host:     getEnv("SMTP_HOST", "smtp.example.com"),
		Port:     smtpPort,
		Username: getEnv("SMTP_USER", "user@example.com"),
		Password: getEnv("SMTP_PASSWORD", "secret"),
		From:     getEnv("SMTP_FROM", "MyService <no-reply@example.com>"),
	}

	templatesDir := getEnv("EMAIL_TEMPLATES_DIR", "./templates")

	emailSender, err := email.NewGomailSender(smtpCfg, templatesDir)
	if err != nil {
		fmt.Printf("FAILED to initialize Email Sender: %v\n", err)
		return
	}

	fmt.Println("Email sender initialized successfully")

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("SECRET_KEY environment variable is not set")
		return
	}

	accessTTL := os.Getenv("ACCESS_TOKEN_TTL")
	if accessTTL == "" {
		fmt.Println("ACCESS_TOKEN_TTL environment variable is not set")
		return
	}
	accessTTLInt, err := strconv.Atoi(accessTTL)
	if err != nil {
		fmt.Println("Invalid ACCESS_TOKEN_TTL:", err)
		return
	}
	userservice := services.NewUserService(userRepo, refreshTokenRepo, emailSender, logger, secretKey, time.Duration(accessTTLInt)*time.Second)
	profileservice := services.NewProfileService(profileRepo)

	fmt.Println("Services initialized successfully")

	handler := rout.NewHandler(userservice, profileservice, logger)
	router := handler.InitRoutes()

	serverConfig := config.MustLoad()
	serverHost := serverConfig.Server.Host
	serverPort := serverConfig.Server.Port

	fmt.Printf("Starting server on %s:%s\n", serverHost, serverPort)
	if err := http.ListenAndServe(serverHost+":"+serverPort, router); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
