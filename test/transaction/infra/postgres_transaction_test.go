package infra

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "ruantiengo/database/generated"
	"ruantiengo/internal/transaction/domain"
	"ruantiengo/internal/transaction/infra"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostgresTransactionRepository_Upsert(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := infra.NewPostgresTransactionRepository(mockDB)

	now := time.Now()
	sampleTransaction := domain.Transaction{
		BankSlipUuid:  uuid.New(),
		Status:        domain.TransactionStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
		PaymentMethod: domain.PaymentMethodCreditCard,
	}

	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(
			sampleTransaction.BankSlipUuid,
			db.TransactionStatus(sampleTransaction.Status),
			sampleTransaction.CreatedAt,
			sampleTransaction.UpdatedAt,
			db.PaymentMethod(sampleTransaction.PaymentMethod),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Upsert(context.Background(), sampleTransaction)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestPostgresTransactionRepository_Upsert_Error(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := infra.NewPostgresTransactionRepository(mockDB)

	now := time.Now()
	sampleTransaction := domain.Transaction{
		BankSlipUuid:  uuid.New(),
		Status:        domain.TransactionStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
		PaymentMethod: domain.PaymentMethodCreditCard,
	}

	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(
			sampleTransaction.BankSlipUuid,
			db.TransactionStatus(sampleTransaction.Status),
			sampleTransaction.CreatedAt,
			sampleTransaction.UpdatedAt,
			db.PaymentMethod(sampleTransaction.PaymentMethod),
		).
		WillReturnError(sql.ErrConnDone)

	err = repo.Upsert(context.Background(), sampleTransaction)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
