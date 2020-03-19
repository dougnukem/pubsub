package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	// This test will create a number of PubSub subscribers and publishers that will run concurrently
	// the test asserts that the exact number of messages are received as expected

	numSubscribers := 20
	numPublishers := 20
	numMessages := 50
	expectedTotalMsgs := numMessages * numPublishers * numSubscribers
	addr := "127.0.0.1:9911"

	l := log.New(os.Stdout, "[PUBSUB SERVER]: ", log.Ldate|log.Ltime)
	s := NewServer(l, addr)
	go func() {
		err := s.ListenAndServe()
		l.Printf("server listen error: %s", err)
	}()
	// HACK to wait for PubSub HTTP server to be up
	time.Sleep(1 * time.Second)

	// used to track expected messages received
	var wg sync.WaitGroup
	var msgsReceived uint64
	// expected messages to receive per PubSub client
	wg.Add(expectedTotalMsgs)

	// connect Subscribers
	for i := 0; i < numSubscribers; i++ {
		cl := NewClient(addr)
		wsCl, msgCh, errCh, err := cl.Subscribe()
		if err != nil {
			t.Errorf("Failed to setup PubSub client[%d]: %s", i, err)
			return
		}
		l := log.New(os.Stdout, fmt.Sprintf("[PUBSUB CLIENT: %s]: ", wsCl.LocalAddr()), log.Ldate|log.Ltime)
		go func() {
			for {
				select {
				case msg, ok := <-msgCh:
					if !ok {
						l.Printf("client subscription shutdown")
						return
					}
					// received message
					atomic.AddUint64(&msgsReceived, 1)
					l.Printf("client received message [totalmsgs=%d]: %s", atomic.LoadUint64(&msgsReceived), msg)
					wg.Done()

				case err, ok := <-errCh:
					if !ok {
						l.Printf("client subscription shutdown")
						return
					}
					l.Printf("client error: %v", err)
				}
			}
		}()
	}

	// create Publishers and send messages
	for i := 0; i < numPublishers; i++ {
		pubID := i
		go func() {
			l := log.New(os.Stdout, fmt.Sprintf("[PUBSUB PUBLISHER: %d]: ", pubID), log.Ldate|log.Ltime)
			p := NewPublisher(addr)
			for j := 0; j < numMessages; j++ {
				err := p.Publish(fmt.Sprintf("[Publisher: %d] Message [%d]", pubID, j))
				if err != nil {
					l.Printf("Publisher failed to send message: %s", err)
					return
				}
			}
		}()
	}

	// wait to receive all expected messages or timeout (test failure)
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		return
	case <-time.After(30 * time.Second):
		t.Errorf("Expected %d messages to be received but received %d before timing out", expectedTotalMsgs, atomic.LoadUint64(&msgsReceived))
		return
	}
}
