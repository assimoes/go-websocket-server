package main

import(
  "github.com/gorilla/websocket"
  "log"
  "net/http"
  "time"
)

const (
  writeWait = 10 * time.Second
  pongWait = 60 * time.Second
  pingPeriod = (pongWait * 9) / 10
  maxMessageSize = 512
)


var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool { return true},
}

type connection struct {
  ws *websocket.Conn
  send chan []byte
}

type incmessage struct {
  channel string
}

func (c *connection) readPump() {
  defer func() {
    h.unregister <- c
    c.ws.Close()
  }()

  c.ws.SetReadLimit(maxMessageSize)
  c.ws.SetReadDeadline(time.Now().Add(pongWait))
  c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil})

  for {
    _, message, err := c.ws.ReadMessage()
    if err != nil {
      break
    }

    msg := incmessage{channel: "mainlobby"}
log.Println("incoming message step 1")
    if msg.channel == "mainlobby" {
      if _,allowed := h.channels["lobby"].connections[c]; allowed {
        log.Println("broadcast to channel")
        h.broadcastToChannel <- message
      }
    } else {
      log.Println("broadcast to all")
      h.broadcast <- message
    }
  }
}

func (c *connection) write(mt int, payload []byte) error {
  c.ws.SetWriteDeadline(time.Now().Add(writeWait))
  return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump() {
  ticker := time.NewTicker(pingPeriod)
  defer func() {
    ticker.Stop()
    c.ws.Close()
  }()

  for {
    select {
    case message, ok := <-c.send:
      if !ok {
        c.write(websocket.CloseMessage, []byte{})
        return
      }
      if err := c.write(websocket.TextMessage, message); err != nil {
        return
      }
    case <- ticker.C:
      if err := c.write(websocket.PingMessage, []byte{}); err != nil {
        return
      }
    }
  }
}

func serveWs(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return
  }

  ws, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println(err)
    return
  }

  c := &connection{send: make(chan []byte), ws: ws}
  h.register <- c

  go c.writePump()
  c.readPump()
}