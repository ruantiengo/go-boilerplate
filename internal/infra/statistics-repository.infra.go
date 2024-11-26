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

// UpdateBranchDailyStats implements repository.StatsRepository.
func (r *postgresStatsRepository) UpdateBranchDailyStats(ctx context.Context, branchId string, date time.Time, transaction domain.Transaction) error {
	params := db.UpdateBranchDailyStatsParams{
		TotalBoletos:      sql.NullInt32{Int32: 1, Valid: true},
		TotalPagos:        sql.NullInt32{Int32: 0, Valid: true},
		ValorEmitido:      sql.NullFloat64{Float64: transaction.Total, Valid: true},
		ValorRecebido:     sql.NullFloat64{Float64: 0, Valid: true},
		BoletosCancelados: sql.NullInt32{Int32: 0, Valid: true},
		BranchID:          branchId,
		Date:              transaction.CreatedAt,
		TenantID:          transaction.TenantId,
	}

	switch transaction.Status {
	case domain.TransactionStatusApproved:
		params.TotalPagos.Int32 = 1
		params.ValorRecebido = sql.NullFloat64{Valid: true, Float64: transaction.Total}
	case domain.TransactionStatusCancelled:
		params.BoletosCancelados.Int32 = 1
	}

	return r.queries.UpdateBranchDailyStats(ctx, params)
}

func NewPostgresStatsRepository(sqlDB *sql.DB) repository.StatsRepository {
	return &postgresStatsRepository{
		queries: db.New(sqlDB),
	}
}

func (r *postgresStatsRepository) UpdateCustomerMonthlyStats(ctx context.Context, customerDocumentNumber string, month string, transaction domain.Transaction) error {
	params := db.UpdateCustomerMonthlyStatsParams{
		CustomerDocumentNumber: customerDocumentNumber,
		Month:                  time.Date(transaction.CreatedAt.Year(), transaction.CreatedAt.Month(), 1, 0, 0, 0, 0, time.UTC),
		TotalBoletos:           sql.NullInt32{Int32: 1, Valid: true},
		TotalPagos:             sql.NullInt32{Int32: 0, Valid: true},
		ValorEmitido:           sql.NullFloat64{Float64: transaction.Total, Valid: true},
		ValorRecebido:          sql.NullFloat64{Float64: 0, Valid: true},
		TenantID:               transaction.TenantId,
	}

	if transaction.Status == domain.TransactionStatusApproved {
		params.TotalPagos.Int32 = 1
		params.ValorRecebido = sql.NullFloat64{Float64: transaction.Total, Valid: true}
	}

	return r.queries.UpdateCustomerMonthlyStats(ctx, params)
}

func formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// CalculateCompanyStatistics - Calcula estatísticas por empresa
func (r *postgresStatsRepository) CalculateCompanyStatistics(ctx context.Context, tenantID, branchID string, startDate, endDate time.Time) (repository.CompanyStatisticsResult, error) {

	data, err := r.queries.GetCompanyStatistics(ctx, db.GetCompanyStatisticsParams{
		TenantID: tenantID,
		Date:     startDate,
		Date_2:   endDate,
	})
	if err != nil {
		return repository.CompanyStatisticsResult{}, err
	}

	var result repository.CompanyStatisticsResult
	result.TotalBoletos = int(data.TotalBoletos)
	result.TotalPagos = int(data.TotalPagos)
	result.ValorEmitido = float64(data.ValorEmitido)
	result.ValorRecebido = float64(data.ValorRecebido)
	result.BoletosCancelados = int(data.BoletosCancelados)
	result.ValorCancelado = float64(data.ValorCancelado)
	if result.TotalBoletos > 0 {
		result.TaxaPagamento = float64(result.TotalPagos) / float64(result.TotalBoletos) * 100
	}
	return result, nil
}
