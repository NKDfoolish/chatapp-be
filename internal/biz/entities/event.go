package entities

import "time"

type Event struct {
	EventID   string
	RoomID    string
	Title     string
	Location  *string
	Capacity  *int
	StartTime time.Time
	EndTime   time.Time
	CreatedAt time.Time
	CreatedBy string
}

type EventParticipation struct {
	EventID string
	UserID  string
}

// Business methods
func (e *Event) IsActive() bool {
	now := time.Now()
	return now.After(e.StartTime) && now.Before(e.EndTime)
}

func (e *Event) IsUpcoming() bool {
	return time.Now().Before(e.StartTime)
}

func (e *Event) HasCapacityLimit() bool {
	return e.Capacity != nil && *e.Capacity > 0
}

func (e *Event) Duration() time.Duration {
	return e.EndTime.Sub(e.StartTime)
}

func (e *Event) IsPast() bool {
	return time.Now().After(e.EndTime)
}
