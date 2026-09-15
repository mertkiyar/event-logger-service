package main

import "testing"

func TestEventName_IsValid(t *testing.T) {
	event := Event{}
	event.EventType = "test-event-type"

	if event.EventType.IsValid() {
		t.Error("EventType.IsValid() should return false, because event type is invalid")
	}

	eventBlank := Event{}

	if eventBlank.EventType.IsValid() {
		t.Error("EventType.IsValid() should return false, because event type is empty")
	}

	eventCaseSensitive := Event{}
	eventCaseSensitive.EventType = "CLICK"
	if eventCaseSensitive.EventType.IsValid() {
		t.Error("EventType.IsValid() should return false, because event type is invalid")
	}

	// All Valid Event Types should be true
	eventLogin := Event{}
	eventLogin.EventType = EventLogin

	if !eventLogin.EventType.IsValid() {
		t.Error("EventType.IsValid() should return true, because event type is valid")
	}

	eventLogout := Event{}
	eventLogout.EventType = EventLogout

	if !eventLogout.EventType.IsValid() {
		t.Error("EventType.IsValid() should return true, because event type is valid")
	}

	eventSignup := Event{}
	eventSignup.EventType = EventSignup

	if !eventSignup.EventType.IsValid() {
		t.Error("EventType.IsValid() should return true, because event type is valid")
	}

	eventClick := Event{}
	eventClick.EventType = EventClick

	if !eventClick.EventType.IsValid() {
		t.Error("EventType.IsValid() should return true, because event type is valid")
	}
}
