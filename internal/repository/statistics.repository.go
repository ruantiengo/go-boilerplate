package repository

import (
	"context"
	"ruantiengo/internal/domain"
)

type StatsRepository interface {
	UpdateBranchDailyStats(ctx context.Context, branchId string, date string, transaction domain.Transaction) error
	UpdateCustomerMonthlyStats(ctx context.Context, customerId string, month string, transaction domain.Transaction) error
}
