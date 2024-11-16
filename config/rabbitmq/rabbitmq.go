package rabbit_config

import (
	"fmt"
	"os"
	"time"

	logger "ruantiengo/config/log"

	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Uri string
}

type RabbitMQManager struct {
	Config *RabbitMQConfig
	Conn   *amqp.Connection
	Ch     *amqp.Channel
}

func NewRabbitMQConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		Uri: os.Getenv("RABBITMQ_URI"),
	}
}

func (c *RabbitMQConfig) ConnectionString() string {
	return fmt.Sprintf("%s", c.Uri)
}

func NewRabbitMQManager(config *RabbitMQConfig) *RabbitMQManager {
	return &RabbitMQManager{
		Config: config,
	}
}

func (m *RabbitMQManager) Connect() error {
	var err error
	m.Conn, err = amqp.Dial(m.Config.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	m.Ch, err = m.Conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	go m.handleReconnect()

	return nil
}

func (m *RabbitMQManager) handleReconnect() {
	for {
		reason, ok := <-m.Conn.NotifyClose(make(chan *amqp.Error))
		if !ok {
			logger.Message(logger.Warning, "RabbitMQ connection closed. Attempting to reconnect...")
			break
		}
		logger.Message(logger.Error, "RabbitMQ connection closed. Reason: %v", reason)

		for {
			time.Sleep(5 * time.Second)

			err := m.Connect()
			if err == nil {
				logger.Message(logger.Info, "Successfully reconnected to RabbitMQ")
				break
			}

			logger.Message(logger.Error, "Failed to reconnect to RabbitMQ: %v", err)
		}
	}
}

func (m *RabbitMQManager) GetChannel() (*amqp.Channel, error) {
	if m.Ch == nil {
		return nil, fmt.Errorf("RabbitMQ channel is not initialized")
	}
	return m.Ch, nil
}

func (m *RabbitMQManager) Close() {
	if m.Ch != nil {
		m.Ch.Close()
	}
	if m.Conn != nil {
		m.Conn.Close()
	}
}
