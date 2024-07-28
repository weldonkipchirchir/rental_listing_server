package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

type createReviewRequest struct {
	ListingID int32  `json:"listingId" binding:"required"`
	Rating    int32  `json:"rating" binding:"required"`
	Comment   string `json:"comment" binding:"required"`
}

func (s *Server) CreateReview(c *gin.Context) {
	var request createReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized users only"})
		return
	}

	_, err = s.q.GetListingByID(c, request.ListingID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	_, err = s.q.CreateReview(c, db.CreateReviewParams{
		UserID:    user.ID,
		ListingID: request.ListingID,
		Rating:    request.Rating,
		Comment:   sql.NullString{String: request.Comment, Valid: true},
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Review created successfully"})

}

type getReviews struct {
	ID        int32     `json:"id" binding:"required,alphanum"`
	UserID    int32     `json:"user_id" binding:"required,alphanum"`
	ListingID int32     `json:"listing_id" binding:"required,alphanum"`
	Rating    int32     `json:"rating" binding:"required,alphanum"`
	Comment   string    `json:"comment" binding:"required,alphanum"`
	CreatedAt time.Time `json:"created_at" binding:"required,alphanum"`
}

func (s *Server) GetReviewByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized users only"})
		return
	}

	arg := db.GetReviewByIDParams{
		ID:     int32(id),
		UserID: user.ID,
	}

	review, err := s.q.GetReviewByID(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, []getReviews{})
			return
		}
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getReviews{
		ID:        review.ID,
		UserID:    review.UserID,
		ListingID: review.ListingID,
		Rating:    review.Rating,
		Comment:   review.Comment.String,
		CreatedAt: review.CreatedAt.Time,
	})
}

type getReviewsByListing struct {
	ListingID int32 `json:"listing_id" binding:"required"`
}

func (s *Server) GetReviewsByListing(c *gin.Context) {
	var req getReviewsByListing
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized admin only"})
		return
	}

	arg := db.GetReviewsByListingIDParams{
		ListingID: req.ListingID,
		AdminID:   admin.ID,
	}

	reviews, err := s.q.GetReviewsByListingID(c, arg)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusOK, []getReviews{})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (s *Server) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorized users only"})
		return
	}

	arg1 := db.GetReviewByIDParams{
		ID:     int32(id),
		UserID: user.ID,
	}

	_, err = s.q.GetReviewByID(c, arg1)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.DeleteReviewParams{
		ID:     int32(id),
		UserID: user.ID,
	}

	err = s.q.DeleteReview(c, arg)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
