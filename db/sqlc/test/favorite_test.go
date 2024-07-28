package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

func createRandomFavorite(t *testing.T) db.Favorite {
	user := CreateRandomUser(t)
	listing := CreateListing(t)

	arg := db.CreateFavoriteParams{
		UserID:    user.ID,
		ListingID: listing.ID,
	}

	favorite, err := testQueries.CreateFavorite(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, favorite)

	require.Equal(t, arg.UserID, favorite.UserID)
	require.Equal(t, arg.ListingID, favorite.ListingID)
	require.NotZero(t, favorite.CreatedAt.Valid)

	return db.Favorite{
		ID:        favorite.ID,
		UserID:    favorite.UserID,
		ListingID: favorite.ListingID,
		CreatedAt: favorite.CreatedAt,
	}
}

func TestCreateFavorite(t *testing.T) {
	createRandomFavorite(t)
}

func TestGetFavorite(t *testing.T) {
	favorite := createRandomFavorite(t)
	favoriteFromDB, err := testQueries.GetFavorite(context.Background(), favorite.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, favoriteFromDB)

	require.Equal(t, favorite.ID, favoriteFromDB[0].ID)
	require.Equal(t, favorite.UserID, favoriteFromDB[0].UserID)
	require.Equal(t, favorite.ListingID, favoriteFromDB[0].ListingID)
	require.Equal(t, favorite.CreatedAt, favoriteFromDB[0].CreatedAt)
}
