package main

import (
  "log";
  "github.com/marekgalovic/go-http-server"
)

func main() {
  config := server.NewConfig()
  s := server.NewServer(config)

  s.Get("/greeting", SayHi, nil)

  err := s.Listen()
  if err != nil {
    log.Fatal(err)
  }
}

func SayHi(request *server.Request) *server.Response {
  return request.Response().Plain("Hi!")
}