# Event Logger Service

A simple REST API built with Go to learn backend development. It receives user events, processes them in the background, saves them to MongoDB, and provides basic statistics.

## Concepts Learned

Through building this project, I explored and practiced several core Go concepts:

* **HTTP Handlers:** Creating routes and handling basic GET and POST requests.
* **In-Memory Data Structures:** Using slices to store data and maps to count and group events.
* **Custom Types & Validation:** Accepting specific event types (login, logout, signup, click) and checking required user IDs.
* **Basic Concurrency:** Using Goroutines and Channels to take incoming requests and process them in the background without making the user wait.
* **MongoDB:** Saving events and loading stored records when the service starts.
* **Middleware:** Adding a basic per-IP rate limiter to `/events`.
* **Testing:** Using HTTP test helpers and table-driven tests.

## API Endpoints

### 1. Create an Event
* **Method:** POST
* **Path:** `/events`
* **Body:** 
  ```json
  {
    "user_id": "mert",
    "product_id": "SM-798",
    "event_type": "click"
  }
  ```
* **Validation:** `user_id` is required and cannot be blank. `event_type` must be one of: login, logout, signup, click, using lowercase letters. `product_id` is optional; include it to count clicks for a specific product. Invalid JSON, an invalid event type, or a missing user ID returns `400 Bad Request`.
* **Response:** `202 Accepted` with the message `The event has been accepted`. This means the event was queued, not that MongoDB has already saved it.

Example request:

```bash
curl -i -X POST http://localhost:8080/events \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"mert","product_id":"SM-798","event_type":"click"}'
```

### 2. List Events
* **Method:** GET
* **Path:** `/events`
* **Optional Filter:** `/events?user_id=mert`
* **Returns:** A JSON list of events currently loaded in memory. No matching events returns an empty list.

### 3. Event Statistics
* **Method:** GET
* **Path:** `/stats`
* **Returns:** A JSON object showing the total count for each event type. Example:
  ```json
  {
    "login": 3,
    "signup": 1
  }
  ```

### 4. Product Click Statistics

* **Method:** GET
* **Path:** `/product-stats?user_id=mert&product_id=SM-798`
* **Returns:** The click count for that user and product. Other event types are not counted.

```json
{
  "user_id": "mert",
  "product_id": "SM-798",
  "click_count": 2
}
```

Both query parameters are required; a missing parameter returns `400 Bad Request`. No matching clicks returns `click_count: 0`.

```bash
curl 'http://localhost:8080/product-stats?user_id=mert&product_id=SM-798'
```

## How to Run

Requirements: Go 1.26.1 or later and a running MongoDB server. The MongoDB address is currently set to `mongodb://localhost:27017` in `main.go`. Events are stored in the `events` collection of the `event_logger` database.

1. Clone this repository to your desktop:
   ```bash
   git clone https://github.com/mertkiyar/event-logger-service.git
   cd event-logger-service
   ```

2. Start MongoDB. If you installed MongoDB 8.0 with Homebrew on macOS:

   ```bash
   brew services start mongodb-community@8.0
   ```

3. Start the server:
   ```bash
   go run .
   ```

4. The service listens on `http://localhost:8080`. MongoDB must be reachable before the HTTP server can start. You can test the API using curl or Postman.

## How Events Are Processed

1. The POST handler validates the request and puts the event into a channel.
2. A background worker reads the channel and inserts the event into MongoDB.
3. After a successful insert, the worker adds the event to the in-memory slice used by the GET endpoints.
4. When the service restarts, it loads stored MongoDB events into that slice before starting the HTTP server.

To check persistence, send one event and wait for the worker's `new event processed` message. Stop the API with Ctrl+C, run `go run .` again, and repeat the GET request without sending another POST. The saved event should still appear.

## Tests

```bash
go test ./...
```

Tests cover event type validation, event handlers, filtering, statistics, and the product click counter. They do not require MongoDB; database persistence is checked separately using the restart steps above.

## Current Scope

This is a small learning project, not a production-ready service.

* GET endpoints read from memory, not directly from MongoDB. Changes made outside the service are loaded on the next restart.
* Background processing means a newly accepted event may not appear immediately. Queued events can be lost if the process stops before they are saved; failed inserts are logged but not retried.
* The shared event slice is not yet protected against concurrent reads and writes, so the project is not intended for heavy concurrent use.
* The basic rate limiter applies only to `/events`. It currently allows 99 requests per IP and rejects the 100th with `429`. The counter resets after more than 10 seconds without a request; this is not a fixed 10-second window.
