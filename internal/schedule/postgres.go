package schedule

import (
	"context"
	"database/sql"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) Reservations(ctx context.Context, eventID, dateID string) ([]ReservationSeat, error) {
	rows, err := p.db.QueryContext(ctx, reservationSQL, eventID, dateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ReservationSeat{}
	for rows.Next() {
		var item ReservationSeat
		if err := rows.Scan(&item.SeatID); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (p *Postgres) Areas(ctx context.Context, eventID string) ([]Area, error) {
	rows, err := p.db.QueryContext(ctx, areaSQL, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Area{}
	for rows.Next() {
		var area Area
		if err := rows.Scan(&area.ID, &area.Label, &area.SVG, &area.Price); err != nil {
			return nil, err
		}
		out = append(out, area)
	}
	return out, rows.Err()
}
