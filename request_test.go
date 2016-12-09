package server

import (
  "testing";
  "net/http";
  "net/http/httptest";

  "github.com/stretchr/testify/suite";
  "github.com/stretchr/testify/assert"
)

type RequestTestSuite struct{
  suite.Suite
  request *Request
  res *httptest.ResponseRecorder
}

func (suite *RequestTestSuite) SetupTest() {
  req, _ := http.NewRequest("GET", "/test?key=value", nil)
  suite.res = httptest.NewRecorder()
  suite.request = newRequest(req, suite.res)
}

func (suite *RequestTestSuite) TestRespondWritesPlainTextDataAsResponse() {
  go suite.request.Respond("response")
  suite.request.response.write()

  assert.Equal(suite.T(), "response", suite.res.Body.String())
}

func (suite *RequestTestSuite) TestRespondJsonWritesJsonStringAsResponseAndSetsCorrectContentTypeHeader() {
  go suite.request.RespondJson(&map[string]string{"foo": "bar"})
  suite.request.response.write()

  assert.Equal(suite.T(), "application/json", suite.res.HeaderMap.Get("Content-Type"))
  assert.Equal(suite.T(), `{"foo":"bar"}`, suite.res.Body.String())
}

func (suite *RequestTestSuite) TestErrorWritesErrorStringAsAResonseAndSetsCorrectCode() {
  go suite.request.Error(404, "Resource not found")
  suite.request.response.write()

  assert.Equal(suite.T(), 404, suite.res.Code)
  assert.Equal(suite.T(), "Resource not found", suite.res.Body.String())
}

func TestRequestTestSuite(t *testing.T) {
  suite.Run(t, new(RequestTestSuite))
}