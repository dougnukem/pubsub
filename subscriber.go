package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Subscriber is a HTTP WebSocket client that is subscribed for published messages on the PubSub server
type Subscriber struct {
	logger  *log.Logger
	conn    *websocket.Conn
	message chan []byte
}

// NewSubscriber creates a new *Subscriber to manage receiving messages to be broadcast to Subscriber
func NewSubscriber(conn *websocket.Conn, logger *log.Logger) *Subscriber {
	// Create a buffered channel to allow clients to buffer as they are broadcasting and not block other subscribers
	return &Subscriber{conn: conn, message: make(chan []byte, 10), logger: logger}
}

// SendMessage sends a message to the subscriber
func (s *Subscriber) SendMessage(message []byte) {
	s.message <- message
}

// receiveMessages will block and read messages sent to the s.message channel until closed or WebSocket client is disconnected
func (s *Subscriber) receiveMessages() {
	for {
		select {
		case msg := <-s.message:
			err := s.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				s.logger.Printf("Error writing client websocket message: %s\n", err)
				return
			}
		}
	}
}
