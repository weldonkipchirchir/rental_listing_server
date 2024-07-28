package db

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weldonkipchirchir/rental_listing/token"
	"github.com/weldonkipchirchir/rental_listing/util"
)

func TestJwtToken(t *testing.T) {
	user := CreateRandomUser(t)

	accessToken, refreshToken, err := token.GenerateJwtToken(user.ID, user.Email, user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	claims, msg := token.VerifyTokenString(accessToken)
	require.NotEmpty(t, claims)
	require.Empty(t, msg)

	claims2, msg2 := token.VerifyTokenString(refreshToken)
	require.NotEmpty(t, claims2)
	require.Empty(t, msg2)
}

func TestInvalidJwtTo(t *testing.T) {
	sample := util.RandomString(50)
	claims, msg := token.VerifyTokenString(sample)
	require.Empty(t, claims)
	require.NotEmpty(t, msg)
}
