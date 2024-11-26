-- name: CreateTransaction :one

-- ############################
-- # QUERIES PARA TRANSACTION #
-- ############################

-- Criação de uma nova transação
-- name: CreateTransaction :exec
INSERT INTO Transaction (
    bank_slip_uuid, status, created_at, updated_at, due_date, total, customer_document_number, tenant_id, branch_id, payment_method
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetTransactionByUUID :one
SELECT * FROM Transaction
WHERE id = $1;

-- name: UpdateTransactionStatus :exec
UPDATE Transaction
SET status = $1, updated_at = $2
WHERE id = $3;

-- name: UpdateTransaction :exec
UPDATE Transaction
SET
    status = COALESCE($1, status),
    updated_at = COALESCE($2, updated_at),
    due_date = COALESCE($3, due_date),
    total = COALESCE($4, total),
    payment_method = COALESCE($5, payment_method)
WHERE id = $6;

-- name: DeleteTransaction :exec
DELETE FROM Transaction
WHERE id = $1;

-- name: GetAllTransactions :many
SELECT * FROM Transaction;

-- name: UpsertTransaction :one
INSERT INTO Transaction (
    bank_slip_uuid, status, created_at, updated_at, due_date, total, customer_document_number, tenant_id, branch_id, payment_method
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (id) DO UPDATE
SET
    bank_slip_uuid = EXCLUDED.bank_slip_uuid,
    status = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at,
    due_date = EXCLUDED.due_date,
    total = EXCLUDED.total,
    customer_document_number = EXCLUDED.customer_document_number,
    tenant_id = EXCLUDED.tenant_id,
    branch_id = EXCLUDED.branch_id,
    payment_method = EXCLUDED.payment_method
RETURNING *;

-- name: GetTransactionsByCustomerDocument :many
SELECT * FROM Transaction
WHERE customer_document_number = $1
ORDER BY created_at DESC;

-- name: GetTransactionsByTenantId :many
SELECT * FROM Transaction
WHERE tenant_id = $1
ORDER BY created_at DESC;

-- name: GetTransactionsByBranchId :many
SELECT * FROM Transaction
WHERE branch_id = $1
ORDER BY created_at DESC;

-- name: GetTransactionsByStatus :many
SELECT * FROM Transaction
WHERE status = $1
ORDER BY created_at DESC;

-- name: GetTransactionsByDateRange :many
SELECT * FROM Transaction
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC;

-- name: GetTransactionsTotalByTenantId :one
SELECT SUM(total) as total_amount
FROM Transaction
WHERE tenant_id = $1;

-- name: GetTransactionsTotalByBranchId :one
SELECT SUM(total) as total_amount
FROM Transaction
WHERE branch_id = $1;

-- name: GetTransactionsCountByStatus :many
SELECT status, COUNT(*) as count
FROM Transaction
GROUP BY status;

-- name: GetLatestTransactions :many
SELECT * FROM Transaction
ORDER BY created_at DESC
LIMIT $1;


-- ############################
-- # QUERIES PARA BRANCHDAILYSTATS #
-- ############################

-- Total de boletos gerados por filial
-- name: GetTotalBoletosGerados :one
SELECT SUM(total_boletos) AS total_boletos_gerados
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- Total de boletos pagos por filial
-- name: GetTotalBoletosPagos :one
SELECT SUM(total_pagos) AS total_boletos_pagos
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- Taxa de pagamento por filial
-- name: GetTaxaPagamento :one
SELECT 
    SUM(total_pagos)::DECIMAL / NULLIF(SUM(total_boletos), 0) * 100 AS taxa_pagamento
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- Valor total emitido por filial
-- name: GetValorTotalEmitido :one
SELECT SUM(valor_emitido) AS valor_total_emitido
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- Valor total recebido por filial
-- name: GetValorTotalRecebido :one
SELECT SUM(valor_recebido) AS valor_total_recebido
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- Média de atraso no pagamento por filial
-- name: GetMediaAtrasoPagamento :one
SELECT 
    SUM(total_dias_atraso)::DECIMAL / NULLIF(SUM(boletos_atrasados), 0) AS media_atraso
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- Percentual de boletos atrasados por filial
-- name: GetPercentualBoletosAtrasados :one
SELECT 
    SUM(boletos_atrasados)::DECIMAL / NULLIF(SUM(total_pagos), 0) * 100 AS percentual_atrasados
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- Boletos cancelados por filial
-- name: GetBoletosCancelados :one
SELECT 
    SUM(boletos_cancelados) AS total_boletos_cancelados,
    SUM(valor_cancelado) AS valor_total_cancelado
FROM BranchDailyStats
WHERE date = CURRENT_DATE;

-- ############################
-- # QUERIES PARA CUSTOMERMONTHLYSTATS #
-- ############################

-- Pontualidade dos clientes
-- name: GetCustomerPontualidade :many
SELECT 
    customer_document_number,
    (SUM(total_pagos) - SUM(boletos_atrasados))::DECIMAL / NULLIF(SUM(total_pagos), 0) * 100 AS pontualidade
FROM CustomerMonthlyStats
WHERE month = DATE_TRUNC('month', CURRENT_DATE)
GROUP BY customer_document_number;

-- Valor médio dos boletos por cliente
-- name: GetCustomerValorMedio :many
SELECT 
    customer_document_number,
    SUM(valor_emitido)::DECIMAL / NULLIF(SUM(total_boletos), 0) AS valor_medio
FROM CustomerMonthlyStats
WHERE month = DATE_TRUNC('month', CURRENT_DATE)
GROUP BY customer_document_number;

-- Tempo médio de atraso por cliente
-- name: GetCustomerMediaAtraso :many
SELECT 
    customer_document_number,
    SUM(total_dias_atraso)::DECIMAL / NULLIF(SUM(total_pagos), 0) AS tempo_medio_atraso
FROM CustomerMonthlyStats
WHERE month = DATE_TRUNC('month', CURRENT_DATE)
GROUP BY customer_document_number;

-- Taxa de atraso por cliente
-- name: GetCustomerTaxaAtraso :many
SELECT 
    customer_document_number,
    SUM(boletos_atrasados)::DECIMAL / NULLIF(SUM(total_pagos), 0) * 100 AS taxa_atraso
FROM CustomerMonthlyStats
WHERE month = DATE_TRUNC('month', CURRENT_DATE)
GROUP BY customer_document_number;

-- Frequência de pagamentos por cliente
-- name: GetCustomerFrequencia :many
SELECT 
    customer_document_number,
    SUM(total_pagos) AS num_pagamentos
FROM CustomerMonthlyStats
WHERE month = DATE_TRUNC('month', CURRENT_DATE)
GROUP BY customer_document_number;

-- Concentração de receita por cliente
-- name: GetCustomerConcentracaoReceita :many
WITH total_recebido AS (
    SELECT SUM(valor_recebido) AS valor_total_recebido
    FROM CustomerMonthlyStats
    WHERE month = DATE_TRUNC('month', CURRENT_DATE)
)
SELECT 
    customer_document_number,
    SUM(valor_recebido)::DECIMAL / (SELECT valor_total_recebido FROM total_recebido) * 100 AS concentracao_receita
FROM CustomerMonthlyStats
WHERE month = DATE_TRUNC('month', CURRENT_DATE)
GROUP BY customer_document_number;

-- Classificação de risco por cliente
-- name: GetCustomerClassificacaoRisco :many
SELECT 
    customer_document_number,
    CASE
        WHEN SUM(total_dias_atraso) = 0 THEN 'Bom Pagador'
        WHEN SUM(total_dias_atraso) <= 7 THEN 'Pagador Regular'
        ELSE 'Mau Pagador'
    END AS classificacao_risco
FROM CustomerMonthlyStats
WHERE month = DATE_TRUNC('month', CURRENT_DATE)
GROUP BY customer_document_number;

-- ############################
-- # CONSULTAS EXTRAS #
-- ############################

-- Boletos gerados por filial
-- name: GetBoletosPorFilial :many
SELECT 
    tenant_id,
    branch_id,
    SUM(total_boletos) AS total_boletos_gerados
FROM BranchDailyStats
WHERE date = CURRENT_DATE
GROUP BY tenant_id, branch_id;

-- Taxa de pagamento por filial
-- name: GetTaxaPagamentoPorFilial :many
SELECT 
    tenant_id,
    branch_id,
    SUM(total_pagos)::DECIMAL / NULLIF(SUM(total_boletos), 0) * 100 AS taxa_pagamento
FROM BranchDailyStats
WHERE date = CURRENT_DATE
GROUP BY tenant_id, branch_id;

-- Valor médio dos boletos por filial
-- name: GetValorMedioPorFilial :many
SELECT 
    tenant_id,
    branch_id,
    SUM(valor_emitido)::DECIMAL / NULLIF(SUM(total_boletos), 0) AS valor_medio
FROM BranchDailyStats
WHERE date = CURRENT_DATE
GROUP BY tenant_id, branch_id;

-- Tempo médio de recebimento por filial
-- name: GetTempoMedioRecebimentoFilial :many
SELECT 
    tenant_id,
    branch_id,
    SUM(total_dias_atraso)::DECIMAL / NULLIF(SUM(total_pagos), 0) AS tempo_medio_recebimento
FROM BranchDailyStats
WHERE date = CURRENT_DATE
GROUP BY tenant_id, branch_id;

-- name: UpdateBranchDailyStats :exec
INSERT INTO BranchDailyStats (
    tenant_id,
    branch_id,
    date,
    total_boletos,
    total_pagos,
    valor_emitido,
    valor_recebido,
    boletos_cancelados,
    valor_cancelado,
    boletos_atrasados,
    total_dias_atraso
) VALUES (
    $1, -- tenant_id
    $2, -- branch_id
    $3, -- date
    $4, -- delta_total_boletos
    $5, -- delta_total_pagos
    $6, -- delta_valor_emitido
    $7, -- delta_valor_recebido
    $8, -- delta_boletos_cancelados
    $9, -- delta_valor_cancelado
    $10, -- delta_boletos_atrasados
    $11 -- delta_total_dias_atraso
) ON CONFLICT (tenant_id, branch_id, date) DO UPDATE SET
    total_boletos = BranchDailyStats.total_boletos + EXCLUDED.total_boletos,
    total_pagos = BranchDailyStats.total_pagos + EXCLUDED.total_pagos,
    valor_emitido = BranchDailyStats.valor_emitido + EXCLUDED.valor_emitido,
    valor_recebido = BranchDailyStats.valor_recebido + EXCLUDED.valor_recebido,
    boletos_cancelados = BranchDailyStats.boletos_cancelados + EXCLUDED.boletos_cancelados,
    valor_cancelado = BranchDailyStats.valor_cancelado + EXCLUDED.valor_cancelado,
    boletos_atrasados = BranchDailyStats.boletos_atrasados + EXCLUDED.boletos_atrasados,
    total_dias_atraso = BranchDailyStats.total_dias_atraso + EXCLUDED.total_dias_atraso;

-- name: UpdateCustomerMonthlyStats :exec
INSERT INTO CustomerMonthlyStats (
    customer_document_number,
    tenant_id,
    month,
    total_boletos,
    total_pagos,
    valor_emitido,
    valor_recebido,
    boletos_atrasados,
    total_dias_atraso
) VALUES (
    $1, -- customer_document_number
    $2, -- tenant_id
    $3, -- month
    $4, -- delta_total_boletos
    $5, -- delta_total_pagos
    $6, -- delta_valor_emitido
    $7, -- delta_valor_recebido
    $8, -- delta_boletos_atrasados
    $9  -- delta_total_dias_atraso
) ON CONFLICT (customer_document_number, month) DO UPDATE SET
    total_boletos = CustomerMonthlyStats.total_boletos + EXCLUDED.total_boletos,
    total_pagos = CustomerMonthlyStats.total_pagos + EXCLUDED.total_pagos,
    valor_emitido = CustomerMonthlyStats.valor_emitido + EXCLUDED.valor_emitido,
    valor_recebido = CustomerMonthlyStats.valor_recebido + EXCLUDED.valor_recebido,
    boletos_atrasados = CustomerMonthlyStats.boletos_atrasados + EXCLUDED.boletos_atrasados,
    total_dias_atraso = CustomerMonthlyStats.total_dias_atraso + EXCLUDED.total_dias_atraso;

-- name: GetBranchDailyStats :many
SELECT
    branch_id,
    date,
    total_boletos,
    total_pagos,
    valor_emitido,
    valor_recebido,
    boletos_cancelados,
    valor_cancelado,
    boletos_atrasados,
    total_dias_atraso
FROM
    branchdailystats
WHERE
    tenant_id = $1
    AND branch_id = $2
    AND date BETWEEN $3 AND $4;

-- name: GetCompanyStatistics :one
SELECT
    SUM(total_boletos) AS total_boletos,
    SUM(total_pagos) AS total_pagos,
    SUM(valor_emitido) AS valor_emitido,
    SUM(valor_recebido) AS valor_recebido,
    SUM(boletos_cancelados) AS boletos_cancelados,
    SUM(valor_cancelado) AS valor_cancelado,
    SUM(boletos_atrasados) AS boletos_atrasados,
    SUM(total_dias_atraso) AS total_dias_atraso
FROM
    branchdailystats
WHERE
    tenant_id = $1
    AND date BETWEEN $2 AND $3;
