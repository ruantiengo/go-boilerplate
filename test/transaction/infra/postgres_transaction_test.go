package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	db "ruantiengo/database/generated"
	"ruantiengo/internal/transaction/domain"
	"ruantiengo/internal/transaction/infra"
)

func TestPostgresTransactionRepository(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	repo := infra.NewPostgresTransactionRepository(sqlDB)

	t.Run("Save", func(t *testing.T) {
		ctx := context.Background()
		transaction := domain.Transaction{
			BankSlipUuid:  uuid.New(),
			Status:        domain.TransactionStatusApproved,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PaymentMethod: domain.PaymentMethodCreditCard,
		}

		mock.ExpectExec("INSERT INTO Transaction").
			WithArgs(
				transaction.BankSlipUuid,
				db.TransactionStatus(transaction.Status),
				transaction.CreatedAt,
				transaction.UpdatedAt,
				db.PaymentMethod(transaction.PaymentMethod),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(ctx, transaction)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Save Error", func(t *testing.T) {
		ctx := context.Background()
		transaction := domain.Transaction{
			BankSlipUuid:  uuid.New(),
			Status:        domain.TransactionStatusApproved,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			PaymentMethod: domain.PaymentMethodCreditCard,
		}

		mock.ExpectExec("INSERT INTO Transaction").
			WithArgs(
				transaction.BankSlipUuid,
				db.TransactionStatus(transaction.Status),
				transaction.CreatedAt,
				transaction.UpdatedAt,
				db.PaymentMethod(transaction.PaymentMethod),
			).
			WillReturnError(errors.New("database error"))

		err := repo.Save(ctx, transaction)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
