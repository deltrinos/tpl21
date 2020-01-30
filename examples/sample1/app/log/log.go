package log

import "github.com/deltrinos/tpl21/log"

func init() {
	log.Default()
	log.Debug().Msgf("log switch to default mode")
}
