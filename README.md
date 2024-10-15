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


