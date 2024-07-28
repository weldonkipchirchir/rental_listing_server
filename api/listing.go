package api

import (
	"context"
	"database/sql"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
)

// Define request and response structs
type createListingRequest struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Price       string `json:"price" form:"price"`
	Location    string `json:"location" form:"location"`
	Available   bool   `json:"available" form:"available"`
}

type createListingResponse struct {
	ID          int32     `json:"id"`
	AdminID     int32     `json:"admin_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Location    string    `json:"location"`
	Available   bool      `json:"available"`
	Imagelinks  []string  `json:"imagelink"`
	CreatedAt   time.Time `json:"created_at"`
}

// Create a helper function to upload multiple images
func UploadMultipleToCloudinary(c *gin.Context, files []*multipart.FileHeader) ([]string, error) {
	cloudinaryURL := os.Getenv("cloudnary")
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		return nil, err
	}

	var imageURLs []string

	for _, file := range files {
		localPath := "assets/uploads/" + file.Filename
		if err := os.MkdirAll("assets/uploads", os.ModePerm); err != nil {
			return nil, err
		}
		if err := c.SaveUploadedFile(file, localPath); err != nil {
			return nil, err
		}
		defer os.Remove(localPath)

		ctx := context.Background()
		resp, err := cld.Upload.Upload(ctx, localPath, uploader.UploadParams{PublicID: "my_avatar-" + file.Filename + "-" + uuid.New().String()})
		if err != nil {
			return nil, err
		}

		imageURLs = append(imageURLs, resp.SecureURL)
	}

	return imageURLs, nil
}

func (s *Server) CreateListing(c *gin.Context) {
	var req createListingRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized admins only"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one image is required"})
		return
	}

	imageURLs, err := UploadMultipleToCloudinary(c, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload images"})
		return
	}

	arg := db.CreateListingParams{
		AdminID:     admin.ID,
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Price:       req.Price,
		Location:    sql.NullString{String: req.Location, Valid: req.Location != ""},
		Available:   sql.NullBool{Bool: req.Available, Valid: true},
		Column7:     imageURLs,
	}

	listing, err := s.q.CreateListing(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := createListingResponse{
		ID:          listing.ID,
		AdminID:     listing.AdminID,
		Title:       listing.Title,
		Description: listing.Description.String,
		Price:       listing.Price,
		Location:    listing.Location.String,
		Available:   listing.Available.Bool,
		Imagelinks:  listing.Imagelinks,
		CreatedAt:   listing.CreatedAt.Time,
	}
	c.JSON(http.StatusCreated, rsp)
}

type getAllListingsResponse struct {
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

func (s *Server) GetAllListings(c *gin.Context) {
	var listings []db.Listing
	rows, err := s.q.GetListings(c)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listings = make([]db.Listing, len(rows))
	for i, row := range rows {
		listings[i] = db.Listing{
			ID:          row.ID,
			AdminID:     row.AdminID,
			Title:       row.Title,
			Description: sql.NullString{String: row.Description.String, Valid: row.Description.Valid},
			Price:       row.Price,
			Location:    sql.NullString{String: row.Location.String, Valid: row.Location.Valid},
			Available:   sql.NullBool{Bool: row.Available.Bool, Valid: row.Available.Valid},
			Imagelinks:  row.Imagelinks,
			CreatedAt:   row.CreatedAt,
		}
	}

	response := make([]getAllListingsResponse, len(listings))
	for i, listing := range listings {
		response[i] = getAllListingsResponse{
			ID:          listing.ID,
			AdminID:     listing.AdminID,
			Title:       listing.Title,
			Description: listing.Description.String,
			Price:       listing.Price,
			Location:    listing.Location.String,
			Available:   listing.Available.Bool,
			Imagelink:   listing.Imagelinks,
			CreatedAt:   listing.CreatedAt.Time,
		}
	}

	c.JSON(http.StatusOK, response)
}

type getFavoriteListingsResponse struct {
	ID             int32     `json:"id"`
	AdminID        int32     `json:"admin_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Price          string    `json:"price"`
	Location       string    `json:"location"`
	Available      bool      `json:"available"`
	Imagelink      []string  `json:"imagelink"`
	BookmarkStatus bool      `json:"bookmarkStatus"`
	CreatedAt      time.Time `json:"created_at"`
}

