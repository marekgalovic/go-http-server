package server

type AuthProvider interface {
  Verify(*Request) error
}