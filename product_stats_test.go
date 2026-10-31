package main

import "testing"

func TestCountProductClicks(t *testing.T) {
	items := []Event{
		{UserID: "u1", ProductID: "p1", EventType: EventClick},
		{UserID: "u1", ProductID: "p1", EventType: EventClick},
		{UserID: "u1", ProductID: "p2", EventType: EventClick},
		{UserID: "u2", ProductID: "p1", EventType: EventClick},
		{UserID: "u1", ProductID: "p1", EventType: EventLogin},
	}

	tests := []struct {
		name      string
		userID    string
		productID string
		want      int
	}{
		{"matching user and product", "u1", "p1", 2},
		{"no matching product", "u1", "missing", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countProductClicks(items, tt.userID, tt.productID)
			if got != tt.want {
				t.Errorf("click count = %d, want %d", got, tt.want)
			}
		})
	}
}
