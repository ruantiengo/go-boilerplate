package main

import (
	"context"
	"ruantiengo/config"
	"ruantiengo/internal/transaction/repository"
	logger "ruantiengo/log"
	"time"
)

var rabbitMQManager *config.RabbitMQManager

func main() {
	config.New()

	db, err := initDatabase()
	if err != nil {
		logger.Message(logger.Error, "Erro ao inicializar o banco de dados: %v", err)
		return
	}
	logger.Message(logger.Info, "✔️ Banco de dados inicializado com sucesso")

	transactionRepo := repository.TransactionRepository()

	rabbitMQManager = initRabbitMQ()

	defer rabbitMQManager.Close()
	defer db.Close()
}

func initDatabase() (*config.PostgresDB, error) {
	dbConfig := config.NewPostgresConfig()
	db, err := config.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Message(logger.Error, "Erro ao inicializar o banco de dados: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.TestConnection(ctx); err != nil {
		logger.Message(logger.Error, "Erro ao testar conexão com o banco de dados: %v", err)
		return nil, err
	}

	return db, nil
}

func initRabbitMQ() *config.RabbitMQManager {
	rabbitMQConfig := config.NewRabbitMQConfig()
	manager := config.NewRabbitMQManager(rabbitMQConfig)

	err := manager.Connect()
	if err != nil {
		logger.Message(logger.Error, "Erro ao conectar ao RabbitMQ: %v", err)
		logger.Message(logger.Info, "A aplicação continuará funcionando sem o RabbitMQ")
	} else {
		logger.Message(logger.Info, "✔️ Conexão com RabbitMQ estabelecida com sucesso")
	}

	return manager
}
