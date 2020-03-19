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
$ pubsub server
[PUBSUB SERVER]: 2020/03/19 01:36:03 listening on: :8080
[PUBSUB SERVER]: 2020/03/19 01:36:10 new subscriber [count=1]: 127.0.0.1:55127
[PUBSUB SERVER]: 2020/03/19 01:36:18 publish to [count=1] subscribers: Hello there
[PUBSUB SERVER]: 2020/03/19 01:37:55 new subscriber [count=2]: 127.0.0.1:55218
```

Endpoints:
- `/subscribe` - WebSocket endpoint for a client to subscribe and be delivered published messages
- `/publish` - HTTP POST endpoint for messages to be published and sent to all subscribed WebSocket clients


## `client` - subscribes to websocket server and is deliverd published messages
```
$ pubsub client
[PUBSUB CLIENT: 127.0.0.1:55127]: 2020/03/19 01:36:10 subscribed waiting for published messages
[PUBSUB CLIENT: 127.0.0.1:55127]: 2020/03/19 01:36:18 client received message: Hello there
```
## `publish` messages

```
Usage:
$ pubsub publish -message="Hello there"
[PUBSUB PUBLISHER]: 2020/03/19 01:39:14 Published message[Hello there] to [:8080]
```
