package infra

import (
	"context"
	db "ruantiengo/database/generated"
	"ruantiengo/internal/transaction/domain"
	"ruantiengo/internal/transaction/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTransactionRepository struct {
	queries *db.Queries
}

func NewPostgresTransactionRepository(dbPool *pgxpool.Pool) repository.TransactionRepository {
	return &postgresTransactionRepository{
		queries: db.New(dbPool),
	}
}

func (r *postgresTransactionRepository) Save(ctx context.Context, transaction domain.Transaction) error {
	params := db.CreateTransactionParams{
		BankSlipUuid:  transaction.BankSlipUuid,
		Status:        db.TransactionStatus(transaction.Status),
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		PaymentMethod: db.PaymentMethod(transaction.PaymentMethod),
	}

	_, err := r.queries.CreateTransaction(ctx, params)
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
		BankSlipUuid:  dbTransaction.BankSlipUuid,
		Status:        domain.TransactionStatus(dbTransaction.Status),
		CreatedAt:     dbTransaction.CreatedAt,
		UpdatedAt:     dbTransaction.UpdatedAt,
		PaymentMethod: domain.PaymentMethod(dbTransaction.PaymentMethod),
	}, nil
}
