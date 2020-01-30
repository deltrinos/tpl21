package router

import (
	"github.com/gin-gonic/gin"
)

func (r *Router) NoCache() {
	r.Engine.Use(func(c *gin.Context) {
		c.Header("Cache-control", "no-cache, no-store, must-revalidate")
	})
}
