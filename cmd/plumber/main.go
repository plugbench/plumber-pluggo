package main

import (
	"log"

	"github.com/plugbench/nats_cli"

	"github.com/plugbench/plumber-pluggo/plumber"
)

func main() {
	natsCfg, err := nats_cli.LoadConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	p, err := plumber.New(natsCfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
