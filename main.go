package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	// DefaultAddr is the default pubsub server HTTP listen address
	DefaultAddr = ":8080"
)

func main() {
	subCmdUsage := fmt.Sprintf("expected subcommand: server | client | publish\n")
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, subCmdUsage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		serverCmd(os.Args[2:]...)
	case "client":
		clientCmd(os.Args[2:]...)
	case "publish":
		publishCmd(os.Args[2:]...)
	default:
		fmt.Fprintf(os.Stderr, subCmdUsage)
		os.Exit(1)
	}
}

func serverCmd(args ...string) {
	var addr string

	cmd := flag.NewFlagSet("server", flag.ExitOnError)
	cmd.StringVar(&addr, "addr", DefaultAddr, "pubsub server address to listen on")
	mustParseCmd(cmd, args)

	l := log.New(os.Stdout, "[PUBSUB SERVER]: ", log.Ldate|log.Ltime)
	s := NewServer(l, addr)

	if err := s.ListenAndServe(); err != nil {
		l.Printf("HTTP server shutting down: %s\n", err)
	}
}

func clientCmd(args ...string) {
	var addr string

	cmd := flag.NewFlagSet("client", flag.ExitOnError)
	cmd.StringVar(&addr, "addr", DefaultAddr, "pubsub server address to subscribe to")
	l := log.New(os.Stdout, "[PUBSUB CLIENT]: ", log.Ldate|log.Ltime)
	mustParseCmd(cmd, args)

	cl := NewClient(addr)
	wsCl, messages, errCh, err := cl.Subscribe()
	l.SetPrefix(fmt.Sprintf("[PUBSUB CLIENT: %s]: ", wsCl.LocalAddr()))
	l.Printf("subscribed waiting for published messages")
	if err != nil {
		l.Fatalf("Error subscribing: %v", err)
	}

	// block and read published messages from PubSub client
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				l.Printf("client subscription shutdown")
				return
			}
			l.Printf("client received message: %s", msg)
		case err, ok := <-errCh:
			if !ok {
				l.Printf("client subscription shutdown")
				return
			}
			l.Fatalf("client error: %v", err)
		}
	}
}

func publishCmd(args ...string) {
	var (
		addr    string
		message string
	)

	cmd := flag.NewFlagSet("publish", flag.ExitOnError)
	cmd.StringVar(&addr, "addr", DefaultAddr, "pubsub server address to publish to")
	cmd.StringVar(&message, "message", "", "message to publish to subscribed clients")
	mustParseCmd(cmd, args)

	if message == "" {
		printfUsageExit(cmd, "error: requires -message\n")
	}

	l := log.New(os.Stdout, "[PUBSUB PUBLISHER]: ", log.Ldate|log.Ltime)

	p := NewPublisher(addr)
	err := p.Publish(message)
	if err != nil {
		l.Fatalf("Error publishing message: %s", err)
	}

	fmt.Printf("pubsub %s: addr[%s] message[%s]\n", cmd.Name(), addr, message)
}

func printfUsageExit(cmd *flag.FlagSet, fmtMsg string, argv ...interface{}) {
	fmt.Fprintf(os.Stderr, fmtMsg, argv...)
	fmt.Printf("usage: pubsub %s\n", cmd.Name())
	cmd.PrintDefaults()
	os.Exit(1)
}

func mustParseCmd(cmd *flag.FlagSet, args []string) {
	err := cmd.Parse(args)
	if err != nil {
		printfUsageExit(cmd, "error parsing command: %s\n", err)
	}
}
