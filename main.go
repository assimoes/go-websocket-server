package main

import (
  "flag"
  "log"
  "net/http"
  "encoding/json"
  "os"
)

type Configuration struct {
  Address string
}

func main() {

  file, _ := os.Open("config.json")
  decoder := json.NewDecoder(file)

  config := Configuration{}
  err := decoder.Decode(&config)

  if err != nil {
    panic("Unable to parse config file")
  }

  var addr = flag.String("addr", config.Address, "Muzzley Websocket Server Address")

  flag.Parse()
  go h.run()
  http.HandleFunc("/ws", serveWs)

  _err := http.ListenAndServe(*addr, nil)
  if _err != nil {
    log.Fatal("Error staring the server: ", err)
  }

}