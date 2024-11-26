package usecase

import (
	"context"
	"time"

	"ruantiengo/internal/domain"
	"ruantiengo/internal/repository"
)

type StatisticsService struct {
	statsRepo repository.StatsRepository
}

func NewUpdateStatistics(statsRepo repository.StatsRepository) *StatisticsService {
	return &StatisticsService{
		statsRepo: statsRepo,
	}
}

func (u *StatisticsService) Execute(ctx context.Context, transaction domain.Transaction) error {
	month := transaction.CreatedAt.Format("2006-01")

	err := u.statsRepo.UpdateBranchDailyStats(ctx, transaction.BranchId, transaction.CreatedAt, transaction)
	if err != nil {
		return err
	}

	err = u.statsRepo.UpdateCustomerMonthlyStats(ctx, transaction.CustomerDocumentNumber, month, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *StatisticsService) GetCompanyStatistics(ctx context.Context, tenantID, branchID string, startDate, endDate time.Time) (repository.CompanyStatisticsResult, error) {
	return s.statsRepo.CalculateCompanyStatistics(ctx, tenantID, branchID, startDate, endDate)
}

// func (s *StatisticsService) GetCustomerStatistics(ctx context.Context, tenantID, customerID string, startDate, endDate time.Time) (repository.CustomerStatisticsResult, error) {
// 	return s.statsRepo.CalculateCustomerStatistics(ctx, tenantID, customerID, startDate, endDate)
// }
