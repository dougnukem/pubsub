[![codecov](https://codecov.io/gh/dougnukem/pubsub/branch/master/graph/badge.svg)](https://codecov.io/gh/dougnukem/pubsub)

# pubsub server
Implements a PubSub HTTP server that manages a collection of in-memory WebSocket clients that subscribe, and a publish endpoint that will deliver a message to all subscribed HTTP WebSocket clients.

## Automated test

`TestPubSub` is an integration test that will start a PubSub server and spawn concurrent subscribers and publishers to ensure the expected number of messages are received.

```
$ go test -race -v -count 1 .
...
[PUBSUB SERVER]: 2020/03/19 01:08:21 publish to [count=50] subscribers: [Publisher: 24] Message [4]
...
[PUBSUB CLIENT: 127.0.0.1:49333]: 2020/03/19 01:08:21 client received message: [Publisher: 48] Message [4]
--- PASS: TestPubSub (1.41s)
PASS
ok  	github.com/dougnukem/pubsub	1.554s
```

# Commands

Commands can be run directly from checked out repo, or via:
```
$ go get github.com/dougnukem/pubsub

# Assuming $GOPATH/bin is in your $PATH
$ pubsub server -help
Usage of server:
  -addr string
    	pubsub server address to listen on (default ":8080")
```

## `server` - starts an HTTP server that listens for WebSocket clients to subscribe and messages to be published
```
$ go run -race main.go server -addr=:8080
```

Endpoints:
- `/subscribe` - WebSocket endpoint for a client to subscribe and be delivered published messages
- `/publish` - HTTP POST endpoint for messages to be published and sent to all subscribed WebSocket clients


## `client` - subscribes to websocket server and is deliverd published messages
```
$ go run -race main.go client -addr=:8080
```
## `publish` messages

```
Usage:
$ go run -race main.go publish -addr=:8080 -message="Message to publish to clients
```
