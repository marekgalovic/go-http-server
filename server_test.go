package server

import (
  "errors";
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

type SampleAuthProvider struct{}

func (p *SampleAuthProvider) Verify(request *Request) error {
  return errors.New("You don't have an access to this resource.")
}

func (suite *ServerTestSuite) SetupSuite() {
  suite.server = NewServer(NewConfig())
}

func (suite *ServerTestSuite) DummyEndpoint(request *Request) {
  request.Respond("ok%s%s", request.Get(":greeting"), request.Get(":id"))
}

func (suite *ServerTestSuite) TestHandleExecutesCorrectAction() {
  suite.server.Get("/dummy", suite.DummyEndpoint, nil)

  req, _ := http.NewRequest("GET", "/dummy", nil)
  res := httptest.NewRecorder()
  suite.server.Handle(res, req)

  assert.Equal(suite.T(), "ok", res.Body.String())
}

func (suite *ServerTestSuite) TestHandleExecutesCorrectActionAndPassesCorrectParams() {
  suite.server.Get("/dummy/:greeting", suite.DummyEndpoint, nil)

  req, _ := http.NewRequest("GET", "/dummy/hi", nil)
  res := httptest.NewRecorder()
  suite.server.Handle(res, req)

  assert.Equal(suite.T(), "okhi", res.Body.String())
}

func (suite *ServerTestSuite) TestHandleExecutesCorrectActionAndPassesMultipleCorrectParams() {
  suite.server.Get("/dummy/:id/:greeting", suite.DummyEndpoint, nil)

  req, _ := http.NewRequest("GET", "/dummy/5/hola", nil)
  res := httptest.NewRecorder()
  suite.server.Handle(res, req)

  assert.Equal(suite.T(), "okhola5", res.Body.String())
}

func (suite *ServerTestSuite) TestHandleReturns401IfAuthenticatedRouteFailsToVerify() {
  suite.server.Get("/protected", suite.DummyEndpoint, &SampleAuthProvider{})

  req, _ := http.NewRequest("GET", "/protected", nil)
  res := httptest.NewRecorder()
  suite.server.Handle(res, req)

  assert.Equal(suite.T(), 401, res.Code)
  assert.Equal(suite.T(), "You don't have an access to this resource.", res.Body.String())
}

func TestServerTestSuite(t *testing.T) {
  suite.Run(t, new(ServerTestSuite))
}