package rabbit_config

import "fmt"

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

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	return nil
}
