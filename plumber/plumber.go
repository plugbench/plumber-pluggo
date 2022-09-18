package plumber

import (
	"errors"
	"net/url"
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

	out.Data = msg.Data
	if base := msg.Header.Get("Base"); base != "" {
		baseURL, err := url.Parse(base)
		if err == nil {
			absoluteURL, err := baseURL.Parse(string(msg.Data))
			if err == nil {
				out.Data = []byte(absoluteURL.String())
			}
		}
	}

	if browserUrl.Match(msg.Data) {
		out.Subject = "browser.open"
	}

	return out, nil
}

func (p *Plumber) Run() error {
	return errors.New("Not implemented")
}
