package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/util"
)

func CreateListing(t *testing.T) db.Listing {
	admin := createRandomAdmin(t)
	priceInt := util.RandomFloat(1, 1000000)
	price := fmt.Sprintf("%.2f", priceInt)

	arg := db.CreateListingParams{
		Title:       util.RandomString(10),
		Price:       price,
		AdminID:     admin.ID,
		Description: sql.NullString{String: util.RandomString(10), Valid: true},
		Location:    sql.NullString{String: util.RandomString(10), Valid: true},
		Available:   sql.NullBool{Bool: true, Valid: true},
		Column7:     []string{util.RandomString(10)},
	}

	listing, err := testQueries.CreateListing(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, listing)

	require.Equal(t, arg.Title, listing.Title)
	require.Equal(t, arg.Price, listing.Price)
	require.Equal(t, arg.AdminID, listing.AdminID)
	require.Equal(t, arg.Description, listing.Description)
	require.Equal(t, arg.Location, listing.Location)
	require.Equal(t, arg.Available, listing.Available)
	require.Equal(t, arg.Column7, listing.Imagelinks)
	require.NotZero(t, listing.CreatedAt.Valid)

	return db.Listing{
		ID:          listing.ID,
		Title:       listing.Title,
		Price:       listing.Price,
		AdminID:     listing.AdminID,
		Description: listing.Description,
		Location:    listing.Location,
		Available:   listing.Available,
		Imagelinks:  listing.Imagelinks,
		CreatedAt:   listing.CreatedAt,
	}
}

func TestCreateListing(t *testing.T) {
	CreateListing(t)
}

func TestGetListing(t *testing.T) {
	listing := CreateListing(t)
	listingFromDB, err := testQueries.GetListings(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, listingFromDB)

	require.Equal(t, listing.Title, listingFromDB[0].Title)
	require.Equal(t, listing.Price, listingFromDB[0].Price)
	require.Equal(t, listing.AdminID, listingFromDB[0].AdminID)
	require.Equal(t, listing.Description, listingFromDB[0].Description)
	require.Equal(t, listing.Location, listingFromDB[0].Location)
	require.Equal(t, listing.Available, listingFromDB[0].Available)
	require.Equal(t, listing.Imagelinks, listingFromDB[0].Imagelinks)
	require.NotZero(t, listingFromDB[0].CreatedAt.Valid)
}
