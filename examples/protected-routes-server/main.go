package main

import (
  "log";
  "errors";
  "github.com/marekgalovic/go-http-server"
)

type MyAuthProvider struct{}

func (p *MyAuthProvider) Verify(request *server.Request) error {
  if request.Empty("user_token") {
    return errors.New("Please provide user_token param.")
  }
  if request.Get("user_token") != "secrettoken" {
    return errors.New("Invalid user_token param.")
  }
  return nil
}

func main() {
  config := server.NewConfig()
  auth := &MyAuthProvider{}
  s := server.NewServer(config)

  s.Get("/greeting", SayHiOnlyToMarek, auth)

  err := s.Listen()
  if err != nil {
    log.Fatal(err)
  }
}

func SayHiOnlyToMarek(request *server.Request) {
  request.Respond("Hi Marek!")
}