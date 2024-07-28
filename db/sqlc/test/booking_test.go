package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/util"
)

func createUserBooking(t *testing.T) db.Booking {
	user := CreateRandomUser(t)
	listing := CreateListing(t)
	checkIn := util.GenerateDate(2024, time.January, 1, 0, 0, 0)
	checkOut := util.GenerateDate(2024, time.December, 1, 0, 0, 0)
	totalAmountInt := util.RandomFloat(10000, 1000000)
	totalAmount := fmt.Sprintf("%.2f", totalAmountInt)

	arg := db.CreateBookingParams{
		UserID:       user.ID,
		ListingID:    listing.ID,
		CheckInDate:  checkIn,
		CheckOutDate: checkOut,
		TotalAmount:  totalAmount,
	}

	booking, err := testQueries.CreateBooking(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, booking)

	require.Equal(t, arg.UserID, booking.UserID)
	require.Equal(t, arg.ListingID, booking.ListingID)
	require.Equal(t, arg.CheckInDate.UTC(), booking.CheckInDate.UTC())
	require.Equal(t, arg.CheckOutDate.UTC(), booking.CheckOutDate.UTC())
	require.Equal(t, arg.TotalAmount, booking.TotalAmount)
	require.NotZero(t, booking.CreatedAt.Valid)

	return db.Booking{
		ID:           booking.ID,
		UserID:       booking.UserID,
		ListingID:    booking.ListingID,
		CheckInDate:  booking.CheckInDate,
		CheckOutDate: booking.CheckOutDate,
		TotalAmount:  booking.TotalAmount,
		CreatedAt:    booking.CreatedAt,
	}
}

func TestCreateBooking(t *testing.T) {
	createUserBooking(t)
}

func TestGetBooking(t *testing.T) {
	booking := createUserBooking(t)
	arg := db.GetUserBookingByIDParams{
		ID:     booking.ID,
		UserID: booking.UserID,
	}
	bookingFromDB, err := testQueries.GetUserBookingByID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bookingFromDB)

	require.Equal(t, booking.UserID, bookingFromDB.UserID)
	require.Equal(t, booking.ListingID, bookingFromDB.ListingID)
	require.Equal(t, booking.CheckInDate.UTC(), bookingFromDB.CheckInDate.UTC())
	require.Equal(t, booking.CheckOutDate.UTC(), booking.CheckOutDate.UTC())
	require.Equal(t, booking.TotalAmount, bookingFromDB.TotalAmount)
	require.Equal(t, booking.CreatedAt.Time.UTC(), bookingFromDB.CreatedAt.Time.UTC())
}
