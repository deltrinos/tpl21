package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

func DefaultWithSession() *Router {
	r := Default()

	store := cookie.NewStore([]byte("my-project"))
	r.Engine.Use(sessions.Sessions("my-session", store))
	return r
}
