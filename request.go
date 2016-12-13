package server

import (
  "io";
  "fmt";
  "time";
  "errors";
  "net/url";
  "net/http";
  "encoding/json";
  "path/filepath";
)

func newRequest(request *http.Request, responseWriter http.ResponseWriter, server *Server) *Request {
  request.ParseForm()
  return &Request{
    Method: request.Method,
    Path: request.URL.Path,
    Host: request.Host,
    Params: request.Form,
    Header: request.Header,
    RemoteAddr: request.RemoteAddr,
    Body: request.Body,
    server: server,
    rawRequest: request,
    responseWriter: responseWriter,
    createdAt: time.Now().UTC(),
  }
}

type Request struct {
  Method string
  Path string
  Host string
  Params url.Values
  Header http.Header
  RemoteAddr string
  Body io.ReadCloser
  server *Server
  rawRequest *http.Request
  responseWriter http.ResponseWriter
  createdAt time.Time
}

func (r *Request) Response() *Response {
  return newResponse(r)
}

func (r *Request) Get(param string) string {
  return r.Params.Get(param)
}

func (r *Request) Empty(param string) bool {
  return r.Get(param) == ""
}

func (r *Request) Json(value interface{}) error {
  decoder := json.NewDecoder(r.Body)
  return decoder.Decode(&value)
}

func (r *Request) Route(path string) string {
  return fmt.Sprintf("http://%s", filepath.Join(r.Host, path))
}

func (r *Request) GetCookie(name string) string {
  cookie, err := r.rawRequest.Cookie(name)
  if err != nil {
    return ""
  }
  return cookie.Value
}

func (r *Request) SetCookie(name string, value string, duration time.Duration) {
  cookie := &http.Cookie{Name: name, Value: value, Expires: time.Now().Add(duration)}
  http.SetCookie(r.responseWriter, cookie)
}

func (r *Request) Session(name string) (*Session, error) {
  if r.server.sessions == nil {
    return nil, errors.New("Session store is not present.")
  }
  session, err := r.server.sessions.Get(r.rawRequest, name)
  if err != nil {
    return nil, err
  }
  return newSession(session, r.rawRequest, r.responseWriter), nil
}