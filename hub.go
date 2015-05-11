package main

import (
  "log"
)

type hub struct {
  connections map[*connection]bool
  broadcast chan []byte
  register chan *connection
  unregister chan *connection
  broadcastToChannel chan []byte
  channels map[string]channel
}

var h = hub {
  broadcast: make(chan []byte),
  register: make(chan *connection),
  unregister: make(chan *connection),
  connections: make(map[*connection]bool),
  broadcastToChannel: make(chan []byte),
  channels: make(map[string]channel),
}

func (h *hub) run() {

  chann := channel{
    key: "main lobby",
    connections: make(map[*connection]bool),
  }

  h.channels["lobby"] = chann

  for {
    select {
    case c := <-h.register:
      h.connections[c] = true
      log.Println("Joined: ", c)
      if num :=len(h.channels["lobby"].connections); num < 10 {
        h.channels["lobby"].connections[c]=true
      }
    case c:= <-h.unregister:
      if _,ok := h.connections[c]; ok {
        log.Println("Left:", c)
        delete(h.connections,c)
        close(c.send)
      }
      if _,ok := h.channels["lobby"].connections[c]; ok {
        delete(h.channels["lobby"].connections,c)
      }
    case m := <-h.broadcast:
      log.Println(string(m));
      for c := range h.connections {
        select {
        case c.send <-m :
        default:
          close(c.send)
          delete(h.connections, c)
        }
      }
    case m := <-h.broadcastToChannel:
      log.Println("sending message to lobby users ", string(m))
      for c := range h.channels["lobby"].connections {
        select {
        case c.send <-m :
        default:
          close(c.send)
          delete(h.connections,c)
        }
      }
    }
  }
}