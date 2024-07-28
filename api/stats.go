package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Stat represents a statistics entity.
type Stat struct {
	ID            int32    `json:"id"`
	Title         string   `json:"title"`
	TotalViews    int32    `json:"total_views"`
	TotalBookings int64    `json:"total_bookings"`
	AverageRating *float64 `json:"average_rating"`
	Income        *float64 `json:"income"`
}

// GetAllStats retrieves all stat entries.
func (s *Server) GetAllStats(c *gin.Context) {
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

	stats, err := s.q.Statistics(c, admin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	statistics := make([]Stat, 0, len(stats))
	for _, stat := range stats {
		var avgRating *float64
		if stat.AverageRating != nil {
			switch v := stat.AverageRating.(type) {
			case float64:
				avgRating = &v
			case []uint8:
				f, err := strconv.ParseFloat(string(v), 64)
				if err == nil {
					avgRating = &f
				}
			}
		}
		if avgRating == nil {
			defaultValue := 0.0
			avgRating = &defaultValue
		}
		var income *float64
		if stat.TotalConfirmedAmount != nil {
			switch v := stat.TotalConfirmedAmount.(type) {
			case float64:
				income = &v
			case []uint8:
				f, err := strconv.ParseFloat(string(v), 64)
				if err == nil {
					income = &f
				}
			}
		}
		statistics = append(statistics, Stat{
			ID:            stat.ID,
			Title:         stat.Title,
			TotalViews:    stat.TotalViews.Int32,
			TotalBookings: stat.TotalBookings,
			AverageRating: avgRating,
			Income:        income,
		})
	}

	c.JSON(http.StatusOK, statistics)
}
