package infra

import (
	"context"
	"encoding/json"
	logger "ruantiengo/config/log"
	"ruantiengo/internal/domain"
	usecase "ruantiengo/internal/usecases"

	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type TransactionConsumer struct {
	channel            *amqp.Channel
	queueName          string
	processTransaction *usecase.ProcessTransaction
	processStats       *usecase.StatisticsService
}

func NewTransactionConsumer(channel *amqp.Channel, queueName string, processTransaction *usecase.ProcessTransaction, proccessStats *usecase.StatisticsService) *TransactionConsumer {
	return &TransactionConsumer{
		channel:            channel,
		queueName:          queueName,
		processTransaction: processTransaction,
		processStats:       proccessStats,
	}
}

func (c *TransactionConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",    // consumer
		false, // auto-ack
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
	logger.Message(logger.Debug, "Mensagem recebida: %s", msg.Body)
	var transactionDTO TransactionDTO
	err := json.Unmarshal(msg.Body, &transactionDTO)
	if err != nil {
		logger.Message(logger.Error, "Erro ao decodificar mensagem: %v", err)
		return
	}

	transaction := mapDTOToDomain(transactionDTO)

	// Save transaction in database
	err = c.saveTransaction(context.Background(), transaction)
	if err != nil {
		logger.Message(logger.Error, "Erro ao salvar transação: %v", err)
		return
	}

	// Update incremental statistics
	err = c.incrementStatistics(context.Background(), transaction)
	if err != nil {
		logger.Message(logger.Error, "Erro ao atualizar estatísticas: %v", err)
		return
	}
	msg.Ack(true)

	logger.Message(logger.Info, "Transação processada com sucesso: %s", transaction.BankSlipUuid)
}

func (c *TransactionConsumer) saveTransaction(ctx context.Context, transaction domain.Transaction) error {
	err := c.processTransaction.Execute(ctx, transaction)
	return err
}

func (c *TransactionConsumer) incrementStatistics(ctx context.Context, transaction domain.Transaction) error {

	err := c.processStats.Execute(ctx, transaction)

	if err != nil {
		return err
	}

	return err
}

func conditionalIncrement(status string, targetStatus string) int32 {
	if status == targetStatus {
		return 1
	}
	return 0
}

func conditionalValue(status string, targetStatus string, value float64) float64 {
	if status == targetStatus {
		return value
	}
	return 0
}

type TransactionDTO struct {
	BankSlipUuid           string  `json:"bank_slip_uuid"`
	Status                 string  `json:"status"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
	DueDate                string  `json:"due_date"`
	Total                  float64 `json:"total"`
	CustomerDocumentNumber string  `json:"customer_document_number"`
	TenantId               string  `json:"tenant_id"`
	BranchId               string  `json:"branch_id"`
	PaymentMethod          string  `json:"payment_method"`
}

func mapDTOToDomain(dto TransactionDTO) domain.Transaction {
	return domain.Transaction{
		BankSlipUuid:           uuid.MustParse(dto.BankSlipUuid),
		Status:                 domain.TransactionStatus(dto.Status),
		CreatedAt:              parseTime(dto.CreatedAt),
		UpdatedAt:              parseTime(dto.UpdatedAt),
		DueDate:                parseTime(dto.DueDate),
		Total:                  dto.Total,
		CustomerDocumentNumber: dto.CustomerDocumentNumber,
		TenantId:               dto.TenantId,
		BranchId:               dto.BranchId,
		PaymentMethod:          domain.PaymentMethod(dto.PaymentMethod),
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
