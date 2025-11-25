package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Server struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"8080"`
}

type ServerConfig struct {
	Server Server `yaml:"server"`
}

func MustLoad() *ServerConfig {
	configPath := os.Getenv("APP_CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg ServerConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
