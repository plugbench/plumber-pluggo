package plumber

import (
	"github.com/nats-io/nats.go"
	"github.com/plugbench/nats_cli"
)

type Plumber struct {
        natsCfg nats_cli.Config
}

func New(natsCfg nats_cli.Config) (*Plumber, error) {
	return &Plumber{natsCfg: natsCfg}, nil
}

func (p *Plumber) Run() error {
	nc, err := p.natsCfg.Connect()
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
