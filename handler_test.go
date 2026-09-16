package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
	 Test Function: createEventHandler(http.ResponseWriter, http.Request)
		Input: invalid json format ("{")
	 	Expected Output: 400, "text/plain; charset=utf-8", "Invalid JSON format!\n"
*/
func TestCreateEventHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader("{"),
	)
	w := httptest.NewRecorder()

	createEventHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}

	if got, want := resp.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("content type = %q, want %q", got, want)
	}

	if got, want := string(body), "Invalid JSON format!\n"; got != want {
		t.Errorf("response body = %q, want %q", got, want)
	}
}

/*
	 Test Function: createEventHandler(http.ResponseWriter, http.Request)
		Input: nil body
	 	Expected Output: 400, "text/plain; charset=utf-8", "Invalid JSON format!\n"
*/
func TestCreateEventHandlerEmptyBody(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		nil,
	)
	w := httptest.NewRecorder()

	createEventHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}

	if got, want := resp.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("content type = %q, want %q", got, want)
	}

	if got, want := string(body), "Invalid JSON format!\n"; got != want {
		t.Errorf("response body = %q, want %q", got, want)
	}
}

/*
	 Test Function: createEventHandler(http.ResponseWriter, http.Request)
		Input: invalid event type
	 	Expected Output: 400, "text/plain; charset=utf-8", "Invalid event type!\n"
*/
func TestCreateEventHandlerInvalidEventType(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader(
			`{"user_id":"test-user","event_type":"unknown"}`,
		),
	)
	w := httptest.NewRecorder()

	// create a separate buffered channel so this test does not use events from other tests
	oldChannel := eventChan
	eventChan = make(chan Event, 1)
	t.Cleanup(func() { eventChan = oldChannel }) // restore the original channel when the test finishes

	createEventHandler(w, req)

	// check that the invalid event was not sent to the channel
	select {
	case got := <-eventChan:
		t.Errorf("invalid event was queued: %+v", got)
	default:
		// expected: invalid event was not queued
	}

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}

	if got, want := resp.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("content type = %q, want %q", got, want)
	}

	if got, want := string(body), "Invalid event type!\n"; got != want {
		t.Errorf("response body = %q, want %q", got, want)
	}
}

/*
	 Test Function: createEventHandler(http.ResponseWriter, http.Request)
		Input: empty user id
	 	Expected Output: 400, "text/plain; charset=utf-8", "User ID is required!\n"
*/
func TestCreateEventHandlerEmptyUserID(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader(
			`{"user_id":"","event_type":"click"}`,
		),
	)
	w := httptest.NewRecorder()

	createEventHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}

	if got, want := resp.Header.Get("Content-Type"), "text/plain; charset=utf-8"; got != want {
		t.Errorf("content type = %q, want %q", got, want)
	}

	if got, want := string(body), "User ID is required!\n"; got != want {
		t.Errorf("response body = %q, want %q", got, want)
	}
}

/*
	 Test Function: createEventHandler(http.ResponseWriter, http.Request)
		Input: valid event
	 	Expected Output: 201, "The event has been processed"
*/
func TestCreateEventHandlerValidEvent(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader(
			`{"user_id":"test-user","event_type":"click"}`,
		),
	)
	w := httptest.NewRecorder()

	// create a separate buffered channel so this test does not use events from other tests
	oldChannel := eventChan
	eventChan = make(chan Event, 1)
	t.Cleanup(func() { eventChan = oldChannel }) // restore the original channel when the test finishes

	createEventHandler(w, req)

	// check if the handler sent the correct event to the channel
	select {
	case got := <-eventChan:
		if got.UserId != "test-user" || got.EventType != EventClick {
			t.Errorf("queued event = %+v, want user=test-user and type=click", got)
		}
	default:
		t.Fatal("handler did not queue an event")
	}

	resp := w.Result()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code = %d, want %d", got, want)
	}

	if got, want := string(body), "The event has been processed"; got != want {
		t.Errorf("response body = %q, want %q", got, want)
	}
}
