package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/util"
)

func createRandomReviews(t *testing.T) db.Review {
	user := CreateRandomUser(t)
	listing := CreateListing(t)

	arg := db.CreateReviewParams{
		UserID:    user.ID,
		ListingID: listing.ID,
		Rating:    int32(util.RandomFloat(1, 5)),
		Comment:   sql.NullString{String: util.RandomString(10), Valid: true},
	}

	review, err := testQueries.CreateReview(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, review)

	return review
}

func TestCreateReview(t *testing.T) {
	createRandomReviews(t)
}
