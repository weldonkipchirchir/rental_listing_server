package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(7)

	hashPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword1)

	err = VerifyPassword(password, hashPassword1)
	require.NoError(t, err)

	//assertain 2 password are not the same
	wrongPassword := RandomString(7)
	err = VerifyPassword(wrongPassword, hashPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	//generated hash password is different from each generated
	hashPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)
	require.NotEqual(t, hashPassword1, hashPassword2)
}
