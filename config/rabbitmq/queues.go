package rabbit_config

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Queue string

const (
	TransactionQueue Queue = "transaction_queue"
)

func (q Queue) String() string {
	return string(q)
}

func (m *RabbitMQManager) DeclareQueue(queueName string) error {
	ch, err := m.GetChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %v", err)
	}

	err = ch.Qos(
		150,   // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	_, err = ch.QueueDeclare(
		queueName,    // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		amqp.Table{}, // arguments with no ack
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	return nil
}
