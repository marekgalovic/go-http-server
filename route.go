package server

import (
  "fmt";
  "regexp";
  "strings"
)

const (
  GET string = "GET";
  POST string = "POST";
  PUT string = "PUT";
  DELETE string = "DELETE";
)

func newRoute(method string, path string, handler func(*Request), authentication AuthProvider) *Route {
  return &Route{Method: method, Path: path, handler: handler, authentication: authentication}
}

type Route struct {
  Method string
  Path string
  Pattern *regexp.Regexp
  ParamNames []string
  handler func(*Request)
  authentication AuthProvider
}

func (r *Route) execute(request *Request) {
  for paramName, paramValue := range r.parseRequestParams(request.Path) {
    request.Params.Set(paramName, paramValue)
  }

  r.handler(request)
}

func (r *Route) compile() error {
  paramNames, paramNamesRegexp, err := r.compilePath()
  if err != nil {
    return err
  }
  r.ParamNames = paramNames
  routeRegex := paramNamesRegexp.ReplaceAllString(r.Path, `([\w-]+)`)
  routeRegex = strings.Replace(routeRegex, "/", "\\/", -1)
  r.Pattern, err = regexp.Compile(fmt.Sprintf("^%s$", routeRegex))
  if err != nil {
    return err
  }
  return nil
}

func (r *Route) compilePath() ([]string, *regexp.Regexp, error) {
  paramNamesRegexp, err := regexp.Compile(`(?::([\w]+))`)
  if err != nil {
    return nil, nil, err
  }
  params := paramNamesRegexp.FindAllString(r.Path, -1)
  return params, paramNamesRegexp, nil
}

func (r *Route) parseRequestParams(path string) map[string]string {
  params := make(map[string]string)
  matched := r.Pattern.FindAllStringSubmatch(path, -1)
  for paramIndex, paramName := range r.ParamNames {
    params[paramName] = matched[0][paramIndex + 1]
  }
  return params
}

func (r *Route) checkAuthentication(request *Request) error {
  if r.isProtected() {
    return r.authentication.Verify(request)
  }
  return nil
}

func (r *Route) isProtected() bool {
  return r.authentication != nil
}

func (r *Route) asString() string {
  var protected string
  if r.isProtected() {
    protected = "PROTECTED"
  }
  return fmt.Sprintf("%s[%s]%s", protected, r.Method, r.Path)
}