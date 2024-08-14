package service

import (
	"sea-study/api/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// CreateUser
func TestCreateUser(t *testing.T) {
	db, mock := setupTestDB(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO \"users"\ (.+) VALUES (.+)`).
		WillReturnRows(mock.NewRows([]string{"user_id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	var user models.User
	err := CreateUser(db, &user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// GetUserByEmail
func TestGetUserByEmail(t *testing.T) {
	db, mock := setupTestDB(t)

	email := "test@example.com"
	expectedUser := &models.User{
		UserID:    uuid.New(),
		Email:     email,
		Name:      "Test User",
		Password:  "hashed_password",
		Balance:   100.00,
		Role:      models.RoleEnum("user"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := mock.NewRows([]string{"user_id", "name", "email", "password", "balance", "role", "created_at", "updated_at"}).
		AddRow(expectedUser.UserID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.Balance, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(email, 1).
		WillReturnRows(rows)

	var user models.User
	err := GetUserByEmail(db, &user, email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.UserID, user.UserID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.Equal(t, expectedUser.Balance, user.Balance)
	assert.Equal(t, expectedUser.Role, user.Role)
	assert.Equal(t, expectedUser.CreatedAt.Unix(), user.CreatedAt.Unix())
	assert.Equal(t, expectedUser.UpdatedAt.Unix(), user.UpdatedAt.Unix())

	assert.NoError(t, mock.ExpectationsWereMet())
}

// GetUserByID
func TestGetUserByID(t *testing.T) {
	db, mock := setupTestDB(t)

	userId := uuid.New()
	expectedUser := &models.User{
		UserID:    userId,
		Email:     "dimas@gmail.com",
		Name:      "Test User",
		Password:  "hashed_password",
		Balance:   100.00,
		Role:      models.RoleEnum("user"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := mock.NewRows([]string{"user_id", "name", "email", "password", "balance", "role", "created_at", "updated_at"}).
		AddRow(expectedUser.UserID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.Balance, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(userId, 1).
		WillReturnRows(rows)

	var user models.User
	err := GetUserByID(db, &user, userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.UserID, user.UserID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Password, user.Password)
	assert.Equal(t, expectedUser.Balance, user.Balance)
	assert.Equal(t, expectedUser.Role, user.Role)
	assert.Equal(t, expectedUser.CreatedAt.Unix(), user.CreatedAt.Unix())
	assert.Equal(t, expectedUser.UpdatedAt.Unix(), user.UpdatedAt.Unix())

	assert.NoError(t, mock.ExpectationsWereMet())
}

// ValidateUserExists
func TestValidateUserExists(t *testing.T) {
	db, mock := setupTestDB(t)

	userId := uuid.New()
	expectedUser := models.User{
		UserID: userId,
	}

	rows := mock.NewRows([]string{"user_id"}).AddRow(expectedUser.UserID)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(userId, 1).
		WillReturnRows(rows)

	err := ValidateUserExists(db, userId)

	assert.NoError(t, err)
	assert.True(t, err == nil)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// GetUserBalance
func TestGetUserBalance(t *testing.T) {
	db, mock := setupTestDB(t)

	userId := uuid.New()
	expectedBalance := 100.00

	rows := mock.NewRows([]string{"user_id", "name", "email", "password", "balance", "role", "created_at", "updated_at"}).
		AddRow(userId, "Test User", "test@example.com", "hashed_password", expectedBalance, "user", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(userId, 1).
		WillReturnRows(rows)

	balance, err := GetUserBalance(db, userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// UpdateUserBalance
func TestUpdateUserBalance(t *testing.T) {
	db, mock := setupTestDB(t)

	userId := uuid.New()
	initialBalance := 100.00
	updateAmount := 50.00
	expectedBalance := 150.00

	selectRows := mock.NewRows([]string{"user_id", "name", "email", "password", "balance", "role", "created_at", "updated_at"}).
		AddRow(userId, "Test User", "test@example.com", "hashed_password", initialBalance, "user", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(userId, 1).
		WillReturnRows(selectRows)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "name"=\$1,"email"=\$2,"password"=\$3,"balance"=\$4,"role"=\$5,"created_at"=\$6,"updated_at"=\$7 WHERE "user_id" = \$8`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectedBalance, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := UpdateUserBalance(db, userId, updateAmount)

	assert.NoError(t, err)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(userId, 1).
		WillReturnRows(mock.NewRows([]string{"balance"}).AddRow(expectedBalance))

	newBalance, err := GetUserBalance(db, userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, newBalance)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserBalanceInsufficientFunds(t *testing.T) {
	db, mock := setupTestDB(t)

	userId := uuid.New()
	initialBalance := 100.00
	updateAmount := -150.00 

	selectRows := mock.NewRows([]string{"user_id", "name", "email", "password", "balance", "role", "created_at", "updated_at"}).
		AddRow(userId, "Test User", "test@example.com", "hashed_password", initialBalance, "user", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE user_id = \$1 ORDER BY "users"."user_id" LIMIT \$2`).
		WithArgs(userId, 1).
		WillReturnRows(selectRows)

	err := UpdateUserBalance(db, userId, updateAmount)

	assert.Error(t, err)
	assert.Equal(t, "insufficient balance", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}