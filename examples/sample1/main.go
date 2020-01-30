package main

import (
	_ "github.com/deltrinos/tpl21/examples/sample1/app/init.d"
	"github.com/deltrinos/tpl21/log"
)

var (
	Version string
	Build   string
	Module  string
)

func main() {
	log.Debug().Msgf("Version(%s) Build(%s) Module(%s)", Version, Build, Module)
}
