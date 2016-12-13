package server

import (
  "net/http";
  "github.com/gorilla/sessions";
)

func newSession(session *sessions.Session, request *http.Request, writer http.ResponseWriter) *Session {
  return &Session{session, request, writer}
}

type Session struct {
  session *sessions.Session
  request *http.Request
  writer http.ResponseWriter
}

func (s *Session) Get(key interface{}) interface{} {
  return s.session.Values[key]
}

func (s *Session) Set(key interface{}, value interface{}) error {
  s.session.Values[key] = value
  return s.save()
}

func (s *Session) Delete(key interface{}) error {
  delete(s.session.Values, key)
  return s.save()
}

func (s *Session) save() error {
  return s.session.Save(s.request, s.writer)
}