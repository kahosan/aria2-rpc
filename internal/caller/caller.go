package caller

import (
	"fmt"
	"net/url"
)

type Caller struct {
	Call  func(method string, params, reply any) error
	Close func() error
}

func NewCaller(host *url.URL) (*Caller, error) {
	rpc := &Caller{}

	switch host.Scheme {
	case "http", "https":
		h, err := newHttpCaller(host.String())
		if err != nil {
			return nil, err
		}

		rpc.Call = h.call
		rpc.Close = h.close
	case "ws", "wss":
		w, err := newWsCaller(host.String())
		if err != nil {
			return nil, err
		}

		rpc.Call = w.call
		rpc.Close = w.close
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", host.Scheme)
	}
	return rpc, nil
}
