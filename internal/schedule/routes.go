package schedule

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Service) ReservationStatus(c *gin.Context) {
	eventID, dateID, ok := s.request(c)
	if !ok {
		return
	}
	name := queueName(eventID, dateID)
	job := "reservation:status:" + name
	if !s.startJob(job, 5*time.Second, func() bool {
		if s.queueEmpty(c, name) {
			return false
		}
		seats, err := s.store.Reservations(c.Request.Context(), eventID, dateID)
		if err == nil {
			_ = s.cache.Publish(c.Request.Context(), queueBroadcast, map[string]any{
				"room": name, "type": "seatsInfo", "payload": map[string]any{"seatsInfo": seats},
			})
		}
		return err == nil
	}) {
		conflict(c)
		return
	}
	created(c, job)
}

func (s *Service) QueueStatus(c *gin.Context) {
	eventID, dateID, ok := s.request(c)
	if !ok {
		return
	}
	name := queueName(eventID, dateID)
	job := "queue:status:" + name
	if !s.startJob(job, time.Second, func() bool {
		if s.queueEmpty(c, name) {
			return false
		}
		_ = s.cache.Publish(c.Request.Context(), queueBroadcast, map[string]any{
			"room": name, "type": "queueStatus",
		})
		return true
	}) {
		conflict(c)
		return
	}
	created(c, job)
}

func (s *Service) SeatReservationStatistic(c *gin.Context) {
	eventID, dateID, ok := s.request(c)
	if !ok {
		return
	}
	room := roomName(eventID, dateID)
	job := "seat:reservation:statistic:" + room
	if !s.startJob(job, 2*time.Second, func() bool {
		count, err := s.cache.RoomCount(c.Request.Context(), room)
		if err != nil || count == 0 {
			return false
		}
		return s.publishAreaStats(c, eventID, dateID, room) == nil
	}) {
		conflict(c)
		return
	}
	created(c, job)
}
