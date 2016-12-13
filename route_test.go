package server

import (
  "testing";
  "net/http";
  "net/http/httptest";

  "github.com/stretchr/testify/suite";
  "github.com/stretchr/testify/assert"
)

type SampleAuthProvider struct{}

func (p *SampleAuthProvider) Verify(request *Request) *Response {
  return request.Response().Error(401, "You don't have an access to this resource.")
}

type RouteTestSuite struct {
  suite.Suite
  route *Route
}

func (suite *RouteTestSuite) SetupSuite() {
  suite.route = newRoute(GET, "/test/:id/resource/:resource_id", func(r *Request) *Response {
    return nil
  }, nil)
}

func (suite *RouteTestSuite) TestCompileSetsCorrectParamNames() {
  expectedParamNames := []string{":id", ":resource_id"}
  suite.route.compile()

  assert.Equal(suite.T(), expectedParamNames, suite.route.ParamNames)
}

func (suite *RouteTestSuite) TestParseRequestParamsReturnsCorrectValues() {
  expected := map[string]string{":id": "5", ":resource_id": "10"}
  actual := suite.route.parseRequestParams("/test/5/resource/10")

  assert.Equal(suite.T(), expected, actual)
}

func (suite *RouteTestSuite) TestExecuteChecksRouteAuthenticationBeforeCallingHandler() {
  protectedRoute := newRoute(GET, "/protected", func(r *Request) *Response {
    return nil
  }, &SampleAuthProvider{})

  req, _ := http.NewRequest("GET", "/test?key=value", nil)
  res := httptest.NewRecorder()
  request := newRequest(req, res, nil)
  response := protectedRoute.execute(request)

  assert.Equal(suite.T(), 401, response.Code)
  assert.Equal(suite.T(), "You don't have an access to this resource.", response.Body)
}

func TestRouteTestSuite(t *testing.T) {
  suite.Run(t, new(RouteTestSuite))
}