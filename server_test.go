package server

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

func (suite *ServerTestSuite) SetupTest() {
  suite.server = NewServer(NewConfig())
}

func (suite *ServerTestSuite) DummyEndpoint(request *Request) *Response {
  return request.Response().Plain("ok%s%s", request.Get(":greeting"), request.Get(":id"))
}

func (suite *ServerTestSuite) MockHandler(pattern string, path string) *httptest.ResponseRecorder {
  suite.server.Get(pattern, suite.DummyEndpoint, nil)

  req, _ := http.NewRequest("GET", path, nil)
  res := httptest.NewRecorder()
  suite.server.Handle(res, req)
  return res
}

func (suite *ServerTestSuite) TestHandleExecutesCorrectAction() {
  res := suite.MockHandler("/dummy", "/dummy")

  assert.Equal(suite.T(), "ok", res.Body.String())
}

func (suite *ServerTestSuite) TestHandleExecutesCorrectActionAndPassesCorrectParams() {
  res := suite.MockHandler("/dummy/:greeting", "/dummy/hi")

  assert.Equal(suite.T(), "okhi", res.Body.String())
}

func (suite *ServerTestSuite) TestHandleExecutesCorrectActionAndPassesMultipleCorrectParams() {
  res := suite.MockHandler("/dummy/:id/:greeting", "/dummy/5/hola")

  assert.Equal(suite.T(), "okhola5", res.Body.String())
}

func (suite *ServerTestSuite) TestHandleMatchesWildcartRoutes() {
  res := suite.MockHandler("/dummy/*/:greeting", "/dummy/wildcard/hola")

  assert.Equal(suite.T(), "okhola", res.Body.String())
}

func TestServerTestSuite(t *testing.T) {
  suite.Run(t, new(ServerTestSuite))
}