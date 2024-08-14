package service

import (
	"sea-study/api/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// CreateTopup
func TestCreateTopup(t *testing.T) {
	db, mock := setupTestDB(t)

	userID := uuid.New()
	amount := 100000.0
	paymentMethod := "credit_card"

	mock.ExpectBegin() 

	mock.ExpectQuery(`INSERT INTO "topup_histories" \("user_id","amount","status","payment_method","created_date"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "topup_id"`).
		WithArgs(sqlmock.AnyArg(), amount, "completed", paymentMethod, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"topup_id"}).AddRow(1))

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "balance"}).AddRow(userID, 0))

	mock.ExpectExec(`UPDATE "users" SET "name"=\$1,"email"=\$2,"password"=\$3,"balance"=\$4,"role"=\$5,"created_at"=\$6,"updated_at"=\$7 WHERE "user_id" = \$8`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), amount, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	topup, err := CreateTopup(db, userID.String(), amount, paymentMethod)

	assert.NoError(t, err)
	assert.NotNil(t, topup)
	assert.Equal(t, userID, topup.UserID)
	assert.Equal(t, amount, topup.Amount)
	assert.Equal(t, models.CompletedStatus, topup.Status)
	assert.Equal(t, paymentMethod, topup.PaymentMethod)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// GetTopupHistory
func TestGetTopupHistory(t *testing.T) {
	db, mock := setupTestDB(t)

	userID := uuid.New()

	mock.ExpectQuery(`SELECT \* FROM "topup_histories" WHERE user_id = \$1`).
		WithArgs(userID.String()).
		WillReturnRows(sqlmock.NewRows([]string{"topup_id", "user_id", "amount", "status", "payment_method", "created_date"}).
			AddRow(1, userID, 1000.0, "completed", "credit_card", time.Now()).
			AddRow(2, userID, 500.0, "completed", "credit_card", time.Now()))

	histories, err := GetTopupHistory(db, userID.String())

	assert.NoError(t, err)
	assert.Len(t, histories, 2)
	assert.Equal(t, userID.String(), histories[0].UserID.String())
	assert.Equal(t, 1000.0, histories[0].Amount)
	assert.Equal(t, models.CompletedStatus, histories[0].Status)
	assert.Equal(t, "credit_card", histories[0].PaymentMethod)
	assert.Equal(t, userID.String(), histories[1].UserID.String())
	assert.Equal(t, 500.0, histories[1].Amount)
	assert.Equal(t, models.CompletedStatus, histories[1].Status)
	assert.Equal(t, "credit_card", histories[1].PaymentMethod)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}