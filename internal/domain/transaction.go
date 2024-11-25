package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCancelled TransactionStatus = "cancelled"
	TransactionStatusExpired   TransactionStatus = "expired"
	TransactionStatusApproved  TransactionStatus = "approved"
)

const (
	StatusPending   TransactionStatus = "pending"
	StatusApproved  TransactionStatus = "approved"
	StatusCancelled TransactionStatus = "cancelled"
	StatusExpired   TransactionStatus = "expired"
)

type PaymentMethod string
type Transaction struct {
	Id                     uuid.UUID
	BankSlipUuid           uuid.UUID
	Status                 TransactionStatus
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DueDate                time.Time
	Total                  float64
	CustomerDocumentNumber string
	TenantId               string
	BranchId               string
	PaymentMethod          PaymentMethod
}
