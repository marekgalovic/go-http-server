package server

import (
  "io";
  "fmt";
  "time";
  "errors";
  "net/url";
  "net/http";
  "encoding/json"
)

var (
  SessionStoreNotPresent error = errors.New("Session store is not present. Please initialize your server with session store.")
)

func newRequest(request *http.Request, responseWriter http.ResponseWriter) *Request {
  return &Request{
    Method: request.Method,
    Path: request.URL.Path,
    Params: request.Form,
    Header: request.Header,
    RemoteAddr: request.RemoteAddr,
    Body: request.Body,
    rawRequest: request,
    response: newResponse(responseWriter),
  }
}

type Request struct {
  Method string
  Path string
  Params url.Values
  Header http.Header
  RemoteAddr string
  Body io.ReadCloser
  rawRequest *http.Request
  response *Response
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

func (r *Request) Respond(data string, params ...interface{}) {
  r.response.setBody(fmt.Sprintf(data, params...))
}

func (r *Request) RespondJson(data interface{}) {
  r.response.setHeader("Content-Type", "application/json")
  marshaled, err := json.Marshal(data)
  if err != nil {
    r.Error(500, "Unable to encode response.")
    return
  }
  r.Respond(string(marshaled))
}

func (r *Request) Error(code int, message string, params ...interface{}) {
  r.response.setCode(code)
  r.Respond(message, params...)
}

func (r *Request) ErrorJson(code int, data interface{}) {
  r.response.setCode(code)
  r.RespondJson(data)
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
  http.SetCookie(r.response.writer, cookie)
}