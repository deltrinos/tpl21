package api

import (
	"github.com/deltrinos/tpl21/examples/sample1/app/router"
	"github.com/deltrinos/tpl21/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	router.Router.Engine.GET("/api/test", func(c *gin.Context) {
		log.Debug().Msgf("test")
		c.JSON(http.StatusOK, map[string]string{
			"test": "42",
		})
	})
}
