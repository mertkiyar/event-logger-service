package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
