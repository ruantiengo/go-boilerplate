package usecase

import (
	"context"
	"ruantiengo/internal/transaction/domain"
)

type ProcessTransaction struct {
	repo TransactionRepository
}

func NewProcessTransaction(repo TransactionRepository) *ProcessTransaction {
	return &ProcessTransaction{repo: repo}
}

func (pt *ProcessTransaction) Execute(ctx context.Context, transaction domain.Transaction) error {
	return pt.repo.Save(ctx, transaction)
}

type TransactionRepository interface {
	Save(ctx context.Context, transaction domain.Transaction) error
}
