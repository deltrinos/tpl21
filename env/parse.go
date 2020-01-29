package env

import (
	"github.com/caarlos0/env"
	"github.com/deltrinos/tpl21/log"
)

func FillEnv(e interface{}) {
	err := env.Parse(e)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to env.Parse %v", err)
	}
}
