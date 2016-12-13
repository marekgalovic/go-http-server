package server

import (
  "os";
  "fmt";
  "log";
  "errors";
  "regexp";
  "net/http";
  "github.com/gorilla/sessions"
)

var Logger = log.New(os.Stdout, "[http-server] ", log.Flags())

func NewServer(config *Config) *Server {
  return &Server{
    config: config,
    routes: make(map[string]map[*regexp.Regexp]*Route),
    notifications: make(chan *Notification),
  }
}

type Server struct {
  config *Config
  routes map[string]map[*regexp.Regexp]*Route
  sessions sessions.Store
  notifications chan *Notification
}

func (s *Server) Listen() error {
  http.HandleFunc("/", s.Handle)

  Logger.Printf("Listening on: %s, TLS: %t", s.bindAddress(), s.usesTls())
  return s.bindListener()
}

func (s *Server) Get(route string, handler func(*Request)*Response, authentication AuthProvider) error {
  return s.setRoute(newRoute(GET, route, handler, authentication))
}

func (s *Server) Post(route string, handler func(*Request)*Response, authentication AuthProvider) error {
  return s.setRoute(newRoute(POST, route, handler, authentication))
}

func (s *Server) Put(route string, handler func(*Request)*Response, authentication AuthProvider) error {
  return s.setRoute(newRoute(PUT, route, handler, authentication))
}

func (s *Server) Delete(route string, handler func(*Request)*Response, authentication AuthProvider) error {
  return s.setRoute(newRoute(DELETE, route, handler, authentication))
}

func (s *Server) UseSessionStore(store sessions.Store) {
  s.sessions = store
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
  request := newRequest(req, w, s)
  route := s.getRoute(req.Method, req.URL.Path)

  var response *Response
  if route != nil {
    response = route.execute(request)
  } else {
    response = request.Response().Error(404, "Route [%s]%s does not exist.", req.Method, req.URL.Path)
  }

  response.write()
  go s.logHandlerResult(request, response)
  go s.sendNotification(request, response)
}

func (s *Server) logHandlerResult(request *Request, response *Response) {
  Logger.Printf(
    "%s. METHOD: %s, PATH: %s, REMOTE_ADDR: %s, RESPONSE_CODE: %d, DURATION: %s",
    response.responseType, request.Method, request.Path, request.RemoteAddr, response.Code, response.Duration,
  )
}

func (s *Server) Notifications() <- chan *Notification {
  return s.notifications
}

func (s *Server) sendNotification(request *Request, response *Response) {
  select {
  case s.notifications <- &Notification{Request: request, Response: response}:
  default:
  }
}