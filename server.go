package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

// Server is an HTTP server that manages an in-memory collection of WebSocket clients that subscribe for messages via a /subscribe endpoint.
// The server also allows messages to be published to the subscribed clients via a /publish endpoint
type Server struct {
	subscribers map[*Subscriber]struct{}
	subscribe   chan *Subscriber
	unsubscribe chan *Subscriber
	publish     chan []byte
	httpAddr    string
	httpMux     *http.ServeMux
	logger      *log.Logger
}

// NewServer creates a new *Server that can handle HTTP WebSocket subscribers and publish messages to subscribed clients
func NewServer(l *log.Logger, httpAddr string) *Server {
	s := &Server{
		subscribers: map[*Subscriber]struct{}{},
		subscribe:   make(chan *Subscriber),
		unsubscribe: make(chan *Subscriber),
		publish:     make(chan []byte),
		httpAddr:    httpAddr,
		httpMux:     http.NewServeMux(),
		logger:      l,
	}

	// HTTP PubSub routes for subscribing and publishing
	s.httpMux.HandleFunc("/subscribe", s.handleSubscribe)
	s.httpMux.HandleFunc("/publish", s.handlePublish)

	return s
}

// ListenAndServe starts HTTP server and go routines for managing PubSub Websocket subscribers and Publishing
func (s *Server) ListenAndServe() error {
	// handle pub sub channel management in a go routine
	go s.handlePubSub()

	s.logger.Printf("listening on: %s", s.httpAddr)

	return http.ListenAndServe(s.httpAddr, s.httpMux)
}

// handleSubscribe handles a WebSocket client subscribing for messages that will be published and broadcast to all subscribed clients via this PubSub server
func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Printf("Error upgrading: %v", err)
		return
	}

	defer c.Close()

	subscriber := NewSubscriber(c, s.logger)

	// subscribe this client to the collection of subscribers that will receive published messages
	s.subscribe <- subscriber

	// unsusbcribe when WebSocket is done receiving messages (e.g. closed by client)
	defer func() {
		s.unsubscribe <- subscriber
	}()

	// block and receive messages
	subscriber.receiveMessages()
}

// publish handles publishing a message to all subscribed WebSocket clients connnected to this PubSub server
func (s *Server) handlePublish(w http.ResponseWriter, r *http.Request) {
	// only handle HTTP POST method
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.Printf("publish: Error reading HTTP body: %s\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Broadcast message to all subscribed clients
	s.publish <- message
}

func (s *Server) handlePubSub() {
	for {
		select {
		case subscriber := <-s.subscribe:
			// add new subscriber to in-memory set that will receive published messages
			s.subscribers[subscriber] = struct{}{}
			s.logger.Printf("new subscriber [count=%d]: %s\n", len(s.subscribers), subscriber.conn.RemoteAddr())
		case unsubscriber := <-s.unsubscribe:
			s.logger.Printf("unsubscribe: %s\n", unsubscriber.conn.RemoteAddr())
			delete(s.subscribers, unsubscriber)
		case msg := <-s.publish:
			s.logger.Printf("publish to [count=%d] subscribers: %s\n", len(s.subscribers), msg)
			for subscriber := range s.subscribers {
				subscriber.SendMessage(msg)
			}
		}
	}
}
