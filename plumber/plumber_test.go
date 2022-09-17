package plumber

import (
	"testing"

	"github.com/nats-io/nats.go"
)

type routeTest struct {
	t  *testing.T
	in *nats.Msg
}

type routeOption = func(r *routeTest)

func data(s string) routeOption { return func(rt *routeTest) { rt.in.Data = []byte(s) } }

func msg(t *testing.T, subject string, opts ...routeOption) *routeTest {
	rt := &routeTest{
		t:  t,
		in: nats.NewMsg(subject),
	}
	for _, opt := range opts {
		opt(rt)
	}
	return rt
}

func (rt *routeTest) routesTo(subject string) *routeTest {
	p, err := New()
	if err != nil {
		rt.t.Fatal(err)
	}
	out, err := p.Route(rt.in)
	if err != nil {
		rt.t.Fatal(err)
	}
	if out.Subject != subject {
		rt.t.Errorf("expected subject %q, but got %q", subject, out.Subject)
	}
	return rt
}

func Test_HTTPS_and_HTTP_URLs_go_to_the_browser(t *testing.T) {
	msg(t, "plumb.click", data("https://eraserhead.net/foo")).routesTo("browser.open")
	msg(t, "plumb.click", data("http://eraserhead.net/foo")).routesTo("browser.open")
}
