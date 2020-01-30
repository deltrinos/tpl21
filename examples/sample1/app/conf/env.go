package conf

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/deltrinos/tpl21/env"
	"github.com/deltrinos/tpl21/log"
)

type AppEnv struct {
	ConnStr  string `env:"CONN_STR" envDefault:"host=localhost port=5432 user=postgres dbname=postgres password=postgres"`
	ConnType string `env:"CONN_TYPE" envDefault:"postgres"`
	Addr     string `env:"ADDR" envDefault:":3000"`
}

var Env *AppEnv

func init() {
	Env = &AppEnv{}

	log.Debug().Msgf("scanning application env...")

	err := env.FillEnv(Env)
	if err != nil {
		log.Error().Msgf("failed to FillEnv variables: %v", err)
	} else {
		log.Debug().Msg(spew.Sdump(Env))
	}

	log.Debug().Msgf("scanning application env... DONE.")
}
