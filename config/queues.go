package config

import "fmt"

func (m *RabbitMQManager) DeclareQueue(queueName string) error {
	ch, err := m.GetChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %v", err)
	}

	// Declare the queue
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
