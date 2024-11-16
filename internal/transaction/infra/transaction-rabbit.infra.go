package infra

import (
	"context"
	"encoding/json"
	logger "ruantiengo/config/log"
	"ruantiengo/internal/transaction/domain"
	usecase "ruantiengo/internal/transaction/usecases"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type TransactionConsumer struct {
	channel            *amqp.Channel
	queueName          string
	processTransaction *usecase.ProcessTransaction
}

func NewTransactionConsumer(channel *amqp.Channel, queueName string, processTransaction *usecase.ProcessTransaction) *TransactionConsumer {
	return &TransactionConsumer{
		channel:            channel,
		queueName:          queueName,
		processTransaction: processTransaction,
	}
}

func (c *TransactionConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			c.processMessage(msg)
		}
	}()

	return nil
}

func (c *TransactionConsumer) processMessage(msg amqp.Delivery) {
	var transactionDTO TransactionDTO
	err := json.Unmarshal(msg.Body, &transactionDTO)
	if err != nil {
		logger.Message(logger.Error, "Erro ao decodificar mensagem: %v", err)
		return
	}

	transaction := mapDTOToDomain(transactionDTO)

	err = c.processTransaction.Execute(context.Background(), transaction)
	if err != nil {
		logger.Message(logger.Error, "Erro ao processar transação: %v", err)
		// Implement retry logic or dead-letter queue here
		return
	}

	logger.Message(logger.Info, "Transação processada com sucesso: %s", transaction.BankSlipUuid)
}

type TransactionDTO struct {
	BankSlipUuid  string `json:"bank_slip_uuid"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PaymentMethod string `json:"payment_method"`
}

func mapDTOToDomain(dto TransactionDTO) domain.Transaction {
	return domain.Transaction{
		BankSlipUuid:  uuid.MustParse(dto.BankSlipUuid),
		Status:        domain.TransactionStatus(dto.Status),
		CreatedAt:     parseTime(dto.CreatedAt),
		UpdatedAt:     parseTime(dto.UpdatedAt),
		PaymentMethod: domain.PaymentMethod(dto.PaymentMethod),
	}
}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		logger.Message(logger.Warning, "Erro ao fazer parse do tempo: %v", err)
		return time.Time{}
	}
	return t
}
