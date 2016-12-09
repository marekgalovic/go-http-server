package server

import (
  "fmt";
  "time";
  "net/http";
)

func newResponse(writer http.ResponseWriter) *Response {
  return &Response{
    Code: 200,
    writer: writer,
    bodyChan: make(chan string),
    createdAt: time.Now().UTC(),
  }
}

type Response struct {
  Code int
  Body string
  Duration time.Duration
  bodyChan chan string
  createdAt time.Time
  writer http.ResponseWriter
}

func (r *Response) setBody(body string) {
  r.bodyChan <- body
}

func (r *Response) setCode(code int) {
  r.Code = code
  r.writer.WriteHeader(r.Code)
}

func (r *Response) setHeader(key string, value string) {
  r.writer.Header().Set(key, value)
}

func (r *Response) write() *Response {
  r.Body = <- r.bodyChan
  r.Duration = time.Now().Sub(r.createdAt)

  fmt.Fprint(r.writer, r.Body)
  return r
}