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

type PaymentMethod string

const (
	PaymentMethodBill       PaymentMethod = "bill"
	PaymentMethodPix        PaymentMethod = "pix"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
)

type Transaction struct {
	BankSlipUuid  uuid.UUID
	Status        TransactionStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
	PaymentMethod PaymentMethod
}
