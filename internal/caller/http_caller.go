package caller

import (
	"context"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/jhttp"
)

type httpCaller struct {
	host  string
	rc    *jrpc2.Client
	close func() error
}

func newHttpCaller(host string) (*httpCaller, error) {
	ch := jhttp.NewChannel(host, nil)
	rc := jrpc2.NewClient(ch, nil)

	c := &httpCaller{
		host:  host,
		rc:    rc,
		close: rc.Close,
	}

	return c, nil
}

// http call
func (h *httpCaller) call(method string, params, reply any) error {
	if reply == nil {
		_, err := h.rc.Call(context.Background(), method, params)
		return err
	}

	return h.rc.CallResult(context.Background(), method, params, reply)
}
