package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type PostgresDB struct {
	Config *PostgresConfig
	Pool   *pgxpool.Pool
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

func NewPostgresDB(config *PostgresConfig) (*PostgresDB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Database)

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %v", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 5 * time.Minute
	poolConfig.MaxConnIdleTime = 15 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err)
	}

	return &PostgresDB{
		Config: config,
		Pool:   pool,
	}, nil
}

func (db *PostgresDB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *PostgresDB) TestConnection(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("connection pool is not initialized")
	}
	return db.Pool.Ping(ctx)
}
