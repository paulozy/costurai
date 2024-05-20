package server

import (
	"github.com/gin-gonic/gin"
	"github.com/paulozy/costurai/internal/infra/server/middlewares"
)

type Handler struct {
	Path   string
	Method string
	Auth   bool
	Func   gin.HandlerFunc
}

type Server struct {
	Host     string
	Port     string
	Router   *gin.Engine
	Handlers []Handler
}

func NewServer(host, port string) *Server {
	server := &Server{
		Host:   host,
		Port:   port,
		Router: gin.Default(),
	}

	return server
}

func (s *Server) AddHandlers() {
	s.Handlers = append(s.Handlers, Routes...)
}

func (s *Server) Start() {
	for _, h := range s.Handlers {
		if h.Auth {
			s.Router.Handle(h.Method, h.Path, middlewares.EnsureAuthenticated(), h.Func)
		} else {
			s.Router.Handle(h.Method, h.Path, h.Func)
		}
	}

	s.Router.Run(s.Host + ":" + s.Port)
}
