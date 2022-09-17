package plumber

import (
	"testing"

	"github.com/nats-io/nats.go"
)

type routeTest struct {
	t   *testing.T
	out *nats.Msg
	err error
}

type msgOption = func(r *nats.Msg)

func data(s string) msgOption { return func(msg *nats.Msg) { msg.Data = []byte(s) } }

func msg(subject string, opts ...msgOption) *nats.Msg {
	msg := nats.NewMsg(subject)
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

func routes(t *testing.T, msg *nats.Msg) *routeTest {
	p, err := New()
	if err != nil {
		t.Fatal(err)
	}
	out, err := p.Route(msg)
	if err != nil {
		t.Fatal(err)
	}
	rt := &routeTest{t: t, out: out, err: err}
	return rt
}

func (rt *routeTest) to(subject string) *routeTest {
	if rt.out.Subject != subject {
		rt.t.Errorf("expected subject %q, but got %q", subject, rt.out.Subject)
	}
	return rt
}

func Test_HTTPS_and_HTTP_URLs_go_to_the_browser(t *testing.T) {
	routes(t, msg("plumb.click", data("https://eraserhead.net/foo"))).to("browser.open")
	routes(t, msg("plumb.click", data("http://eraserhead.net/foo"))).to("browser.open")
}
