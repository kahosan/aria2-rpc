package caller

import (
	"context"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/jhttp"
)

type httpCaller struct {
	host string
	rc   *jrpc2.Client
}

// http call
func (h *httpCaller) call(method string, params, reply any) error {
	ch := jhttp.NewChannel(h.host, nil)
	rc := jrpc2.NewClient(ch, nil)
	defer rc.Close()

	h.rc = rc

	if reply == nil {
		_, err := rc.Call(context.Background(), method, params)
		return err
	}

	return rc.CallResult(context.Background(), method, params, reply)
}

func (h *httpCaller) close() error {
	return h.rc.Close()
}
