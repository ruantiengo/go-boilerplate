package main

import (
	"database/sql"
	"ruantiengo/config"
	"ruantiengo/internal/transaction/infra"
	usecase "ruantiengo/internal/transaction/usecases"
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
	rabbitMQManager = initRabbitMQ()

	transactionRepo := infra.NewPostgresTransactionRepository(db)
	transactionService := usecase.NewProcessTransaction(transactionRepo)
	if rabbitMQManager != nil {
		channel, err := rabbitMQManager.GetChannel()
		if err != nil {
			logger.Message(logger.Error, "Erro ao obter canal do RabbitMQ: %v", err)
		} else {
			transactionConsumer := infra.NewTransactionConsumer(channel, "transaction_queue", transactionService)
			err = transactionConsumer.Start()
			if err != nil {
				logger.Message(logger.Error, "Erro ao iniciar consumidor de transações: %v", err)
			} else {
				logger.Message(logger.Info, "✔️ Consumidor de transações iniciado com sucesso")
			}
		}
	}

	rabbitMQManager = initRabbitMQ()
	time.Sleep(10000 * time.Second)
	defer rabbitMQManager.Close()
	defer db.Close()
}

func initDatabase() (*sql.DB, error) {
	dbConfig := config.NewPostgresConfig()
	db, err := config.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Message(logger.Error, "Erro ao inicializar o banco de dados: %v", err)
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
		return nil
	}
	logger.Message(logger.Info, "✔️ Conexão com RabbitMQ estabelecida com sucesso")

	err = manager.DeclareQueue("transaction_queue")
	if err != nil {
		logger.Message(logger.Error, "Failed to declare queue: %v", err)
		return nil
	}
	return manager
}
