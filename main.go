package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var events []Event

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Event Logger Service working good!")
	})
	http.HandleFunc("/events", createEventHandler)

	fmt.Println("Service working on 8080 port")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Service not started!", err)
	}
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {

	// in this if block, the service allow only post method
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post Method allowed!", http.StatusMethodNotAllowed)
		return
	}

	// create new event to convert json value from response
	var newEvent Event

	// read the json value in the request body and fill out newEvent
	// sending *newEvent because can write the data inside of variable
	err := json.NewDecoder(r.Body).Decode(&newEvent)
	if err != nil {
		http.Error(w, "Invalid JSON format!", http.StatusBadRequest)
		return
	}

	// append the new event in the events slice
	events = append(events, newEvent)

	// write the console 201 success http status code
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "The event saved successfuly, now total event number is: %d", len(events))
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {
	// in this if block, the service allow only get method
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get Method allowed!", http.StatusMethodNotAllowed)
		return
	}

	// return type is json
	w.Header().Set("Content-Type", "application/json")

	// convert the event slice to json format and write inside of response
	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		http.Error(w, "Failed to encode events!", http.StatusInternalServerError)
		return
	}
}
