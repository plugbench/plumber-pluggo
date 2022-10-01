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
	if expect.Header != nil {
		for k, vs := range expect.Header {
			actualVs, ok := rt.out.Header[k]
			if !ok {
				rt.t.Errorf("missing expected header %q", k)
				continue
			}
			equivalent := true
			if len(vs) != len(actualVs) {
				equivalent = false
			}
			for i := range vs {
				if vs[i] != actualVs[i] {
					equivalent = false
				}
			}
			if !equivalent {
				rt.t.Errorf("expected %q header values %v, but got %v", k, vs, actualVs)
			}
		}
	}
	return rt
}

func Test_URLs_are_routed_to_schema_specific_topics(t *testing.T) {
	t.Parallel()
	routes(t, &nats.Msg{
		Subject: "plumb.click",
		Data:    []byte("https://eraserhead.net/foo"),
	}).to(&nats.Msg{
		Subject: "cmd.show.url.https",
		Data:    []byte("https://eraserhead.net/foo"),
	})
	routes(t, &nats.Msg{
		Subject: "plumb.click",
		Data:    []byte("http://eraserhead.net/foo"),
	}).to(&nats.Msg{
		Subject: "cmd.show.url.http",
		Data:    []byte("http://eraserhead.net/foo"),
	})
	routes(t, &nats.Msg{
		Subject: "plumb.click",
		Data:    []byte("file://my-workstation/tmp/foo.txt"),
	}).to(&nats.Msg{
		Subject: "cmd.show.url.file",
		Data:    []byte("file://my-workstation/tmp/foo.txt"),
	})
}

func Test_Plumber_passes_through_Base_header(t *testing.T) {
	t.Parallel()
	routes(t, &nats.Msg{
		Subject: "plumb.click",
		Data:    []byte("file://my-workstation/tmp/foo.txt"),
		Header: map[string][]string{
			"Base": {"file://file-server/tmp/"},
		},
	}).to(&nats.Msg{
		Subject: "cmd.show.url.file",
		Data:    []byte("file://my-workstation/tmp/foo.txt"),
		Header: map[string][]string{
			"Base": {"file://file-server/tmp/"},
		},
	})
}

func Test_Plumber_resolves_relative_URLs(t *testing.T) {
	t.Parallel()
	t.Run("absolute path with no server", func(t *testing.T) {
		routes(t, &nats.Msg{
			Subject: "plumb.click",
			Data:    []byte("/tmp/foo.txt"),
			Header: map[string][]string{
				"Base": {"file://file-server/bar/quux"},
			},
		}).to(&nats.Msg{
			Data: []byte("file://file-server/tmp/foo.txt"),
		})
	})
	t.Run("relative path", func(t *testing.T) {
		routes(t, &nats.Msg{
			Subject: "plumb.click",
			Data:    []byte("quux/foo.txt"),
			Header: map[string][]string{
				"Base": {"file://file-server/bar/"},
			},
		}).to(&nats.Msg{
			Data: []byte("file://file-server/bar/quux/foo.txt"),
		})
	})
	t.Run("absolute, non-URL filepath with no Base is made into a file URL", func(t *testing.T) {
		routes(t, &nats.Msg{
			Subject: "plumb.click",
			Data:    []byte("/tmp/foo.txt"),
		}).to(&nats.Msg{
			Data: []byte("file:///tmp/foo.txt"),
		})
	})
}

func Test_Plumber_converts_line_numbers_to_RFC_5147_fragment_ids(t *testing.T) {
	t.Parallel()
	t.Run("exisiting fragment is passed through", func(t *testing.T) {
		routes(t, &nats.Msg{
			Subject: "plumb.click",
			Data:    []byte("file:///tmp/foo.txt#line=42"),
		}).to(&nats.Msg{
			Data: []byte("file:///tmp/foo.txt#line=42"),
		})
	})
	t.Run("colon-separated line number is converted", func(t *testing.T) {
		routes(t, &nats.Msg{
			Subject: "plumb.click",
			Data:    []byte("/tmp/foo.txt:79"),
		}).to(&nats.Msg{
			Data: []byte("file:///tmp/foo.txt#line=78"),
		})
	})
	t.Run("trailing colon is ignored", func(t *testing.T) {
		routes(t, &nats.Msg{
			Subject: "plumb.click",
			Data:    []byte("/tmp/foo.txt:79:"),
		}).to(&nats.Msg{
			Data: []byte("file:///tmp/foo.txt#line=78"),
		})
	})
	// "line-line is converted"
	t.Run("line:column is converted", func(t *testing.T) {
		routes(t, &nats.Msg{
			Subject: "plumb.click",
			Data:    []byte("/tmp/foo.txt:79:12:"),
		}).to(&nats.Msg{
			Data: []byte("file:///tmp/foo.txt#line=78;char=11"),
		})
	})
	// "line:column-column is converted"
}
