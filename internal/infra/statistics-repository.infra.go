package infra

import (
	"context"
	"database/sql"
	"fmt"
	db "ruantiengo/database/generated"
	"ruantiengo/internal/domain"
	"ruantiengo/internal/repository"
	"time"
)

type postgresStatsRepository struct {
	queries *db.Queries
}

func NewPostgresStatsRepository(sqlDB *sql.DB) repository.StatsRepository {
	return &postgresStatsRepository{
		queries: db.New(sqlDB),
	}
}

func (r *postgresStatsRepository) UpdateBranchDailyStats(ctx context.Context, branchId string, date string, transaction domain.Transaction) error {
	params := db.UpdateBranchDailyStatsParams{
		TotalBoletos:      sql.NullInt32{Int32: 1, Valid: true},
		TotalPagos:        sql.NullInt32{Int32: 0, Valid: true},
		ValorEmitido:      sql.NullString{String: formatFloat(transaction.Total), Valid: true},
		ValorRecebido:     sql.NullString{String: "0", Valid: true},
		BoletosCancelados: sql.NullInt32{Int32: 0, Valid: true},
		BranchID:          branchId,
		Date:              transaction.CreatedAt,
	}

	switch transaction.Status {
	case domain.TransactionStatusApproved:
		params.TotalPagos.Int32 = 1
		params.ValorRecebido.String = formatFloat(transaction.Total)
	case domain.TransactionStatusCancelled:
		params.BoletosCancelados.Int32 = 1
	}

	return r.queries.UpdateBranchDailyStats(ctx, params)
}

func (r *postgresStatsRepository) UpdateCustomerMonthlyStats(ctx context.Context, customerDocumentNumber string, month string, transaction domain.Transaction) error {
	params := db.UpdateCustomerMonthlyStatsParams{
		CustomerDocumentNumber: customerDocumentNumber,
		Month:                  time.Date(transaction.CreatedAt.Year(), transaction.CreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC),
		TotalBoletos:           sql.NullInt32{Int32: 1, Valid: true},
		TotalPagos:             sql.NullInt32{Int32: 0, Valid: true},
		ValorEmitido:           sql.NullString{String: formatFloat(transaction.Total), Valid: true},
		ValorRecebido:          sql.NullString{String: "0", Valid: true},
	}

	if transaction.Status == domain.TransactionStatusApproved {
		params.TotalPagos.Int32 = 1
		params.ValorRecebido.String = formatFloat(transaction.Total)
	}

	return r.queries.UpdateCustomerMonthlyStats(ctx, params)
}

func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}
