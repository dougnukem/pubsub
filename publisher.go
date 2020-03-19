package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Publisher is an http client to publish messages to a PubSub server
type Publisher struct {
	addr string
}

// NewPublisher creates a *Publisher with specified HTTP address
func NewPublisher(addr string) *Publisher {
	return &Publisher{addr: addr}
}

// Publish sends message to PubSub server
func (p *Publisher) Publish(msg string) error {
	resp, err := http.Post(fmt.Sprintf("http://%s/publish", p.addr), "text/html", bytes.NewBufferString(msg))
	if err != nil {
		return fmt.Errorf("Error publishing message: %w", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading publish message response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error with published message [%s] to [%s] HTTP response [status-code=%d] [status=%s]: %s", msg, p.addr, resp.StatusCode, resp.Status, b)
	}

	return nil
}
