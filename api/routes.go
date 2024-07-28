package api

import (
	"github.com/gin-gonic/gin"
	"github.com/weldonkipchirchir/rental_listing/middleware"
)

func (server *Server) initAdminRoutes(router *gin.Engine) {
	router.POST("/api/admin/register", server.CreateAdmin)
	router.POST("/api/admin/login", server.loginAdmin)
	router.POST("/api/admin/forgot-password", server.forgotAdminPassword)

	authRoutes := router.Group("/").Use(middleware.Authentication())
	authRoutes.PUT("/api/admin/update", server.updateAdmin)
	authRoutes.POST("/api/admin/logout", server.logoutAdmin)
}

func (server *Server) initUserRoutes(router *gin.Engine) {
	router.POST("/api/user/register", server.CreateUser)
	router.POST("/api/user/login", server.loginUser)
	router.POST("/api/user/forgot-password", server.forgotPassword)

	authRoutes := router.Group("/").Use(middleware.Authentication())
	authRoutes.PUT("/api/user/update", server.updateUser)
	authRoutes.POST("/api/user/logout", server.logoutUser)
}

func (server *Server) initListingRoutes(router *gin.Engine) {
	router.GET("/api/listing/", server.GetAllListings)
	router.GET("/api/listing/search", server.SearchListings)
	authRoutes := router.Group("/").Use(middleware.Authentication())
	authRoutes.GET("/api/listing/user", server.GetListings)
	authRoutes.POST("/api/listing/create", server.CreateListing)
	authRoutes.GET("/api/listing/:id", server.GetListingByID)
	authRoutes.GET("/api/listing/admin/listing", server.GetAdminListings)
	authRoutes.GET("/api/listing/admin/listing/data", server.GetAdminListingsData)
	authRoutes.GET("/api/listing/admin/listing/:id", server.GetListingsByAdminID)
	authRoutes.PUT("/api/listing/admin/listing/update/:id", server.UpdateListing)
	authRoutes.DELETE("/api/listing/admin/listing/:id", server.deleteListing)
	router.GET("/api/listings/:id/views", server.IncrementListingViews)
	authRoutes.PUT("/api/listing/listing/status/:id", server.UpdateListingStatus)
}

func (s *Server) initBookingRoutes(router *gin.Engine) {
	authRoutes := router.Group("/api").Use(middleware.Authentication())
	authRoutes.POST("/user/bookings", s.CreateBooking)
	authRoutes.GET("/user/bookings", s.GetBookingsByUserID)
	authRoutes.DELETE("/user/bookings/:id", s.DeleteBooking)

	authRoutes.GET("/admin/bookings", s.GetBookingsByAdminID)
	authRoutes.GET("/admin/bookings/:id", s.GetBookingByAdminIDAndID)
	authRoutes.PUT("/admin/bookings/:id", s.UpdateBookingStatus)
	authRoutes.PUT("/user/bookings/:id", s.updateCancelledBooking)
	authRoutes.GET("/admin/bookings/listing/:id", s.GetBookingsByListingID)
}

func (s *Server) initFavoriteRoutes(router *gin.Engine) {
	authRoutes := router.Group("/api").Use(middleware.Authentication())
	authRoutes.POST("/user/favorites", s.CreateFavorite)
	authRoutes.GET("/user/favorites", s.GetAllFavorites)
	authRoutes.DELETE("/user/favorite/:id", s.DeleteFavorite)
	authRoutes.GET("/favorite/search", s.SearchFavorite)
}

func (s *Server) initReviewRoutes(router *gin.Engine) {
	authRoutes := router.Group("/api").Use(middleware.Authentication())
	authRoutes.POST("/user/review", s.CreateReview)
	authRoutes.GET("/user/review/:id", s.GetReviewByID)
	authRoutes.GET("/user/review/listing", s.GetReviewsByListing)
	authRoutes.DELETE("/user/review/:id", s.DeleteReview)
}

func (s *Server) initNotificationRoutes(router *gin.Engine) {
	authRoutes := router.Group("/api").Use(middleware.Authentication())
	authRoutes.POST("/admin/notification", s.CreateNotification)
	authRoutes.POST("/user/notification", s.CreateAdminNotification)
	authRoutes.GET("/user/notification", s.GetNotifications)
	authRoutes.GET("/user/notification/sent", s.GetSentUserNotifications)
	authRoutes.GET("/admin/notification", s.GetAdminNotifications)
	authRoutes.GET("/admin/notification/sent", s.GetSentAdminNotifications)
	authRoutes.GET("/user/notification/:id", s.GetNotificationByID)
	authRoutes.PUT("/user/notification/:id", s.UpdateNotification)
	authRoutes.PUT("/admin/notification/:id", s.UpdateAdminNotification)
	authRoutes.GET("/user/notification/unread", s.GetUserUnreadNotifications)
	authRoutes.GET("/admin/notification/unread", s.GetAdminUnreadNotifications)
	authRoutes.DELETE("/user/notification/:id", s.DeleteNotification)
}

func (server *Server) initVerifyRoutes(router *gin.Engine) {
	router.GET("/api/user/verify/:email/:token", server.VerifyEmailUser)
	router.GET("/api/admin/verify/:email/:token", server.VerifyEmailAdmin)
}

func (s *Server) initPaymentRoutes(router *gin.Engine) {
	authRoutes := router.Group("/api").Use(middleware.Authentication())
	authRoutes.POST("/user/payment", s.HandleCreatePaymentIntent)
	authRoutes.GET("/user/config", s.Config)
	authRoutes.GET("/user/payment-records", s.GetPaymentsByUserID)
	authRoutes.GET("/admin/payment-records", s.GetPaymentsByAdminID)
}

func (s *Server) initStatsRoutes(router *gin.Engine) {
	authRoutes := router.Group("/api").Use(middleware.Authentication())
	authRoutes.GET("/admin/stats", s.GetAllStats)
}
