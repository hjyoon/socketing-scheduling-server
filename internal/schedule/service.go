package schedule

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hjyoon/socketing-scheduling-server/internal/auth"
)

type Service struct {
	cfg   Config
	cache Cache
	store Store
	mu    sync.Mutex
	jobs  map[string]contextCancel
}

type contextCancel func()

func NewService(cfg Config, cache Cache, store Store) *Service {
	return &Service{cfg: cfg, cache: cache, store: store, jobs: map[string]contextCancel{}}
}

func (s *Service) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "The server is alive."})
}

func (s *Service) request(c *gin.Context) (string, string, bool) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	claims, err := auth.Verify(token, s.cfg.JWTSecret)
	if err != nil || claims["sub"] != "scheduling" {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Invalid token subject. Unauthorized request."})
		return "", "", false
	}
	eventID, _ := claims["eventId"].(string)
	dateID, _ := claims["eventDateId"].(string)
	if eventID == "" || dateID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing required token payload: eventId or eventDateId."})
		return "", "", false
	}
	return eventID, dateID, true
}

func (s *Service) startJob(name string, interval time.Duration, tick func() bool) bool {
	s.mu.Lock()
	if s.jobs[name] != nil {
		s.mu.Unlock()
		return false
	}
	stop := make(chan struct{})
	s.jobs[name] = func() { close(stop) }
	s.mu.Unlock()
	go s.runJob(name, interval, tick, stop)
	return true
}

func (s *Service) runJob(name string, interval time.Duration, tick func() bool, stop <-chan struct{}) {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-stop:
			return
		case <-timer.C:
			if !tick() {
				s.mu.Lock()
				delete(s.jobs, name)
				s.mu.Unlock()
				return
			}
			timer.Reset(interval)
		}
	}
}

func created(c *gin.Context, job string) {
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Cron job created successfully.", "data": gin.H{"jobName": job}})
}

func conflict(c *gin.Context) {
	c.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "A cron job is already running for the provided event."})
}
