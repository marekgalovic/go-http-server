package server

import (
  "os";
  "fmt";
  "log";
  "errors";
  "regexp";
  "net/http";
  "path/filepath";
)

var Logger = log.New(os.Stdout, "[http-server] ", log.Flags())

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

func (s *Server) Listen() error {
  http.HandleFunc("/", s.Handle)

  Logger.Printf("Listening on: %s, TLS: %t", s.bindAddress(), s.usesTls())
  return s.bindListener()
}

func (s *Server) bindListener() error {
  if s.usesTls() {
    return http.ListenAndServeTLS(s.bindAddress(), s.config.CertFile, s.config.KeyFile, nil)
  }
  return http.ListenAndServe(s.bindAddress(), nil)
}

func (s *Server) bindAddress() string {
  return fmt.Sprintf("%s:%d", s.config.Address, s.config.Port)
}

func (s *Server) usesTls() bool {
  return s.config.CertFile != "" && s.config.KeyFile != ""
}

func (s *Server) Static(route string, dirName string, authentication AuthProvider) error {
  if string(route[len(route)-1]) != "/" {
    route = fmt.Sprintf("%s/", route)
  }
  http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
    requestedPath := r.URL.Path[len(route):]
    fullRequestedPath := filepath.Join(s.config.StaticRoot, dirName, requestedPath)
    http.ServeFile(w, r, fullRequestedPath)
  })
  Logger.Printf("Registered static handler [%s]%s", route, dirName)
  return nil
}

func (s *Server) Get(route string, handler func(*Request), authentication AuthProvider) error {
  return s.setRoute(newRoute(GET, route, handler, authentication))
}

func (s *Server) Post(route string, handler func(*Request), authentication AuthProvider) error {
  return s.setRoute(newRoute(POST, route, handler, authentication))
}

func (s *Server) Put(route string, handler func(*Request), authentication AuthProvider) error {
  return s.setRoute(newRoute(PUT, route, handler, authentication))
}

func (s *Server) Delete(route string, handler func(*Request), authentication AuthProvider) error {
  return s.setRoute(newRoute(DELETE, route, handler, authentication))
}

func (s *Server) getRoute(method string, path string) *Route {
  for routePattern, route := range s.routes[method] {
    if routePattern.MatchString(path) {
      return route
    }
  }
  return nil
}

func (s *Server) setRoute(route *Route) error {
  err := route.compile()
  if err != nil {
    return errors.New(fmt.Sprintf("Failed to compile route %s", route.Path))
  }
  s.ensureMapForMethod(route.Method)
  s.routes[route.Method][route.Pattern] = route
  Logger.Printf("Registered route %s", route.asString())
  return nil
}

func (s *Server) ensureMapForMethod(method string) {
  methodMap := s.routes[method]
  if methodMap == nil {
    s.routes[method] = make(map[*regexp.Regexp]*Route)
  }
}

func (s *Server) Handle(w http.ResponseWriter, req *http.Request) {
  req.ParseForm()
  request := newRequest(req, w)
  route := s.getRoute(req.Method, req.URL.Path)
  if route != nil {
    err := route.checkAuthentication(request)
    if err != nil {
      go request.Error(401, err.Error())
    } else {
      go route.execute(request)
    }
  } else {
    go request.Error(404, fmt.Sprintf("Route [%s]%s does not exist.", req.Method, req.URL.Path))
  }
  response := request.response.write()
  s.logHandlerResult(request, response)
}

func (s *Server) logHandlerResult(request *Request, response *Response) {
  Logger.Printf(
    "Request. METHOD: %s, PATH: %s, REMOTE_ADDR: %s, RESPONSE_CODE: %d, DURATION: %s",
    request.Method, request.Path, request.RemoteAddr, response.Code, response.Duration,
  )
}