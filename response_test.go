package server

import (
  "testing";
  "net/http";
  "net/http/httptest";

  "github.com/stretchr/testify/suite";
  "github.com/stretchr/testify/assert"
)

type ResponseTestSuite struct {
  suite.Suite
  res *httptest.ResponseRecorder
  request *Request
  response *Response
}

func (suite *ResponseTestSuite) SetupTest() {
  req, _ := http.NewRequest("GET", "/test?key=value", nil)
  suite.res = httptest.NewRecorder()
  suite.request = newRequest(req, suite.res, nil)
  suite.response = newResponse(suite.request)
}

func (suite *ResponseTestSuite) TestPlainSetsPlainTextResponse() {
  suite.response.Plain("Hi there!").write()

  assert.Equal(suite.T(), 200, suite.res.Code)
  assert.Equal(suite.T(), "Hi there!", suite.res.Body.String())
}

func (suite *ResponseTestSuite) TestJsonSerializesResponseToAJsonString() {
  suite.response.Json(&map[string]string{"foo": "bar"}).write()

  assert.Equal(suite.T(), "application/json", suite.res.HeaderMap.Get("Content-Type"))
  assert.Equal(suite.T(), `{"foo":"bar"}`, suite.res.Body.String())
}

func (suite *ResponseTestSuite) TestErrorSetsAResponseWithProvidedErrorCode() {
  suite.response.Error(404, "Not found.").write()

  assert.Equal(suite.T(), 404, suite.res.Code)
  assert.Equal(suite.T(), "Not found.", suite.res.Body.String())
}

func (suite *ResponseTestSuite) TestSetCodeSetsAResponseCode() {
  suite.response.SetCode(501).Plain("").write()

  assert.Equal(suite.T(), 501, suite.res.Code)
}

func (suite *ResponseTestSuite) TestSetHeaderSetsAResponseHeader() {
  suite.response.SetHeader("Keep-Alive", "timeout=5, max=1000").Plain("").write()

  assert.Equal(suite.T(), "timeout=5, max=1000", suite.res.HeaderMap.Get("Keep-Alive"))
}

func TestResponseTestSuite(t *testing.T) {
  suite.Run(t, new(ResponseTestSuite))
}