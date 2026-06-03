package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeService struct{}

func (fakeService) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
func (fakeService) ReservationStatus(c *gin.Context) {
	c.Status(http.StatusCreated)
}
func (fakeService) QueueStatus(c *gin.Context) {
	c.Status(http.StatusAccepted)
}
func (fakeService) SeatReservationStatistic(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func TestRouterLiveness(t *testing.T) {
	r := NewRouter(Config{CORSOrigins: []string{"*"}}, fakeService{})
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/scheduling/liveness", nil))
	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.Code)
	}
}

func TestRouterShortLiveness(t *testing.T) {
	r := NewRouter(Config{CORSOrigins: []string{"*"}}, fakeService{})
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/liveness", nil))
	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.Code)
	}
}

func TestRouterPostRoutes(t *testing.T) {
	r := NewRouter(Config{CORSOrigins: []string{"*"}}, fakeService{})
	tests := map[string]int{
		"/scheduling/reservation/status":         http.StatusCreated,
		"/scheduling/queue/status":               http.StatusAccepted,
		"/scheduling/seat/reservation/statistic": http.StatusNoContent,
	}
	for path, want := range tests {
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, httptest.NewRequest(http.MethodPost, path, nil))
		if resp.Code != want {
			t.Fatalf("%s status=%d want=%d", path, resp.Code, want)
		}
	}
}
