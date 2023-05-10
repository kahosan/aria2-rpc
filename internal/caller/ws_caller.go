package caller

import (
	"context"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/wschannel"
)

type wsCaller struct {
	host  string
	rc    *jrpc2.Client
	close func() error
}

func newWsCaller(host string) (*wsCaller, error) {
	ch, err := wschannel.Dial(host, nil)
	if err != nil {
		return nil, err
	}
	rc := jrpc2.NewClient(ch, nil)

	c := &wsCaller{
		host:  host,
		rc:    rc,
		close: rc.Close,
	}

	return c, nil
}

// websocket call
func (w *wsCaller) call(method string, params, reply any) error {
	if reply == nil {
		_, err := w.rc.Call(context.Background(), method, params)
		return err
	}

	return w.rc.CallResult(context.Background(), method, params, reply)
}
