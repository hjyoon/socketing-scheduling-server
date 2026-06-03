package schedule

import "context"

const (
	queueBroadcast       = "socketing:queue:broadcast"
	reservationBroadcast = "socketing:reservation:broadcast"
)

type Config struct {
	JWTSecret string
}

type Area struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	SVG   string `json:"svg,omitempty"`
	Price int    `json:"price"`
}

type Seat struct {
	ID             string  `json:"id"`
	CX             float64 `json:"cx"`
	CY             float64 `json:"cy"`
	Row            int     `json:"row"`
	Number         int     `json:"number"`
	AreaID         string  `json:"area_id"`
	SelectedBy     *string `json:"selectedBy"`
	ReservedUserID *string `json:"reservedUserId"`
	UpdatedAt      *string `json:"updatedAt"`
	ExpirationTime *string `json:"expirationTime"`
}

type ReservationSeat struct {
	SeatID string `json:"seat_id"`
}

type AreaStat struct {
	AreaID           string `json:"areaId"`
	TotalSeatsNum    int    `json:"totalSeatsNum"`
	ReservedSeatsNum int    `json:"reservedSeatsNum"`
}

type Cache interface {
	Ready(context.Context) error
	Publish(context.Context, string, any) error
	Areas(context.Context, string) ([]Area, error)
	SetArea(context.Context, string, Area) error
	Seats(context.Context, string) ([]Seat, error)
	SetSeat(context.Context, string, Seat) error
	QueueLength(context.Context, string) (int, error)
	RoomCount(context.Context, string) (int, error)
}

type Store interface {
	Reservations(context.Context, string, string) ([]ReservationSeat, error)
	Areas(context.Context, string) ([]Area, error)
	Seats(context.Context, string, string) ([]Seat, error)
}
