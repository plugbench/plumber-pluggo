package plumber

import (
	"errors"

	"github.com/nats-io/nats.go"
)

type Plumber struct {
}

func New() (*Plumber, error) {
	return &Plumber{}, nil
}

func (p *Plumber) Route(msg *nats.Msg) (*nats.Msg, error) {
	out := nats.NewMsg("browser.open")
	out.Data = msg.Data
	return out, nil
}

func (p *Plumber) Run() error {
	return errors.New("Not implemented")
}
