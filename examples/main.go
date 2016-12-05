package main

import (
  service "github.com/marekgalovic/go-http-service"
)

func main() {
  config := service.NewConfig()
  config.Address = "0.0.0.0"
  config.Port = 80

  s := service.NewServer(config)
  s.Get('/ping', PingEndpoint)
  s.Listen()
}

func PingEndpoint(request *service.Request) {
  request.Respond("ok")
}