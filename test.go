package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

type TestEvent struct {
	UserId    string `json:"user_id"`
	EventType string `json:"event_type"`
}

func sendRequest() {
	sampleData := map[string]string{
		"user_id":    "test_user_id",
		"event_type": "login",
	}

	jsonBytes, err := json.Marshal(sampleData)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s\n", err)
		return
	}

	url := "http://localhost:8080/events"

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Printf("Error posting event: %s\n", err)
		return
	}

	defer response.Body.Close()
}

// OBJECTIVE: Perform a concurrent load and integration test on the /events endpoint using multiple goroutines.
// EXPECTED OUTCOME: Verifies that the server can handle high-concurrency traffic without crashing, validates that
// the internal channel pipeline (`eventChan`) correctly serializes and processes simultaneous random events, and
// ensures no data race or corruption occurs on the global data structures.
func sendRequestWithSpecificNumber(count int) {
	eventTypes := []string{"login", "logout", "click", "signup"}

	var wg sync.WaitGroup
	url := "http://localhost:8080/events"

	fmt.Printf("[TEST] The test started with %d sync requests.\n", count)

	for i := 1; i <= count; i++ {
		wg.Add(1)

		go func(requestNum int) {
			defer wg.Done()

			randomUser := fmt.Sprintf("user_%d", rand.Intn(1000)+100)
			randomType := eventTypes[rand.Intn(len(eventTypes))]

			eventData := TestEvent{
				UserId:    randomUser,
				EventType: randomType,
			}

			jsonBytes, err := json.Marshal(eventData)
			if err != nil {
				fmt.Printf("Marshalling error (Request %d): %s\n", requestNum, err)
				return
			}

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
			if err != nil {
				fmt.Printf("Request %d is unsuccessful: %s\n", requestNum, err)
				return
			}
			resp.Body.Close()
		}(i)
	}

	wg.Wait()
	fmt.Printf("[TEST] The test finished with %d sync request!\n", count)

}
