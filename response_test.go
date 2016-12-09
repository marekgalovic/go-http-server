package server

import (
  "testing";
  "net/http/httptest";

  "github.com/stretchr/testify/suite";
  "github.com/stretchr/testify/assert"
)

type ResponseTestSuite struct {
  suite.Suite
  writer *httptest.ResponseRecorder
  response *Response
}

func (suite *ResponseTestSuite) SetupTest() {
  suite.writer = httptest.NewRecorder()
  suite.response = newResponse(suite.writer)
}

func (suite *ResponseTestSuite) TestSetCodeWritesResponseCode() {
  go func() {
    suite.response.setCode(302)
    suite.response.setBody("")
  }()
  suite.response.write()

  assert.Equal(suite.T(), 302, suite.response.Code)
  assert.Equal(suite.T(), 302, suite.writer.Code)
}

func (suite *ResponseTestSuite) TestSetBodyWritesResponseBody() {
  go func() {
    suite.response.setBody("test")
  }()
  suite.response.write()

  assert.Equal(suite.T(), "test", suite.response.Body)
  assert.Equal(suite.T(), "test", suite.writer.Body.String())
}

func (suite *ResponseTestSuite) TestSetHeaderWritesResponseHeader() {
  go func() {
    suite.response.setHeader("X-Custom", "value")
    suite.response.setBody("")
  }()
  suite.response.write()

  assert.Equal(suite.T(), "value", suite.writer.HeaderMap.Get("X-Custom"))
}

func TestResponseTestSuite(t *testing.T) {
  suite.Run(t, new(ResponseTestSuite))
}