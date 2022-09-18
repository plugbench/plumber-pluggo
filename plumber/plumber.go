package plumber

import (
	"errors"
	"regexp"

	"github.com/nats-io/nats.go"
)

var (
	browserUrl = regexp.MustCompile(`^https?://`)

	NoRoute = errors.New("no route")
)

type Plumber struct {
}

func New() (*Plumber, error) {
	return &Plumber{}, nil
}

func (p *Plumber) Route(msg *nats.Msg) (*nats.Msg, error) {
	if browserUrl.Match(msg.Data) {
		out := nats.NewMsg("browser.open")
		out.Data = msg.Data
		return out, nil
	}

	out := nats.NewMsg("editor.open")
	out.Data = msg.Data
	return out, nil
}

func (p *Plumber) Run() error {
	return errors.New("Not implemented")
}
