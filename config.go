package server

func NewConfig() *Config {
  return &Config{MaxWorkers: 4, Address: "127.0.0.1", Port: 5000}
}

type Config struct {
  MaxWorkers int
  Address string
  Port int
  StaticRoot string
  CertFile string
  KeyFile string
}