package plumber

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"

	"github.com/nats-io/nats.go"
)

var (
	filePositions = regexp.MustCompile(`^(.*?):(\d+)(?::(\d+))?(?::[^:]*)?$`)
)

type routeAction struct {
	msg  *nats.Msg
	send func(msg *nats.Msg) error
}

func newRouteCommand(msg *nats.Msg, send func(msg *nats.Msg) error) *routeAction {
	return &routeAction{msg: msg, send: send}
}

func (a *routeAction) route(msg *nats.Msg) (*nats.Msg, error) {
	u, err := a.absoluteURL()
	if err != nil {
		return nil, err
	}
	out := nats.NewMsg(fmt.Sprintf("cmd.show.url.%s", u.Scheme))
	out.Data = []byte(u.String())
	out.Header = msg.Header
	out.Reply = msg.Reply
	return out, nil
}

func (a *routeAction) absoluteURL() (*url.URL, error) {
	base := a.msg.Header.Get("Base")
	if base == "" {
		base = "file://"
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	var line int64
	var haveLine bool
	var col int64
	var haveCol bool
	path := a.msg.Data
	if sub := filePositions.FindSubmatch(a.msg.Data); sub != nil {
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
		return nil, err
	}

	if haveLine {
		absoluteURL.Fragment = fmt.Sprintf("line=%d", line-1)
		if haveCol {
			absoluteURL.Fragment += fmt.Sprintf(".%d", col-1)
		}
	}

	return absoluteURL, nil
}

func (a *routeAction) Execute() {
	log.Printf("recieved %q", string(a.msg.Data))
	next, err := a.route(a.msg)
	if err == nil {
		err = a.send(next)
	}
	if err != nil {
		log.Print(err)
		errReply := nats.NewMsg(a.msg.Reply)
		errReply.Data = []byte(fmt.Sprintf("ERROR: %v", err.Error()))
		if err := a.send(errReply); err != nil {
			log.Printf("error responding: %v", err)
		}
	}
}
