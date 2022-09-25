package plumber

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"

	"github.com/nats-io/nats.go"
)

var (
	browserUrl    = regexp.MustCompile(`^https?://`)
	filePositions = regexp.MustCompile(`^(.*?):(\d+)(?::(\d+))?:?$`)

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
		base = "file://"
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return msg.Data
	}

	var line int64
	var haveLine bool
	var col int64
	var haveCol bool
	path := msg.Data
	if sub := filePositions.FindSubmatch(msg.Data); sub != nil {
		path = sub[1]
		line, _ = strconv.ParseInt(string(sub[2]), 10, 64)
		haveLine = true
		if len(sub[3]) > 0 {
			col, _ = strconv.ParseInt(string(sub[3]), 10, 64)
			haveCol = true
		}
	}

	absoluteURL, err := baseURL.Parse(string(path))
	if err != nil {
		return msg.Data
	}

	if haveLine {
		absoluteURL.Fragment = fmt.Sprintf("line=%d", line-1)
		if haveCol {
			absoluteURL.Fragment += fmt.Sprintf(";char=%d", col-1)
		}
	}

	return []byte(absoluteURL.String())
}
