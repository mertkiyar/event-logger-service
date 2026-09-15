package main

import "testing"

func TestEventName_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		eventName EventName
		want      bool
	}{
		{
			name:      "login is valid",
			eventName: EventLogin,
			want:      true,
		},
		{
			name:      "logout is valid",
			eventName: EventLogout,
			want:      true,
		},
		{
			name:      "signup is valid",
			eventName: EventSignup,
			want:      true,
		},
		{
			name:      "click is valid",
			eventName: EventClick,
			want:      true,
		},
		{
			name:      "empty event name is invalid",
			eventName: "",
			want:      false,
		},
		{
			name:      "unknown event name is invalid",
			eventName: "unknown",
			want:      false,
		},
		{
			name:      "uppercase event name is invalid",
			eventName: "CLICK",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.eventName.IsValid()

			if got != tt.want {
				t.Errorf(
					"EventName(%q).IsValid() = %v, valid %v",
					tt.eventName,
					got,
					tt.want,
				)
			}
		})
	}
}
