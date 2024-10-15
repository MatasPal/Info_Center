# Info Center

Info Center is a backend service developed in Go that enables real-time communication between clients through message sending to various topics. This project implements Server-Sent Events (SSE) for receiving messages and utilizes concurrency to handle multiple clients simultaneously.

## Features

- **Dynamic Topic Management**: Clients can send and receive messages to and from dynamically created topics.
- **Auto-Incrementing Message IDs**: Each message sent to a topic is assigned a unique, auto-incrementing global ID.
- **Server-Sent Events**: Clients can subscribe to topics and receive messages in real-time.
- **CORS Support**: The service allows cross-origin requests to facilitate frontend interactions.
- **Timeout Handling**: Clients are disconnected after a specified time (30 seconds) of inactivity, with timeout events sent before disconnection.

## Getting Started

### Prerequisites

- Go (version 1.21.3 or higher)

### Running the Application

1. Run the main:
   ```bash
   go run main.go
   
2. The server will start on **localhost:8080**

### API Endpoints

#### POST /infocenter/{topic}
Send a message to a specific topic.

- **Request**:
  ```bash
  POST /infocenter/{topic}
  Content-Type: application/x-www-form-urlencoded
  message=Hello

- **Example using curl**:
```bash
curl -X POST http://localhost:8080/infocenter/topic1 --data "message=Hello"
```


#### GET /infocenter/{topic}
Subscribe to a specific topic and receive messages in real-time.

- **Request**:
  ```bash
  GET /infocenter/{topic}

- **Example using curl**:
```bash
curl -N http://localhost:8080/infocenter/topic1
```

### Using the Frontend
The project includes a simple frontend located at InfoCenter/front/index.html, which can be used to interact with the backend and test the real-time messaging functionality.

To use the frontend:

1. **Open the index.html file in your web browser**:

   - Navigate to InfoCenter/front/index.html in your project directory.
   - Double-click the file to open it in your default browser.

2. **The frontend allows you to**:

   - Subscribe to topics.
   - Send messages to specific topics.

Ensure the backend server is running on localhost:8080 for the frontend to communicate properly with the backend.
