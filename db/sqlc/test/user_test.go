package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/util"
)

// createRandomUser creates a random user for testing purposes.
func CreateRandomUser(t *testing.T) db.User {
	hashPassword, err := util.HashPassword(util.RandomString(7))
	require.NoError(t, err)

	arg := db.CreateUserParams{
		Username:     util.RandomString(8),
		Email:        util.RandomEmail(),
		PasswordHash: hashPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt.Valid)

	return db.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashPassword,
		CreatedAt:    user.CreatedAt,
	}
}

// TestCreateUser tests the user creation.
func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

// TestGetUser tests retrieving a user by email.
func TestGetUser(t *testing.T) {
	user := CreateRandomUser(t)
	userFromDB, err := testQueries.GetUser(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, userFromDB)

	require.Equal(t, user.Username, userFromDB.Username)
	require.Equal(t, user.Email, userFromDB.Email)
	require.NotZero(t, userFromDB.CreatedAt.Valid)
	require.Equal(t, user.PasswordHash, userFromDB.PasswordHash)
	require.WithinDuration(t, user.CreatedAt.Time, userFromDB.CreatedAt.Time, 1000*time.Millisecond)
}

// TestGetUser tests retrieving a user by ID.
func TestGetUserID(t *testing.T) {
	user := CreateRandomUser(t)
	userFromDB, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, userFromDB)

	require.Equal(t, user.Username, userFromDB.Username)
	require.Equal(t, user.Email, userFromDB.Email)
	require.NotZero(t, userFromDB.CreatedAt.Valid)
	require.WithinDuration(t, user.CreatedAt.Time, userFromDB.CreatedAt.Time, 1000*time.Millisecond)
}

// TestUpdateUserPassword tests updating a user's password.
func TestUpdateUserPassword(t *testing.T) {
	user := CreateRandomUser(t)
	newHashPassword, err := util.HashPassword(util.RandomString(7))
	require.NoError(t, err)

	arg := db.UpdateUserPasswordParams{
		ID:           user.ID,
		PasswordHash: newHashPassword,
	}

	err = testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)

	userFromDB, err := testQueries.GetUser(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, userFromDB)

	require.Equal(t, user.Username, userFromDB.Username)
	require.Equal(t, user.Email, userFromDB.Email)
	require.NotZero(t, userFromDB.CreatedAt.Valid)
	require.Equal(t, newHashPassword, userFromDB.PasswordHash)
}

// TestDeleteUser tests deleting a user.
func TestDeleteUser(t *testing.T) {
	user := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	userFromDB, err := testQueries.GetUser(context.Background(), user.Email)
	require.Error(t, err)
	require.Empty(t, userFromDB)
}
