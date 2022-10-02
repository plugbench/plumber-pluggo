package plumber

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type routeAction struct {
	msg  *nats.Msg
	send func(msg *nats.Msg) error
}

func newRouteCommand(msg *nats.Msg, send func(msg *nats.Msg) error) *routeAction {
	return &routeAction{msg: msg, send: send}
}

func route(msg *nats.Msg) (*nats.Msg, error) {
	u := router{msg}.absoluteURL()
	out := nats.NewMsg(fmt.Sprintf("cmd.show.url.%s", u.Scheme))
	out.Data = []byte(u.String())
	out.Header = msg.Header
	out.Reply = msg.Reply
	return out, nil
}

func (rc *routeAction) Execute() {
	log.Printf("recieved %q", string(rc.msg.Data))
	next, err := route(rc.msg)
	if err == nil {
		err = rc.send(next)
	}
	if err != nil {
		log.Print(err)
		errReply := nats.NewMsg(rc.msg.Reply)
		errReply.Data = []byte(fmt.Sprintf("ERROR: %v", err.Error()))
		if err := rc.send(errReply); err != nil {
			log.Printf("error responding: %v", err)
		}
	}
}
