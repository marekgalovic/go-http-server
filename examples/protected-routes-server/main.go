package main

import (
  "log";
  "github.com/marekgalovic/go-http-server"
)

type MyAuthProvider struct{}

func (p *MyAuthProvider) Verify(request *server.Request) *server.Response {
  if request.Empty("user_token") {
    return request.Response().Error(400, "Please provide user_token.")
  }
  if request.Get("user_token") != "secrettoken" {
    return request.Response().Error(401, "Invalid user_token.")
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

func SayHiOnlyToMarek(request *server.Request) *server.Response {
  return request.Response().Plain("Hi Marek!")
}