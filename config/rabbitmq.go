package config

import (
	"fmt"
	"os"
	"time"

	logger "ruantiengo/log"

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
	conn, err := amqp.Dial(m.Config.ConnectionString())
	if err != nil {
		return err
	}
	m.Conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	m.Ch = ch

	go m.handleReconnect()

	return nil
}

func (m *RabbitMQManager) handleReconnect() {
	for {
		reason, ok := <-m.Conn.NotifyClose(make(chan *amqp.Error))
		if !ok {
			logger.Message(logger.Warning, "Conexão com RabbitMQ foi fechada. Tentando reconectar...")
			break
		}
		logger.Message(logger.Error, "Conexão com RabbitMQ perdida. Razão: %v", reason)

		for {
			time.Sleep(3 * time.Second)

			err := m.Connect()
			if err == nil {
				logger.Message(logger.Info, "Reconectado ao RabbitMQ com sucesso")
				break
			}

			logger.Message(logger.Error, "Falha ao reconectar com RabbitMQ: %v", err)
		}
	}
}

func (m *RabbitMQManager) Close() {
	if m.Ch != nil {
		m.Ch.Close()
	}
	if m.Conn != nil {
		m.Conn.Close()
	}
}

func (m *RabbitMQManager) GetChannel() *amqp.Channel {
	return m.Ch
}
