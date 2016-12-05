package service

import (
  "testing";
  "net/url";
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
  suite.res = httptest.NewRecorder()
  queryValues, _ := url.ParseQuery("?key=value")
  suite.request = newRequest("GET", "/test", queryValues, nil, suite.res)
}

func (suite *RequestTestSuite) TestRespondWritesPlainTextDataAsResponse() {
  go suite.request.Respond("response")
  suite.request.writeResponse()

  assert.Equal(suite.T(), suite.res.Body.String(), "response")
}

func (suite *RequestTestSuite) TestRespondJsonWritesJsonStringAsResponseAndSetsCorrectContentTypeHeader() {
  go suite.request.RespondJson(&map[string]string{"foo": "bar"})
  suite.request.writeResponse()

  assert.Equal(suite.T(), suite.res.HeaderMap.Get("Content-Type"), "application/json")
  assert.Equal(suite.T(), suite.res.Body.String(), `{"foo":"bar"}`)
}

func (suite *RequestTestSuite) TestErrorWritesErrorStringAsAResonseAndSetsCorrectCode() {
  go suite.request.Error(404, "Resource not found.")
  suite.request.writeResponse()

  assert.Equal(suite.T(), suite.res.Code, 404)
  assert.Equal(suite.T(), suite.res.Body.String(), "Resource not found.")
}

func TestRequestTestSuite(t *testing.T) {
  suite.Run(t, new(RequestTestSuite))
}