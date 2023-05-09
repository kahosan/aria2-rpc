package caller

import (
	"context"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/wschannel"
)

type wsCaller struct {
	host string
	rc   *jrpc2.Client
}

// websocket call
func (w *wsCaller) call(method string, params, reply any) error {
	ch, err := wschannel.Dial(w.host, nil)
	if err != nil {
		return err
	}

	rc := jrpc2.NewClient(ch, nil)

	w.rc = rc

	if reply == nil {
		_, err := rc.Call(context.Background(), method, params)
		return err
	}

	rc.CallResult(context.Background(), method, params, reply)
	return nil
}

func (w *wsCaller) close() error {
	return w.rc.Close()
}
