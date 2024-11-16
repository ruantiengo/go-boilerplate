package postgres_config

import (
	"database/sql"
	"fmt"
	"os"

	logger "ruantiengo/log"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
}

func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

func NewPostgresDB(config *PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	logger.Message(logger.Info, "✔️ Conexão com o PostgreSQL estabelecida com sucesso.")
	return db, nil
}
