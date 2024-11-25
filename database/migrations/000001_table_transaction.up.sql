CREATE TYPE transaction_status AS ENUM (
    'pending',
    'cancelled',
    'expired',
    'approved'
);

CREATE TYPE payment_method AS ENUM (
    'bill',
    'pix',
    'credit_card'
);

CREATE TABLE Transaction (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),       -- id gerado por default como chave primária
    bank_slip_uuid UUID,                                 -- UUID do bank slip
    status transaction_status,                           -- Status da transação (ENUM)
    created_at TIMESTAMP NOT NULL,                       -- Data de criação da transação
    updated_at TIMESTAMP NOT NULL,                       -- Data de atualização da transação
    due_date TIMESTAMP NOT NULL,                         -- Data de vencimento da transação
    total NUMERIC NOT NULL,                              -- Valor total da transação
    customer_id VARCHAR NOT NULL,                        -- ID do cliente
    tenant_id VARCHAR NOT NULL,                          -- ID do locatário
    branch_id VARCHAR NOT NULL,                          -- ID da filial
    payment_method payment_method,                       -- Método de pagamento (ENUM)
    CONSTRAINT check_status CHECK (status IN ('pending', 'cancelled', 'expired', 'approved')),   -- Validação para garantir valores corretos
    CONSTRAINT check_payment_method CHECK (payment_method IN ('bill', 'pix', 'credit_card'))     -- Validação para garantir valores corretos
);
