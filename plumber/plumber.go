package plumber

import (
	"errors"
	"fmt"
	"log"
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
	out.Reply = msg.Reply

	if browserUrl.Match(msg.Data) {
		out.Subject = "browser.open"
	}

	return out, nil
}

func (p *Plumber) Run() error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}
	defer nc.Close()

	ch := make(chan *nats.Msg, 32)
	sub, err := nc.ChanSubscribe("plumb.click", ch)
	if err != nil {
		return err
	}
	defer sub.Drain()

	for {
		msg := <-ch
		log.Printf("recieved %q", string(msg.Data))

		next, err := p.Route(msg)
		if err == nil {
			err = nc.PublishMsg(next)
		}
		if err != nil {
			log.Print(err)
			if err := msg.Respond([]byte(fmt.Sprintf("ERROR: %v", err.Error()))); err != nil {
				log.Printf("error responding: %v", err)
			}
			continue
		}
	}
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
