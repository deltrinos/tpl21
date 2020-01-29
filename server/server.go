package server

import (
	"github.com/deltrinos/tpl21/db"
	"github.com/deltrinos/tpl21/log"
	"github.com/deltrinos/tpl21/router"
	"github.com/gin-gonic/gin"
	"os"
)

type Server struct {
	Router *router.Router
	DB     *db.DB
	Log    *log.Log
	Quit   chan os.Signal
}

func Default() *Server {
	return &Server{
		Router: router.Default(),
		Log:    log.Default(),
		Quit:   make(chan os.Signal),
	}
}

func DefaultWithSession() *Server {
	return &Server{
		Router: router.DefaultWithSession(),
		Log:    log.Default(),
		Quit:   make(chan os.Signal),
	}
}

func (s *Server) WithDB(connType, connStr string) *Server {
	s.DB = db.Default(connType, connStr)
	return s
}

func (s *Server) AddMiddleware(middleware ...gin.HandlerFunc) {
	s.Router.Engine.Use(middleware...)
}

func (s *Server) Start() {
	s.Router.Start()
}
