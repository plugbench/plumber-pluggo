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
	out := nats.NewMsg("editor.open")
	out.Data = msg.Data
	out.Header = msg.Header

	if browserUrl.Match(msg.Data) {
		out.Subject = "browser.open"
	}

	return out, nil
}

func (p *Plumber) Run() error {
	return errors.New("Not implemented")
}
