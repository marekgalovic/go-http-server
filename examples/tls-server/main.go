package main

import (
  "log";
  "github.com/marekgalovic/go-http-server"
)


func main() {
  config := server.NewConfig()
  config.CertFile = "./server.crt"
  config.KeyFile = "./server.key"
  s := server.NewServer(config)

  s.Get("/greeting", EncryptedGreeting, nil)

  err := s.Listen()
  if err != nil {
    log.Fatal(err)
  }
}

func EncryptedGreeting(request *server.Request) {
  request.Respond("Encrypted hi!")
}