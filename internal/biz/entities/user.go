package entities

import "time"

type User struct {
	UserID         string
	Username       string
	Email          string
	Password       string
	IsOnline       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ProfilePicture *string
}

// Business methods
func (u *User) IsActive() bool {
	return u.IsOnline
}

func (u *User) SetOnlineStatus(online bool) {
	u.IsOnline = online
	u.UpdatedAt = time.Now()
}

func (u *User) ValidateEmail() bool {
	// Business logic for email validation
	return len(u.Email) > 0 && len(u.Username) > 0
}
