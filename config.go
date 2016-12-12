package server

func NewConfig() *Config {
  return &Config{Address: "127.0.0.1", Port: 5000}
}

type Config struct {
  Address string
  Port int
  StaticRoot string
  CertFile string
  KeyFile string
}