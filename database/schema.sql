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
    bank_slip_uuid UUID PRIMARY KEY,       -- UUID como chave primária
    status transaction_status,             -- Status da transação (ENUM)
    created_at TIMESTAMP NOT NULL,         -- Data de criação da transação
    updated_at TIMESTAMP NOT NULL,         -- Data de atualização da transação
    payment_method payment_method,         -- Método de pagamento (ENUM)
    CONSTRAINT check_status CHECK (status IN ('pending', 'cancelled', 'expired', 'approved')),   -- Validação para garantir valores corretos
    CONSTRAINT check_payment_method CHECK (payment_method IN ('bill', 'pix', 'credit_card'))   -- Validação para garantir valores corretos
);
