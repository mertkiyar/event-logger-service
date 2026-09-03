package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

var events []Event
var eventChan = make(chan Event, 10000)

func eventWorker() {
	for event := range eventChan {
		events = append(events, event)
		fmt.Println("new event processed:", event.UserId, "->", event.EventType)
	}
}

// mutex provides write data to map at the same time without concurrent map writes
var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

type visitor struct {
	count    int
	lastSeen time.Time
}

func main() {
	go eventWorker()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Event Logger Service working good!")
	})

	http.HandleFunc("/events", rateLimiter(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createEventHandler(w, r)
		case http.MethodGet:
			getEventsHandler(w, r)
		default:
			http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		}
	}))

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

func rateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//excludes port from ip address
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP address!", http.StatusBadRequest)
			return
		}

		// lock mutex to safely access the visitors map
		mu.Lock()
		v, exists := visitors[ip]

		// if visitor does not exist, create new visitor and get ip address
		if !exists {
			v = &visitor{count: 0, lastSeen: time.Now()}
			visitors[ip] = v
		}

		// reset visit count if 10 seconds have passed since lastSeen
		if time.Since(v.lastSeen) > time.Second*10 {
			v.count = 0
		}

		// update lastSeen and increment visit count
		v.lastSeen = time.Now()
		v.count++

		// if request count exceeds 100, throw an error
		if v.count >= 100 {
			mu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		mu.Unlock()

		//route the w and r value to create or get func
		next(w, r)
	}
}
