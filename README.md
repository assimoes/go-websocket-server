# go-websocket-server

## Introduction

Shameless clone of Gorilla's websocket sample. Using this as a base to build a websocket server with support to rooms.

For now, these are hardcoded.


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

Modify and copy ```config.json``` file to /etc/go-websocket-server (This is temporary, shouldn't be hardcoded)

#### Build the server

```
$ go build
```

#### Run it

```
$ ./go-websocket-server
```
