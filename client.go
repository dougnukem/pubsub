package main

import (
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

// Client is a PubSub HTTP websocket client that connects to a PubSub server
type Client struct {
	addr string
}

// NewClient creates a *Client used to subscribe to a HTTP PubSub Websocket server
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// Subscribe connects the client to a PubSub server and subscribes for published messages which can be read off the chan []byte
// errors on the Websocket can be read on the chan error
func (c *Client) Subscribe() (*websocket.Conn, chan []byte, chan error, error) {
	u := url.URL{Scheme: "ws", Host: c.addr, Path: "/subscribe"}

	wsCl, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, nil, nil, err
	}

	messages := make(chan []byte)
	errCh := make(chan error)

	go func() {
		defer func() {
			wsCl.Close()
			close(errCh)
			close(messages)
		}()

		for {
			_, msg, err := wsCl.ReadMessage()
			if err != nil {
				errCh <- fmt.Errorf("error websocket read: %w", err)
				return
			}
			messages <- msg
		}
	}()

	return wsCl, messages, errCh, nil
}
