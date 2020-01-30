package env

import (
	"github.com/caarlos0/env"
)

func FillEnv(e interface{}) error {
	err := env.Parse(e)
	if err != nil {
		return err
	}
	return nil
}
