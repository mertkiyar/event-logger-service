package main

type EventName string

const (
	EventLogin  EventName = "login"
	EventLogout EventName = "logout"
	EventSignup EventName = "signup"
	EventClick  EventName = "click"
)

type Event struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id,omitempty"` // omitempty = can be empty
	EventType EventName `json:"event_type"`
}

func (e EventName) IsValid() bool {
	switch e {
	case EventLogin, EventLogout, EventSignup, EventClick:
		return true
	default:
		return false
	}
}
