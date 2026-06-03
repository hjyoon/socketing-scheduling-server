package schedule

import (
	"context"
	"encoding/json"
)

func (r *Redis) Areas(ctx context.Context, room string) ([]Area, error) {
	raw, err := r.c.HVals(ctx, "areas:"+room).Result()
	if err != nil {
		return nil, err
	}
	out := []Area{}
	for _, item := range raw {
		var area Area
		if json.Unmarshal([]byte(item), &area) == nil {
			out = append(out, area)
		}
	}
	return out, nil
}

func (r *Redis) SetArea(ctx context.Context, room string, area Area) error {
	raw, err := json.Marshal(area)
	if err != nil {
		return err
	}
	return r.c.HSet(ctx, "areas:"+room, area.ID, raw).Err()
}

func (r *Redis) Seats(ctx context.Context, area string) ([]Seat, error) {
	raw, err := r.c.HVals(ctx, "seats:"+area).Result()
	if err != nil {
		return nil, err
	}
	out := []Seat{}
	for _, item := range raw {
		var seat Seat
		if json.Unmarshal([]byte(item), &seat) == nil {
			out = append(out, seat)
		}
	}
	return out, nil
}

func (r *Redis) SetSeat(ctx context.Context, area string, seat Seat) error {
	raw, err := json.Marshal(seat)
	if err != nil {
		return err
	}
	return r.c.HSet(ctx, "seats:"+area, seat.ID, raw).Err()
}
