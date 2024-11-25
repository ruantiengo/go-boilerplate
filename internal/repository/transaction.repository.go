package repository

import (
	"context"
	"ruantiengo/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction domain.Transaction) error
}
