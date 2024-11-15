package infra

import (
	"context"
	"database/sql"
	db "ruantiengo/database/generated"
	"ruantiengo/internal/transaction/domain"
	"ruantiengo/internal/transaction/repository"

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

func (r *postgresTransactionRepository) Save(ctx context.Context, transaction domain.Transaction) error {
	params := db.CreateTransactionParams{
		BankSlipUuid:  transaction.BankSlipUuid,
		Status:        db.NullTransactionStatus{Valid: true, TransactionStatus: db.TransactionStatus(transaction.Status)},
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		PaymentMethod: db.NullPaymentMethod{Valid: true, PaymentMethod: db.PaymentMethod(transaction.PaymentMethod)},
	}

	err := r.queries.CreateTransaction(ctx, params)
	return err
}

func (r *postgresTransactionRepository) GetByBankSlipUuid(ctx context.Context, bankSlipUuid string) (domain.Transaction, error) {
	uuid, err := uuid.Parse(bankSlipUuid)
	if err != nil {
		return domain.Transaction{}, err
	}

	dbTransaction, err := r.queries.GetTransactionByUUID(ctx, uuid)
	if err != nil {
		return domain.Transaction{}, err
	}

	return domain.Transaction{
		BankSlipUuid: dbTransaction.BankSlipUuid,
		Status:       domain.TransactionStatus(dbTransaction.Status.TransactionStatus),
		CreatedAt:    dbTransaction.CreatedAt,
		UpdatedAt:    dbTransaction.UpdatedAt,
		PaymentMethod: func() domain.PaymentMethod {
			if dbTransaction.PaymentMethod.Valid {
				return domain.PaymentMethod(dbTransaction.PaymentMethod.PaymentMethod)
			}
			return domain.PaymentMethod("")
		}(),
	}, nil
}
