package app

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type SchedulingService interface {
	Liveness(*gin.Context)
	ReservationStatus(*gin.Context)
	QueueStatus(*gin.Context)
	SeatReservationStatistic(*gin.Context)
}

func NewRouter(cfg Config, service SchedulingService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))
	r.GET("/scheduling/liveness", service.Liveness)
	r.GET("/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "The server is alive."})
	})
	r.POST("/scheduling/reservation/status", service.ReservationStatus)
	r.POST("/scheduling/queue/status", service.QueueStatus)
	r.POST("/scheduling/seat/reservation/statistic", service.SeatReservationStatistic)
	return r
}
