package schedule

import "context"

func (p *Postgres) Seats(ctx context.Context, dateID, areaID string) ([]Seat, error) {
	rows, err := p.db.QueryContext(ctx, seatSQL, areaID, dateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Seat{}
	for rows.Next() {
		var seat Seat
		if err := rows.Scan(&seat.ID, &seat.CX, &seat.CY, &seat.Row,
			&seat.Number, &seat.AreaID, &seat.ReservedUserID); err != nil {
			return nil, err
		}
		out = append(out, seat)
	}
	return out, rows.Err()
}
