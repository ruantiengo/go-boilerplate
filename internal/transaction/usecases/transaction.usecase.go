package usecase

import (
	"context"
	"ruantiengo/internal/transaction/domain"
	"ruantiengo/internal/transaction/repository"
)

type ProcessTransaction struct {
	repo repository.TransactionRepository
}

func NewProcessTransaction(repo repository.TransactionRepository) *ProcessTransaction {
	return &ProcessTransaction{repo: repo}
}

func (pt *ProcessTransaction) Execute(ctx context.Context, transaction domain.Transaction) error {
	return pt.repo.Upsert(ctx, transaction)
}
