package repository

import (
	"context"
	"ruantiengo/internal/domain"
	"time"
)

type StatsRepository interface {
	UpdateBranchDailyStats(ctx context.Context, branchId string, date time.Time, transaction domain.Transaction) error
	UpdateCustomerMonthlyStats(ctx context.Context, customerId string, month string, transaction domain.Transaction) error
	CalculateCompanyStatistics(ctx context.Context, tenantID string, branchID string, startDate time.Time, endDate time.Time) (CompanyStatisticsResult, error)
}

type CompanyStatisticsResult struct {
	TotalBoletos        int     `json:"total_boletos"`
	TotalPagos          int     `json:"total_pagos"`
	TaxaPagamento       float64 `json:"taxa_pagamento"`
	ValorEmitido        float64 `json:"valor_emitido"`
	ValorRecebido       float64 `json:"valor_recebido"`
	BoletosCancelados   int     `json:"boletos_cancelados"`
	ValorCancelado      float64 `json:"valor_cancelado"`
	MediaDiasAtraso     float64 `json:"media_dias_atraso"`
	PercentualAtrasados float64 `json:"percentual_atrasados"`
}

type CustomerStatisticsResult struct {
	CustomerID           string  `json:"customer_id"`
	Pontualidade         float64 `json:"pontualidade"`
	ValorMedio           float64 `json:"valor_medio"`
	MediaAtrasoDias      float64 `json:"media_atraso_dias"`
	TaxaAtraso           float64 `json:"taxa_atraso"`
	FrequenciaPagamentos int     `json:"frequencia_pagamentos"`
	ConcentracaoReceita  float64 `json:"concentracao_receita"`
	Classificacao        string  `json:"classificacao"`
}
