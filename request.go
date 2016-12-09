package server

import (
  "io";
  "fmt";
  "net/url";
  "net/http";
  "encoding/json"
)

func newRequest(request *http.Request, responseWriter http.ResponseWriter) *Request {
  return &Request{
    Method: request.Method,
    Path: request.URL.Path,
    Params: request.Form,
    Header: request.Header,
    RemoteAddr: request.RemoteAddr,
    Body: request.Body,
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