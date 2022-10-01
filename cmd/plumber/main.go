package main

import (
	"log"

	"github.com/plugbench/plumber-pluggo/plumber"
)

func main() {
	p, err := plumber.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