func (s *Server) GetListings(c *gin.Context) {
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

	var listings []db.Listing
	rows, err := s.q.GetListings(c)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listings = make([]db.Listing, len(rows))
	for i, row := range rows {
		listings[i] = db.Listing{
			ID:          row.ID,
			AdminID:     row.AdminID,
			Title:       row.Title,
			Description: sql.NullString{String: row.Description.String, Valid: row.Description.Valid},
			Price:       row.Price,
			Location:    sql.NullString{String: row.Location.String, Valid: row.Location.Valid},
			Available:   sql.NullBool{Bool: row.Available.Bool, Valid: row.Available.Valid},
			Imagelinks:  row.Imagelinks,
			CreatedAt:   row.CreatedAt,
		}
	}

	var bookmark bool

	response := make([]getFavoriteListingsResponse, len(listings))
	for i, listing := range listings {
		arg := db.GetListingFavoriteByUserParams{
			UserID:    user.ID,
			ListingID: listing.ID,
		}
		favorite, err := s.q.GetListingFavoriteByUser(c, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				bookmark = false
			} else {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
		if favorite.ID != 0 {
			bookmark = true
		}
		response[i] = getFavoriteListingsResponse{
			ID:             listing.ID,
			AdminID:        listing.AdminID,
			Title:          listing.Title,
			Description:    listing.Description.String,
			Price:          listing.Price,
			Location:       listing.Location.String,
			Available:      listing.Available.Bool,
			Imagelink:      listing.Imagelinks,
			BookmarkStatus: bookmark,
			CreatedAt:      listing.CreatedAt.Time,
		}
	}

	c.JSON(http.StatusOK, response)
}

type getAdminListingsResponse struct {
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

func (s *Server) GetAdminListings(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	// Try to get data from cache
	var listings []db.Listing

	rows, err := s.q.GetAdminListings(c, admin.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listings = make([]db.Listing, len(rows))
	for i, row := range rows {
		listings[i] = db.Listing{
			ID:          row.ID,
			AdminID:     row.AdminID,
			Title:       row.Title,
			Description: sql.NullString{String: row.Description.String, Valid: row.Description.Valid},
			Price:       row.Price,
			Location:    sql.NullString{String: row.Location.String, Valid: row.Location.Valid},
			Available:   sql.NullBool{Bool: row.Available.Bool, Valid: row.Available.Valid},
			Imagelinks:  row.Imagelinks,
			CreatedAt:   row.CreatedAt,
		}
	}

	response := make([]getAdminListingsResponse, len(listings))
	for i, listing := range listings {
		response[i] = getAdminListingsResponse{
			ID:          listing.ID,
			AdminID:     listing.AdminID,
			Title:       listing.Title,
			Description: listing.Description.String,
			Price:       listing.Price,
			Location:    listing.Location.String,
			Available:   listing.Available.Bool,
			Imagelink:   listing.Imagelinks,
			CreatedAt:   listing.CreatedAt.Time,
		}
	}

	c.JSON(http.StatusOK, response)
}

type getAdminListingsDataResponse struct {
	ID          int32     `json:"id"`
	AdminID     int32     `json:"admin_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Location    string    `json:"location"`
	Available   bool      `json:"available"`
	Imagelink   []string  `json:"imagelink"`
	Revenue     float64   `json:"revenue"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Server) GetAdminListingsData(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	var listings []getAdminListingsDataResponse

	rows, err := s.q.GetAdminListings(c, admin.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listings = make([]getAdminListingsDataResponse, len(rows))
	for i, row := range rows {
		listings[i] = getAdminListingsDataResponse{
			ID:          row.ID,
			AdminID:     row.AdminID,
			Title:       row.Title,
			Description: row.Description.String,
			Price:       row.Price,
			Location:    row.Location.String,
			Available:   row.Available.Bool,
			Imagelink:   row.Imagelinks,
			CreatedAt:   row.CreatedAt.Time,
		}
	}

	response := make([]getAdminListingsDataResponse, len(listings))
	for i, listing := range listings {
		var totalRevenue float64
		arg := db.GetBookingsByListingIDParams{
			ListingID: listing.ID,
			AdminID:   admin.ID,
		}
		booking, err := s.q.GetBookingsByListingID(c, arg)
		if err != nil {
			if err == sql.ErrNoRows {
				totalRevenue = 0
			} else {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			for _, b := range booking {
				price, err := strconv.ParseFloat(b.TotalAmount, 32)
				if err != nil {
					c.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
				if b.Status.String == "completed" || b.Status.String == "confirmed" {
					totalRevenue += price
				} else {
					totalRevenue += 0
				}
			}
		}

		response[i] = getAdminListingsDataResponse{
			ID:          listing.ID,
			AdminID:     listing.AdminID,
			Title:       listing.Title,
			Description: listing.Description,
			Price:       listing.Price,
			Location:    listing.Location,
			Available:   listing.Available,
			Imagelink:   listing.Imagelink,
			Revenue:     totalRevenue,
			CreatedAt:   listing.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func newListingAdminByIdResponse(listing *db.Listing) getListingsResponse {
	return getListingsResponse{
		ID:          listing.ID,
		AdminID:     listing.AdminID,
		Title:       listing.Title,
		Description: listing.Description.String,
		Price:       listing.Price,
		Location:    listing.Location.String,
		Available:   listing.Available.Bool,
		Imagelink:   listing.Imagelinks,
		CreatedAt:   listing.CreatedAt.Time,
	}
}

type getListingsResponse struct {
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

func (s *Server) GetListingsByAdminID(c *gin.Context) {
	jobId := c.Param("id")
	id, err := strconv.Atoi(jobId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	arg := db.GetListingsByAdminIDParams{
		AdminID: admin.ID,
		ID:      int32(id),
	}

	row, err := s.q.GetListingsByAdminID(c, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	listing := db.Listing{
		ID:          row.ID,
		AdminID:     row.AdminID,
		Title:       row.Title,
		Description: sql.NullString{String: row.Description.String, Valid: row.Description.Valid},
		Price:       row.Price,
		Location:    sql.NullString{String: row.Location.String, Valid: row.Location.Valid},
		Available:   sql.NullBool{Bool: row.Available.Bool, Valid: row.Available.Valid},
		Imagelinks:  row.Imagelinks,
		CreatedAt:   row.CreatedAt,
	}

	res := newListingAdminByIdResponse(&listing)

	c.JSON(http.StatusOK, res)
}

type Review struct {
	ID        int32  `json:"id"`
	Rating    int32  `json:"rating"`
	Username  string `json:"username"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type getListingsByIDResponse struct {
	ID          int32     `json:"id"`
	AdminID     int32     `json:"admin_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Location    string    `json:"location"`
	Available   bool      `json:"available"`
	Imagelink   []string  `json:"imagelink"`
	CreatedAt   time.Time `json:"created_at"`
	Reviews     []Review  `json:"reviews"`
}

func (s *Server) GetListingByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	row, err := s.q.GetListingByID(c, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rows, err := s.q.GetListingReviews(c, row.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reviews := make([]Review, len(rows))
	for i, row := range rows {
		reviews[i] = Review{
			ID:        row.ID,
			Username:  row.Username,
			Rating:    row.Rating,
			Comment:   row.Comment.String,
			CreatedAt: row.CreatedAt.Time.String(),
		}
	}

	listing := getListingsByIDResponse{
		ID:          row.ID,
		AdminID:     row.AdminID,
		Title:       row.Title,
		Description: row.Description.String,
		Price:       row.Price,
		Location:    row.Location.String,
		Available:   row.Available.Bool,
		Imagelink:   row.Imagelinks,
		CreatedAt:   row.CreatedAt.Time,
		Reviews:     reviews,
	}

	c.JSON(http.StatusOK, listing)
}

func (s *Server) UpdateListing(c *gin.Context) {
	listingIDStr := c.Param("id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	arg1 := db.GetListingsByAdminIDParams{
		AdminID: admin.ID,
		ID:      int32(listingID),
	}

	_, err = s.q.GetListingsByAdminID(c, arg1)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "listing not found"})
		return
	}

	arg := db.UpdateListingParams{
		ID:      int32(listingID),
		AdminID: admin.ID,
	}

	// Handle image upload if new images are provided
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Access form fields
	title := form.Value["title"][0]
	description := form.Value["description"][0]
	price := form.Value["price"][0]
	location := form.Value["location"][0]
	available := form.Value["available"][0]

	// Optional fields
	if title != "" {
		arg.Title = sql.NullString{String: title, Valid: true}
	}
	if price != "" {
		arg.Price = sql.NullString{String: price, Valid: true}
	}
	if description != "" {
		arg.Description = sql.NullString{String: description, Valid: true}
	}
	if location != "" {
		arg.Location = sql.NullString{String: location, Valid: true}
	}
	if available != "" {
		availableValue, _ := strconv.ParseBool(form.Value["available"][0])
		arg.Available = sql.NullBool{Bool: availableValue, Valid: true}
	}

	files := form.File["images"]
	if len(files) > 0 {
		imageURLs, err := UploadMultipleToCloudinary(c, files)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload images"})
			return
		}

		arg.Imagelinks = imageURLs
	}

	err = s.q.UpdateListing(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "listing updated successfully"})
}

func (s *Server) deleteListing(c *gin.Context) {
	listingIDStr := c.Param("id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not found"})
		return
	}

	admin, err := s.q.GetAdmin(c, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized admins only"})
		return
	}

	arg1 := db.GetListingsByAdminIDParams{
		AdminID: admin.ID,
		ID:      int32(listingID),
	}

	_, err = s.q.GetListingsByAdminID(c, arg1)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "listing not found"})
		return
	}

	arg := db.DeleteListingParams{
		ID:      int32(listingID),
		AdminID: admin.ID,
	}

	arg2 := db.ListingActiveBookingCountParams{
		ListingID: int32(listingID),
		AdminID:   admin.ID,
	}

	count, err := s.q.ListingActiveBookingCount(c, arg2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "listing has active bookings"})
		return
	}

	err = s.q.DeleteListing(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "listing deleted successfully"})
}

