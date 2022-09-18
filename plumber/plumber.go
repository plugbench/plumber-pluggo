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
	out.Data = router{msg}.absoluteURL()
	out.Header = msg.Header

	if browserUrl.Match(msg.Data) {
		out.Subject = "browser.open"
	}

	return out, nil
}

func (p *Plumber) Run() error {
	return errors.New("Not implemented")
}

type router struct{ *nats.Msg }

func (msg router) absoluteURL() []byte {
	base := msg.Header.Get("Base")
	if base == "" {
		return msg.Data
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return msg.Data
	}
	absoluteURL, err := baseURL.Parse(string(msg.Data))
	if err != nil {
		return msg.Data
	}
	return []byte(absoluteURL.String())
}
