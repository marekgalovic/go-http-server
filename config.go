package service

import (
  "time"
)

func NewConfig() *Config {
  return &Config{}
}

type Config struct {
  Address string
  Port int
  ReadTimeout *time.Duration
  WriteTimeout *time.Duration
  CertFile string
  KeyFile string
}