package service

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

func newRoute(method string, path string, handler func(*Request)) *Route {
  return &Route{Method: method, Path: path, handler: handler}
}

type Route struct {
  Path string
  Pattern *regexp.Regexp
  Method string
  ParamNames []string
  handler func(*Request)
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