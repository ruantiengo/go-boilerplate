-- name: CreateTransaction :exec
INSERT INTO Transaction (bank_slip_uuid, status, created_at, updated_at, payment_method)
VALUES ($1, $2, $3, $4, $5)
RETURNING bank_slip_uuid, status, created_at, updated_at, payment_method;

-- name: GetTransactionByUUID :one
SELECT bank_slip_uuid, status, created_at, updated_at, payment_method
FROM Transaction
WHERE bank_slip_uuid = $1;

-- name: UpdateTransactionStatus :exec
UPDATE Transaction
SET status = $1, updated_at = $2
WHERE bank_slip_uuid = $3;
    
-- name: UpdateTransaction :exec
UPDATE Transaction
SET
    status = COALESCE($1, status),
    updated_at = COALESCE($2, updated_at),
    payment_method = COALESCE($3, payment_method)
WHERE bank_slip_uuid = $4;  -- Added missing semicolon here

-- name: DeleteTransaction :exec
DELETE FROM Transaction
WHERE bank_slip_uuid = $1;

-- name: GetAllTransactions :many
SELECT bank_slip_uuid, status, created_at, updated_at, payment_method
FROM Transaction;