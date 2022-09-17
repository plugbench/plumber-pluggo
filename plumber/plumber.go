package nats_plumber

import (
	"errors"

	"github.com/nats-io/nats.go"
)

type Plumber struct {
}

func NewPlumber() (*Plumber, error) {
	return &Plumber{}, nil
}

func (p *Plumber) Route(msg *nats.Msg) (*nats.Msg, error) {
	return nil, errors.New("Not implemented")
}
