package schedule

import "github.com/gin-gonic/gin"

func roomName(eventID, eventDateID string) string {
	return eventID + "_" + eventDateID
}

func queueName(eventID, eventDateID string) string {
	return "queue:" + roomName(eventID, eventDateID)
}

func (s *Service) queueEmpty(c *gin.Context, name string) bool {
	n, err := s.cache.QueueLength(c.Request.Context(), name)
	return err != nil || n == 0
}
