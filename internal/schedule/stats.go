package schedule

import "github.com/gin-gonic/gin"

func (s *Service) publishAreaStats(c *gin.Context, eventID, dateID, room string) error {
	areas, err := s.cache.Areas(c.Request.Context(), room)
	if err != nil {
		return err
	}
	if len(areas) == 0 {
		areas, err = s.store.Areas(c.Request.Context(), eventID)
		if err != nil {
			return err
		}
		for _, area := range areas {
			_ = s.cache.SetArea(c.Request.Context(), room, area)
		}
	}
	stats := []AreaStat{}
	for _, area := range areas {
		seats, err := s.seats(c, eventID, dateID, area.ID)
		if err != nil {
			return err
		}
		stats = append(stats, areaStat(area.ID, seats))
	}
	return s.cache.Publish(c.Request.Context(), reservationBroadcast, map[string]any{
		"room": room, "type": "reservedSeatsStatistic", "payload": stats,
	})
}

func (s *Service) seats(c *gin.Context, eventID, dateID, areaID string) ([]Seat, error) {
	name := eventID + "_" + dateID + "_" + areaID
	seats, err := s.cache.Seats(c.Request.Context(), name)
	if err != nil {
		return nil, err
	}
	if len(seats) > 0 {
		return seats, nil
	}
	seats, err = s.store.Seats(c.Request.Context(), dateID, areaID)
	if err != nil {
		return nil, err
	}
	for _, seat := range seats {
		_ = s.cache.SetSeat(c.Request.Context(), name, seat)
	}
	return seats, nil
}

func areaStat(areaID string, seats []Seat) AreaStat {
	stat := AreaStat{AreaID: areaID, TotalSeatsNum: len(seats)}
	for _, seat := range seats {
		if seat.ReservedUserID != nil {
			stat.ReservedSeatsNum++
		}
	}
	return stat
}
