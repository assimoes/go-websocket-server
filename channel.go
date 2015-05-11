package main

type channel struct {
  key string
  connections map[*connection]bool
}

