package main

import (
  "log"
  "encoding/json"
)

type hub struct {
  connections map[*connection]bool
  broadcast chan []byte
  register chan *connection
  unregister chan *connection
  broadcastToChannel chan []byte
  channels map[string]channel
  joinChannel chan *connection
}

var h = hub {
  broadcast: make(chan []byte),
  register: make(chan *connection),
  unregister: make(chan *connection),
  connections: make(map[*connection]bool),
  broadcastToChannel: make(chan []byte),
  channels: make(map[string]channel),
  joinChannel: make(chan *connection),
}

func (h *hub) run() {

  chann := channel{
    key: "main lobby",
    connections: make(map[*connection]bool),
  }

  h.channels["lobby"] = chann

  // NEED TO IMPLEMENT CASE FOR JOIN CHANNEL MESSAGES.
  // WHEN SUCH MESSAGE ARRIVES ON THAT CHANNEL, IT SHOULD ADD THE CONNECTION TO THE NEW/EXISTING CHANNEL

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
    case c := <-h.joinChannel:
      if _, ok := h.channels[c.channel]; !ok {
          chann := channel{
            key: c.channel,
            connections: make(map[*connection]bool),
          }
          log.Println("creating channel")
          h.channels[c.channel] = chann
          h.channels[c.channel].connections[c]=true
      } else {
        h.channels[c.channel].connections[c]=true
      }


    case m := <-h.broadcastToChannel:
      var msg IncMessage
      _err := json.Unmarshal(m, &msg)

      if _err != nil {
        log.Println(_err)
      }

      for c := range h.channels[msg.Channel].connections {
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