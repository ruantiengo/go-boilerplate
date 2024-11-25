package infra

import (
	"context"
	"database/sql"
	"fmt"
	db "ruantiengo/database/generated"
	"ruantiengo/internal/domain"
	"ruantiengo/internal/repository"

	"github.com/google/uuid"
)

type postgresTransactionRepository struct {
	queries *db.Queries
}

func NewPostgresTransactionRepository(sqlDB *sql.DB) repository.TransactionRepository {
	return &postgresTransactionRepository{
		queries: db.New(sqlDB),
	}
}

func (r *postgresTransactionRepository) Create(ctx context.Context, transaction domain.Transaction) error {
	params := db.CreateTransactionParams{
		BankSlipUuid:           uuid.NullUUID{UUID: transaction.BankSlipUuid, Valid: true},
		Status:                 db.NullTransactionStatus{Valid: true, TransactionStatus: db.TransactionStatus(transaction.Status)},
		CreatedAt:              transaction.CreatedAt,
		UpdatedAt:              transaction.UpdatedAt,
		DueDate:                transaction.DueDate,
		Total:                  fmt.Sprintf("%f", transaction.Total),
		CustomerDocumentNumber: transaction.CustomerDocumentNumber,
		TenantID:               transaction.TenantId,
		BranchID:               transaction.BranchId,
		PaymentMethod:          db.NullPaymentMethod{Valid: true, PaymentMethod: db.PaymentMethod(transaction.PaymentMethod)},
	}

	_, err := r.queries.CreateTransaction(ctx, params)
	return err
}
