package server

import (
  "fmt";
  "time";
  "net/http";
  "path/filepath";
  "encoding/json";
)

const (
  TEXT string = "TEXT_RESPONSE"
  FILE string = "FILE_RESPONSE"
  REDIRECT string = "REDIRECT_RESPONSE"
)

func newResponse(request *Request) *Response {
  return &Response{
    request: request,
    Code: 200,
  }
}

type Response struct {
  request *Request
  responseType string
  Code int
  Body string
  Duration time.Duration
}

func (r *Response) Plain(data string, params ...interface{}) *Response {
  r.responseType = TEXT
  r.Body = fmt.Sprintf(data, params...)
  return r
}

func (r *Response) Json(data interface{}) *Response {
  r.SetHeader("Content-Type", "application/json")
  marshaled, err := json.Marshal(data)
  if err != nil {
    return r.Error(500, "Unable to encode response. %s", err.Error())
  }
  return r.Plain(string(marshaled))
}

func (r *Response) Error(code int, message string, params ...interface{}) *Response {
  r.SetCode(code)
  return r.Plain(message, params...)
}

func (r *Response) ErrorJson(code int, data interface{}) *Response {
  r.SetCode(code)
  return r.Json(data)
}

func (r *Response) File(path string) *Response {
  r.responseType = FILE
  r.Body = filepath.Join(r.request.server.config.StaticRoot, path)
  return r
}

func (r *Response) Redirect(code int, url string) *Response {
  r.responseType = REDIRECT
  r.Body = url
  return r.SetCode(code)
}

func (r *Response) SetCode(code int) *Response {
  r.Code = code
  return r
}

func (r *Response) SetHeader(key string, value string) *Response {
  r.request.responseWriter.Header().Set(key, value)
  return r
}

func (r *Response) write() {
  r.Duration = time.Now().Sub(r.request.createdAt)
  switch r.responseType {
    case TEXT:
      r.request.responseWriter.WriteHeader(r.Code)
      fmt.Fprint(r.request.responseWriter, r.Body)
    case FILE:
      http.ServeFile(r.request.responseWriter, r.request.rawRequest, r.Body)
    case REDIRECT:
      http.Redirect(r.request.responseWriter, r.request.rawRequest, r.Body, r.Code)
    default:
      fmt.Fprint(r.request.responseWriter, "No response type specified.")
  }
}