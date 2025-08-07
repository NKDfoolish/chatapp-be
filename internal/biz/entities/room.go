package entities

import "time"

type Room struct {
	RoomID        string
	LastMessageAt *time.Time
	CreatedAt     time.Time
	CreatedBy     string
}

type RoomUser struct {
	RoomID string
	UserID string
}

// Business methods
func (r *Room) UpdateLastMessage() {
	now := time.Now()
	r.LastMessageAt = &now
}

func (r *Room) IsCreatedBy(userID string) bool {
	return r.CreatedBy == userID
}

func (r *Room) HasRecentActivity(hours int) bool {
	if r.LastMessageAt == nil {
		return false
	}
	threshold := time.Now().Add(-time.Duration(hours) * time.Hour)
	return r.LastMessageAt.After(threshold)
}

func (r *Room) IsActive() bool {
	return r.HasRecentActivity(24) // active if message within 24h
}