func (s *Server) IncrementListingViews(c *gin.Context) {
	listingIDStr := c.Param("id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Update view count in the stats table
	err = s.q.UpdateTotalViews(c, int32(listingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "view count incremented"})
}

type updateListingStatusRequest struct {
	Available bool `json:"available"`
}

func (s *Server) UpdateListingStatus(c *gin.Context) {
	listingIDStr := c.Param("id")
	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var req updateListingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	listing, err := s.q.GetListingByID(c, int32(listingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateListingStatusParams{
		ID:        listing.ID,
		Available: sql.NullBool{Bool: req.Available, Valid: true},
	}

	err = s.q.UpdateListingStatus(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "listing status updated"})
}

type searchListingsResponse struct {
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

func (s *Server) SearchListings(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}

	searchString := sql.NullString{String: keyword, Valid: true}

	rows, err := s.q.SearchListings(c, searchString)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "no listings found"})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var listings []searchListingsResponse
	for _, row := range rows {
		listings = append(listings, searchListingsResponse{
			ID:          row.ID,
			AdminID:     row.AdminID,
			Title:       row.Title,
			Description: row.Description.String,
			Price:       row.Price,
			Location:    row.Location.String,
			Available:   row.Available.Bool,
			Imagelink:   row.Imagelinks,
			CreatedAt:   row.CreatedAt.Time,
		})
	}

	c.JSON(http.StatusOK, listings)
}
