# Event Logger Service

A simple REST API built with Go. This is a beginner-level project I created over three days to learn the fundamentals of Go backend development. It receives user events, processes them in the background, and provides basic statistics.

## Concepts Learned

Through building this project, I explored and practiced several core Go concepts:

* **HTTP Handlers:** Creating routes and handling basic GET and POST requests.
* **In-Memory Data Structures:** Using slices to store data and maps to count and group events.
* **Custom Types & Validation:** Creating an enum-like structure to accept only specific event types (login, signup, click).
* **Basic Concurrency:** Using Goroutines and Channels to take incoming requests and process them in the background without making the user wait.

## API Endpoints

### 1. Create an Event
* **Method:** POST
* **Path:** `/events`
* **Body:** 
  ```json
  {
    "user_id": "ZXZlbnQtbG9nZ2VyLXNlcnZpY2U=",
    "event_type": "login"
  }
  ```
* **Note:** `event_type` must be one of: login, signup, click. It returns a 201 Created (or 202 Accepted) status.

### 2. List Events
* **Method:** GET
* **Path:** `/events`
* **Optional Filter:** `/events?user_id=ZXZlbnQtbG9nZ2VyLXNlcnZpY2U=`
* **Returns:** A JSON list of recorded events.

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

## How to Run

1. Clone this repository to your desktop:
   ```bash
   git clone https://github.com/mertkiyar/event-logger-service.git
   cd event-logger-service
   ```

2. Start the server:
   ```bash
   go run .
   ```

3. The service will start and listen on port 8080. You can test it using curl or Postman.