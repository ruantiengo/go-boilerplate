package repository

import (
	"context"
	"ruantiengo/internal/transaction/domain"
)

type TransactionRepository interface {
	Save(ctx context.Context, transaction domain.Transaction) error
}
