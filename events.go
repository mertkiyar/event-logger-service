package main

type EventName string

const (
	EventLogin  EventName = "login"
	EventSignup EventName = "signup"
	EventClick  EventName = "click"
)

type Event struct {
	UserId    string    `json:"user_id"`
	EventType EventName `json:"event_type"`
}
