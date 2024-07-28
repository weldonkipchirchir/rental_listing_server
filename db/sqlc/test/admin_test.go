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
func createRandomAdmin(t *testing.T) db.Admin {
	hashPassword, err := util.HashPassword(util.RandomString(7))
	require.NoError(t, err)

	arg := db.CreateAdminParams{
		Username:     util.RandomString(8),
		Email:        util.RandomEmail(),
		PasswordHash: hashPassword,
	}

	admin, err := testQueries.CreateAdmin(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, admin)

	require.Equal(t, arg.Username, admin.Username)
	require.Equal(t, arg.Email, admin.Email)
	require.NotZero(t, admin.CreatedAt.Valid)

	return db.Admin{
		ID:           admin.ID,
		Username:     admin.Username,
		Email:        admin.Email,
		PasswordHash: hashPassword,
		CreatedAt:    admin.CreatedAt,
	}
}

// TestCreateUser tests the user creation.
func TestCreateAdmin(t *testing.T) {
	createRandomAdmin(t)
}

// TestGetUser tests retrieving a user by email.
func TestGetAdmin(t *testing.T) {
	admin := createRandomAdmin(t)
	adminFromDB, err := testQueries.GetAdmin(context.Background(), admin.Email)
	require.NoError(t, err)
	require.NotEmpty(t, adminFromDB)

	require.Equal(t, admin.Username, adminFromDB.Username)
	require.Equal(t, admin.Email, adminFromDB.Email)
	require.NotZero(t, adminFromDB.CreatedAt.Valid)
	require.Equal(t, admin.PasswordHash, adminFromDB.PasswordHash)
	require.WithinDuration(t, admin.CreatedAt.Time, adminFromDB.CreatedAt.Time, 1000*time.Millisecond)
}

// TestGetUser tests retrieving a user by ID.
func TestGetAdminID(t *testing.T) {
	admin := createRandomAdmin(t)
	adminFromDB, err := testQueries.GetAdminByID(context.Background(), admin.ID)
	require.NoError(t, err)
	require.NotEmpty(t, adminFromDB)

	require.Equal(t, admin.Username, adminFromDB.Username)
	require.Equal(t, admin.Email, adminFromDB.Email)
	require.NotZero(t, adminFromDB.CreatedAt.Valid)
	require.WithinDuration(t, admin.CreatedAt.Time, adminFromDB.CreatedAt.Time, 1000*time.Millisecond)
}

// TestUpdateUserPassword tests updating a user's password.
func TestUpdateAdminPassword(t *testing.T) {
	admin := createRandomAdmin(t)
	newHashPassword, err := util.HashPassword(util.RandomString(7))
	newUsername := util.RandomString(8)
	require.NoError(t, err)

	arg := db.UpdateAdminPasswordParams{
		ID:           admin.ID,
		Username:     newUsername,
		PasswordHash: newHashPassword,
	}

	err = testQueries.UpdateAdminPassword(context.Background(), arg)
	require.NoError(t, err)

	adminFromDB, err := testQueries.GetAdmin(context.Background(), admin.Email)
	require.NoError(t, err)
	require.NotEmpty(t, adminFromDB)

	require.Equal(t, newUsername, adminFromDB.Username)
	require.NotZero(t, adminFromDB.CreatedAt.Valid)

	require.Equal(t, newHashPassword, adminFromDB.PasswordHash)
}

// TestDeleteUser tests deleting a user.
func TestDeleteAdmin(t *testing.T) {
	admin := createRandomAdmin(t)
	err := testQueries.DeleteUser(context.Background(), admin.ID)
	require.NoError(t, err)

	adminFromDB, err := testQueries.GetUser(context.Background(), admin.Email)
	require.Error(t, err)
	require.Empty(t, adminFromDB)
}
