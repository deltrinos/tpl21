package router

import (
	"github.com/deltrinos/tpl21/log"
	"github.com/deltrinos/tpl21/router"
)

var Router *router.Router

func init() {
	log.Debug().Msgf("Initialize router...")

	Router = router.Default()
	Router.NoCache()

	log.Debug().Msgf("Initialize router... DONE.")
}
