package service

import (
  "testing";
  "net/http";
  "net/http/httptest";

  "github.com/stretchr/testify/suite";
  "github.com/stretchr/testify/assert";
)

type ServerTestSuite struct{
  suite.Suite
  server *Server
}

func (suite *ServerTestSuite) SetupSuite() {
  suite.server = NewServer()
}

func (suite *ServerTestSuite) DummyEndpoint(request *Request) {
  request.Respond("ok%s%s", request.Get(":greeting"), request.Get(":id"))
}

func (suite *ServerTestSuite) TestHandlerExecutesCorrectAction() {
  suite.server.Get("/dummy", suite.DummyEndpoint)

  req, _ := http.NewRequest("GET", "/dummy", nil)
  res := httptest.NewRecorder()
  suite.server.Call(res, req)

  assert.Equal(suite.T(), res.Body.String(), "ok")
}

func (suite *ServerTestSuite) TestHandlerExecutesCorrectActionAndPassesCorrectParams() {
  suite.server.Get("/dummy/:greeting", suite.DummyEndpoint)

  req, _ := http.NewRequest("GET", "/dummy/hi", nil)
  res := httptest.NewRecorder()
  suite.server.Call(res, req)

  assert.Equal(suite.T(), res.Body.String(), "okhi")
}

func (suite *ServerTestSuite) TestHandlerExecutesCorrectActionAndPassesMultipleCorrectParams() {
  suite.server.Get("/dummy/:id/:greeting", suite.DummyEndpoint)

  req, _ := http.NewRequest("GET", "/dummy/5/hola", nil)
  res := httptest.NewRecorder()
  suite.server.Call(res, req)

  assert.Equal(suite.T(), res.Body.String(), "okhola5")
}

func TestServerTestSuite(t *testing.T) {
  suite.Run(t, new(ServerTestSuite))
}