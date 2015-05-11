# go-websocket-server

## Introduction

Shameless clone of Gorilla's websocket sample. Using this as a base to build a websocket server with rooms support.
~~For now, these are hardcoded.~~

## External Dependencies
```
gorilla/websocket
```
## Setup

You must have Go installed on your machine.

#### Clone the repo

```
$ git clone git@github.com:assimoes/go-websocket-server.git && cd go-websocket-server
```
#### Get the dependencies

```
$ go get 
```

#### Change websocket server address 

Modify ```config.json``` file accordingly or just leave it be.

#### Build the server

```
$ go build
```

#### Run it

```
$ ./go-websocket-server
```

#### Message structure

If you want to send a message to a specific channel, your payload must be a json structure:

```
{
  "channel": "the channel",
  "data": "the data"
}

```
If you don't specify a Channel, then the message is broadcasted to all connections on the hub

To Join a channel you should send the following payload:

```
{
  "channel": "the channel",
  "command": "join"
}
