CREATE TABLE BranchDailyStats (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR NOT NULL,
    branch_id VARCHAR NOT NULL,
    date DATE NOT NULL,
    total_boletos INT DEFAULT 0,
    total_pagos INT DEFAULT 0,
    valor_emitido DOUBLE PRECISION DEFAULT 0.00,
    valor_recebido DOUBLE PRECISION DEFAULT 0.00,
    boletos_cancelados INT DEFAULT 0,
    valor_cancelado DOUBLE PRECISION DEFAULT 0.00,
    boletos_atrasados INT DEFAULT 0,
    total_dias_atraso DOUBLE PRECISION DEFAULT 0.00,
    CONSTRAINT unique_tenant_branch_date UNIQUE (tenant_id, branch_id, date)
);

CREATE TABLE CustomerMonthlyStats (
    id SERIAL PRIMARY KEY,
    customer_document_number VARCHAR NOT NULL,
    tenant_id VARCHAR NOT NULL,
    month DATE NOT NULL,
    total_boletos INT DEFAULT 0,
    total_pagos INT DEFAULT 0,
    valor_emitido DOUBLE PRECISION DEFAULT 0.00,
    valor_recebido DOUBLE PRECISION DEFAULT 0.00,
    boletos_atrasados INT DEFAULT 0,
    total_dias_atraso DOUBLE PRECISION DEFAULT 0.00,
    CONSTRAINT unique_customer_month UNIQUE (customer_document_number, month)
);

