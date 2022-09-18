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

func (rt *routeTest) to(expect *nats.Msg) *routeTest {
	if rt.err != nil {
		rt.t.Errorf("expected routing to succeed, but failed with %v", rt.err)
		return rt
	}
	if expect.Subject != "" && rt.out.Subject != expect.Subject {
		rt.t.Errorf("expected subject %q, but got %q", expect.Subject, rt.out.Subject)
	}
	if expect.Data != nil && string(rt.out.Data) != string(expect.Data) {
		rt.t.Errorf("expected data %q, but got %q", string(expect.Data), string(rt.out.Data))
	}
	return rt
}

func Test_HTTPS_and_HTTP_URLs_go_to_the_browser(t *testing.T) {
	routes(t, &nats.Msg{
		Subject: "plumb.click",
		Data:    []byte("https://eraserhead.net/foo"),
	}).to(&nats.Msg{
		Subject: "browser.open",
		Data:    []byte("https://eraserhead.net/foo"),
	})
	routes(t, &nats.Msg{
		Subject: "plumb.click",
		Data:    []byte("http://eraserhead.net/foo"),
	}).to(&nats.Msg{
		Subject: "browser.open",
		Data:    []byte("http://eraserhead.net/foo"),
	})
}

func Test_Absolute_paths_are_routed_to_the_editor(t *testing.T) {
	routes(t, &nats.Msg{
		Subject: "plumb.click",
		Data:    []byte("/tmp/foo.txt"),
	}).to(&nats.Msg{
		Subject: "editor.open",
		Data:    []byte("/tmp/foo.txt"),
	})
}
