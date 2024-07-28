package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

type createFavoriteRequest struct {
	ListingID int32 `json:"listing_id" binding:"required"`
}

func (s *Server) CreateFavorite(c *gin.Context) {
	var req createFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	listing, err := s.q.GetListingByID(c, req.ListingID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	_, err = s.q.CreateFavorite(c, db.CreateFavoriteParams{
		UserID:    user.ID,
		ListingID: listing.ID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Favorite created successfully"})

}

type getAllListings struct {
	ID          int32     `json:"id"`
	AdminID     int32     `json:"admin_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Location    string    `json:"location"`
	Available   bool      `json:"available"`
	Imagelink   []string  `json:"imagelink"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Server) GetAllFavorites(c *gin.Context) {
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

	favorites, err := s.q.GetFavorite(c, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var favoritesList []getAllListings
	for _, favorite := range favorites {
		listing, err := s.q.GetListingByID(c, favorite.ListingID)
		if err != nil {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		favoritesList = append(favoritesList, getAllListings{
			ID:          listing.ID,
			AdminID:     listing.AdminID,
			Title:       listing.Title,
			Description: listing.Description.String,
			Price:       listing.Price,
			Location:    listing.Location.String,
			Available:   listing.Available.Bool,
			Imagelink:   listing.Imagelinks,
			CreatedAt:   listing.CreatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, favoritesList)
}

func (s *Server) DeleteFavorite(c *gin.Context) {
	favoriteIDStr := c.Param("id")
	favoriteID, err := strconv.Atoi(favoriteIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized users only"})
		return
	}

	arg1 := db.GetFavoriteByListingIDParams{
		UserID:    user.ID,
		ListingID: int32(favoriteID),
	}

	_, err = s.q.GetFavoriteByListingID(c, arg1)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "favorite not found"})
		return
	}

	arg := db.DeleteFavoriteParams{
		ListingID: int32(favoriteID),
		UserID:    user.ID,
	}

	err = s.q.DeleteFavorite(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "favorite deleted successfully"})
}

type searchfavoriteResponse struct {
	ID          int32    `json:"id"`
	AdminID     int32    `json:"admin_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       string   `json:"price"`
	Location    string   `json:"location"`
	Available   bool     `json:"available"`
	Imagelink   []string `json:"imagelink"`
}

func (s *Server) SearchFavorite(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}

	searchString := sql.NullString{String: keyword, Valid: true}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	user, err := s.q.GetUser(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized users only"})
		return
	}

	arg := db.SearchFavoriteParams{
		Column1: searchString,
		UserID:  user.ID,
	}

	rows, err := s.q.SearchFavorite(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, []searchfavoriteResponse{})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var listings []searchfavoriteResponse
	for _, row := range rows {
		listings = append(listings, searchfavoriteResponse{
			ID:          row.ID,
			AdminID:     row.AdminID,
			Title:       row.Title,
			Description: row.Description.String,
			Price:       row.Price,
			Location:    row.Location.String,
			Available:   row.Available.Bool,
			Imagelink:   row.Imagelinks,
		})
	}

	c.JSON(http.StatusOK, listings)
}
