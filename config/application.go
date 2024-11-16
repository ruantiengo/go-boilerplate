package config

import (
	"database/sql"
	postgres_config "ruantiengo/config/database"
	server_config "ruantiengo/config/gim"
	logger "ruantiengo/config/log"
	rabbit_config "ruantiengo/config/rabbitmq"
	"ruantiengo/internal/transaction/infra"
	usecase "ruantiengo/internal/transaction/usecases"

	"github.com/gin-gonic/gin"
)

type Application struct {
	db                  *sql.DB
	rabbitMQManager     *rabbit_config.RabbitMQManager
	transactionConsumer *infra.TransactionConsumer
	router              *gin.Engine
	server              *server_config.Server
}

func NewApplication() *Application {
	CheckVariables()
	return &Application{
		router: gin.Default(),
	}
}

func (app *Application) Run() {
	if err := app.initDatabase(); err != nil {
		logger.Message(logger.Error, "Failed to initialize database: %v", err)
		return
	}

	if err := app.initRabbitMQ(); err != nil {
		logger.Message(logger.Error, "Failed to initialize RabbitMQ: %v", err)
	}

	if err := app.initTransactionConsumer(); err != nil {
		logger.Message(logger.Error, "Failed to initialize transaction consumer: %v", err)
		return
	}

	app.server = server_config.NewServer(app.router)
	app.server.SetupRoutes()
	app.server.Start()
}

func (app *Application) initDatabase() error {
	dbConfig := postgres_config.NewPostgresConfig()
	db, err := postgres_config.NewPostgresDB(dbConfig)
	if err != nil {
		return err
	}
	app.db = db
	logger.Message(logger.Info, "✔️ Database initialized successfully")
	return nil
}

func (app *Application) initRabbitMQ() error {
	rabbitMQConfig := rabbit_config.NewRabbitMQConfig()
	manager := rabbit_config.NewRabbitMQManager(rabbitMQConfig)

	if err := manager.Connect(); err != nil {
		return err
	}

	if err := manager.DeclareQueue(rabbit_config.TransactionQueue.String()); err != nil {
		return err
	}

	app.rabbitMQManager = manager
	logger.Message(logger.Info, "✔️ RabbitMQ connection established successfully")
	return nil
}

func (app *Application) initTransactionConsumer() error {
	if app.rabbitMQManager == nil {
		return nil
	}

	channel, err := app.rabbitMQManager.GetChannel()
	if err != nil {
		return err
	}

	transactionRepo := infra.NewPostgresTransactionRepository(app.db)
	transactionService := usecase.NewProcessTransaction(transactionRepo)
	consumer := infra.NewTransactionConsumer(channel, rabbit_config.TransactionQueue.String(), transactionService)

	if err := consumer.Start(); err != nil {
		return err
	}

	app.transactionConsumer = consumer
	logger.Message(logger.Info, "✔️ Transaction consumer started successfully")
	return nil
}

func (app *Application) Shutdown() {
	if app.db != nil {
		app.db.Close()
	}
	if app.rabbitMQManager != nil {
		app.rabbitMQManager.Close()
	}
	if app.server != nil {
		app.server.Shutdown()
	}
}
