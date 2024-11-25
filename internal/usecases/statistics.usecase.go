package usecase

import (
	"context"

	"ruantiengo/internal/domain"
	"ruantiengo/internal/repository"
)

type UpdateStatistics struct {
	statsRepo repository.StatsRepository
}

func NewUpdateStatistics(statsRepo repository.StatsRepository) *UpdateStatistics {
	return &UpdateStatistics{
		statsRepo: statsRepo,
	}
}

func (u *UpdateStatistics) Execute(ctx context.Context, transaction domain.Transaction) error {
	date := transaction.CreatedAt.Format("2006-01-02")
	month := transaction.CreatedAt.Format("2006-01")

	err := u.statsRepo.UpdateBranchDailyStats(ctx, transaction.BranchId, date, transaction)
	if err != nil {
		return err
	}

	err = u.statsRepo.UpdateCustomerMonthlyStats(ctx, transaction.CustomerDocumentNumber, month, transaction)
	if err != nil {
		return err
	}

	return nil
}
