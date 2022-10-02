package plumber

import (
	"github.com/nats-io/nats.go"
)

type Plumber struct {
}

func New() (*Plumber, error) {
	return &Plumber{}, nil
}

func (p *Plumber) Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	ch := make(chan *nats.Msg, 32)
	sub, err := nc.ChanSubscribe("cmd.show.data.plumb", ch)
	if err != nil {
		return err
	}
	defer sub.Drain()

	for {
		msg := <-ch
		cmd := newRouteCommand(msg, nc.PublishMsg)
		cmd.Execute()
	}
}

