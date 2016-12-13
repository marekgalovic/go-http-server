package server

import (
  "bytes";
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

type SampleRequestBody struct {
  Foo string
}

func (suite *RequestTestSuite) SetupTest() {
  req, _ := http.NewRequest("GET", "/test?key=value", nil)
  suite.res = httptest.NewRecorder()
  suite.request = newRequest(req, suite.res, nil)
}

func (suite *RequestTestSuite) TestGetReturnsParamValue() {
  assert.Equal(suite.T(), "value", suite.request.Get("key"))
}

func (suite *RequestTestSuite) TestJsonScansBodyContentToStruct() {
  req, _ := http.NewRequest("POST", "/resource", bytes.NewBuffer([]byte(`{"foo": "bar"}`)))
  req.Header.Set("Content-Type", "application/json")
  request := newRequest(req, suite.res, nil)

  var body *SampleRequestBody
  request.Json(&body)

  assert.Equal(suite.T(), "bar", body.Foo)
}

func TestRequestTestSuite(t *testing.T) {
  suite.Run(t, new(RequestTestSuite))
}