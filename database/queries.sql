-- name: CreateTransaction :one
INSERT INTO transaction (id, status, created_at, updated_at, due_date, total, customer_id, tenant_id, branch_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transaction WHERE id = $1;

-- name: UpdateTransactionStatus :exec
UPDATE transaction SET status = $2, updated_at = $3 WHERE id = $1;

-- name: DeleteTransaction :exec
DELETE FROM transaction WHERE id = $1;