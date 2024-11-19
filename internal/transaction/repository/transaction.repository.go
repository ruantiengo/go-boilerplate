package repository

import (
	"context"
	"ruantiengo/internal/transaction/domain"
)

type TransactionRepository interface {
	Upsert(ctx context.Context, transaction domain.Transaction) error
}
