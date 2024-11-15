package config

import (
	"os"
	logger "ruantiengo/log"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func New() *Config {
	CheckVariables()
	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
			Mode: os.Getenv("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Name:     os.Getenv("POSTGRES_DB"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
	}

	logger.Message(logger.Info, "✔️ Configurações carregadas com sucesso")
	return config
}
