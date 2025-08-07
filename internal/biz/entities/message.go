package entities

import "time"

type MessageType string

const (
	MessageTypeText  MessageType = "TEXT"
	MessageTypeImage MessageType = "IMAGE"
	MessageTypeEvent MessageType = "EVENT"
)

type Message struct {
	MessageID   string
	RoomID      string
	SenderID    string
	Content     *string
	MessageType MessageType
	MediaURL    *string
	SentAt      time.Time
	IsRead      bool
}

// Business methods
func (m *Message) MarkAsRead() {
	m.IsRead = true
}

func (m *Message) HasMedia() bool {
	return m.MediaURL != nil && *m.MediaURL != ""
}

func (m *Message) IsTextMessage() bool {
	return m.MessageType == MessageTypeText
}

func (m *Message) IsSentBy(userID string) bool {
	return m.SenderID == userID
}

func (m *Message) IsRecentlySent(minutes int) bool {
	threshold := time.Now().Add(-time.Duration(minutes) * time.Minute)
	return m.SentAt.After(threshold)
}
