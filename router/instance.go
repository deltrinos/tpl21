package router

import (
	"context"
	"fmt"
	"github.com/deltrinos/tpl21/db"
	"github.com/deltrinos/tpl21/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Router struct {
	DB          *db.DB
	Engine      *gin.Engine
	Quit        chan os.Signal
	StartUpTime time.Time
}

type Handle struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

func Default() *Router {
	r := gin.Default()

	// handle 404 to /
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		prefix := "/"
		if strings.HasPrefix(path, prefix) {
			path = path[len(prefix):]
		}
		c.Redirect(http.StatusFound, prefix+"/?v="+time.Now().Format(time.RFC3339Nano)+"&redirect="+url.QueryEscape(path))
	})

	return &Router{
		Engine:      r,
		StartUpTime: time.Time{},
		Quit:        make(chan os.Signal),
	}
}

func (r *Router) Start() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Info().Msgf("Start server at %s...", addr)

	srv := http.Server{
		Addr:    addr,
		Handler: r.Engine,
	}

	go func() {
		r.StartUpTime = time.Now()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msgf("Failed to srv.ListenAndServe: %v\n", err)
		}
	}()

	signal.Notify(r.Quit, os.Interrupt)
	<-r.Quit

	log.Info().Msgf("Shutting down Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed to srv.Shutdown: %v\n", err)
	}
	log.Info().Msgf("Server exited.")
}

func (r *Router) Stop() {
	r.Quit <- os.Interrupt
}
