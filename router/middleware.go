package router

import "github.com/gin-gonic/gin"

func (r *Router) AddMiddleware(middleware ...gin.HandlerFunc) {
	r.Engine.Use(middleware...)
}
