package server

import (
  "testing";

  "github.com/stretchr/testify/assert"
)

func TestNewConfigReturnsConfigObjectWithCorrectDefaultParameters(t *testing.T) {
  config := NewConfig()

  assert.Equal(t, "127.0.0.1", config.Address)
  assert.Equal(t, 5000, config.Port)
}