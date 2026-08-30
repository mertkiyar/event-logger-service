package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var events []Event
var eventChan = make(chan Event, 10000)

func eventWorker() {
	for event := range eventChan {
		events = append(events, event)
		fmt.Println("new event processed:", event.UserId, "->", event.EventType)
	}
}

func main() {
	go eventWorker()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Event Logger Service working good!")
	})

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createEventHandler(w, r)
		case http.MethodGet:
			getEventsHandler(w, r)
		default:
			http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/stats", getStatsHandler)

	fmt.Println("Service working on 8080 port")

	go func() {
		sendRequest()
		sendRequestWithSpecificNumber(10000)
	}()

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Service not started!", err)
	}
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {

	// create new event to convert json value from response
	var newEvent Event

	// read the json value in the request body and fill out newEvent
	// sending *newEvent because can write the data inside variable
	err := json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		http.Error(w, "Invalid JSON format!", http.StatusBadRequest)
		return
	}

	// if the event type is not predefined to new event, no new events are created
	// only defined types like login, click and signup allowed
	if !newEvent.EventType.IsValid() {
		http.Error(w, "Invalid event type!", http.StatusBadRequest)
		return
	}

	// throw the new event data inside channel
	eventChan <- newEvent

	// write the console 201 success http status code
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "The event has been processed")
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {

	// return type is json
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		if events != nil {
			// convert the event slice to json format and write inside response
			err := json.NewEncoder(w).Encode(events)
			if err != nil {
				http.Error(w, "Failed to encode events!", http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Fprintln(w, events)
		}
		return
	}

	filteredEvents := []Event{}

	for _, event := range events {
		if event.UserId == userID {
			filteredEvents = append(filteredEvents, event)
		}
	}

	err := json.NewEncoder(w).Encode(filteredEvents)
	if err != nil {
		http.Error(w, "Failed to encode events!", http.StatusInternalServerError)
		return
	}
}
func getStatsHandler(w http.ResponseWriter, r *http.Request) {
	// in this if block, the service allow only get method
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get Method allowed!", http.StatusMethodNotAllowed)
		return
	}

	// return type is json
	w.Header().Set("Content-Type", "application/json")

	// event type and number of this event
	stats := make(map[string]int)

	for _, event := range events {
		stats[string(event.EventType)]++
	}

	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, "Failed to encode stats!", http.StatusInternalServerError)
		return
	}
}
