package service

import (
  "fmt";
  "errors";
  "regexp";
  "net/http";
  "path/filepath";

  log "github.com/Sirupsen/logrus";
)

func NewServer(config *Config) *Server {
  return &Server{
    config: config,
    routes: make(map[string]map[*regexp.Regexp]*Route),
  }
}

type Server struct {
  config *Config
  routes map[string]map[*regexp.Regexp]*Route
}

func (s *Server) Listen() {
  http.HandleFunc("/", s.Call)
  log.Info("Server listening on: ", s.address())
  http.ListenAndServe(s.address(), nil)
}

func (s *Server) Static(route string, dirName string) error {
  if string(route[len(route)-1]) != "/" {
    route = fmt.Sprintf("%s/", route)
  }
  http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
    requestedPath := r.URL.Path[len(route):]
    fullRequestedPath := filepath.Join(config.Get().Root, dirName, requestedPath)
    http.ServeFile(w, r, fullRequestedPath)
  })
  log.Infof("Registered static handler [%s] %s", route, dirName)
  return nil
}

func (s *Server) Get(route string, handler func(*Request)) error {
  return s.setRoute(GET, route, handler)
}

func (s *Server) Post(route string, handler func(*Request)) error {
  return s.setRoute(POST, route, handler)
}

func (s *Server) Put(route string, handler func(*Request)) error {
  return s.setRoute(PUT, route, handler)
}

func (s *Server) Delete(route string, handler func(*Request)) error {
  return s.setRoute(DELETE, route, handler)
}

func (s *Server) address() string {
  return fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port)
}

func (s *Server) getRoute(method string, path string) *Route {
  for routePattern, route := range s.routes[method] {
    if routePattern.MatchString(path) {
      return route
    }
  }
  return nil
}

func (s *Server) setRoute(method string, path string, handler func(*Request)) error {
  route := newRoute(method, path, handler)
  err := route.compile()
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to compile route %s", path))
  }
  s.ensureMapForMethod(route.Method)
  s.routes[route.Method][route.Pattern] = route
  log.Infof("Registered route [%s]%s", method, path)
  return nil
}

func (s *Server) ensureMapForMethod(method string) {
  methodMap := s.routes[method]
  if methodMap == nil {
    s.routes[method] = make(map[*regexp.Regexp]*Route)
  }
}

func (s *Server) Call(w http.ResponseWriter, req *http.Request) {
  req.ParseForm()
  request := newRequest(req.Method, req.URL.Path, req.Form, req.Body, w)
  if route := s.getRoute(req.Method, req.URL.Path); route != nil {
    go route.execute(request)
  } else {
    go request.Error(404, fmt.Sprintf("Route [%s]%s does not exist.", req.Method, req.URL.Path))
  }
  request.writeResponse()
}