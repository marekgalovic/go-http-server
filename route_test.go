package server

import (
  "testing";

  "github.com/stretchr/testify/suite";
  "github.com/stretchr/testify/assert"
)

type RouteTestSuite struct {
  suite.Suite
  route *Route
  protectedRoute *Route
}

func (suite *RouteTestSuite) SetupSuite() {
  suite.route = newRoute(GET, "/test/:id/resource/:resource_id", func(r *Request) {}, nil)
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

func TestRouteTestSuite(t *testing.T) {
  suite.Run(t, new(RouteTestSuite))
}