package service

import (
  "io";
  "fmt";
  "net/url";
  "net/http";
  "encoding/json"
)

func newRequest(method string, path string, params url.Values, body io.ReadCloser, responseWriter http.ResponseWriter) *Request {
  return &Request{method, path, params, body, responseWriter, make(chan string)}
}

type Request struct {
  Method string
  Path string
  Params url.Values
  body io.ReadCloser
  responseWriter http.ResponseWriter
  response chan string
}

func (r *Request) Get(param string) string {
  return r.Params.Get(param)
}

func (r *Request) Json(value interface{}) error {
  decoder := json.NewDecoder(r.body)
  return decoder.Decode(&value)
}

func (r *Request) Respond(data string, params ...interface{}) {
  r.response <- fmt.Sprintf(data, params...)
}

func (r *Request) RespondJson(data interface{}) {
  r.SetHeader("Content-Type", "application/json")
  if marshaled, err := json.Marshal(data); err != nil {
    r.Error(500, "Unable to encode response.")
  } else {
    r.Respond(string(marshaled))
  }
}

func (r *Request) Error(code int, message string, params ...interface{}) {
  r.responseWriter.WriteHeader(code)
  r.Respond(message, params...)
}

func (r *Request) ErrorJson(code int, data interface{}) {
  r.responseWriter.WriteHeader(code)
  r.RespondJson(data)
}

func (r *Request) SetHeader(key string, value string) {
  r.responseWriter.Header().Set(key, value)
}

func (r *Request) writeResponse() {
  fmt.Fprintf(r.responseWriter, <- r.response)
}