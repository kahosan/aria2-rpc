package caller

import (
	"fmt"
	"net/url"
)

type Caller struct {
	Call  func(method string, params, reply any) error
	Close func() error
}

func NewCaller(host string) (*Caller, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	rpc := &Caller{}

	switch u.Scheme {
	case "http", "https":
		h := &httpCaller{
			host: host,
		}

		rpc.Call = h.call
		rpc.Close = h.close
	case "ws", "wss":
		w := &wsCaller{
			host: host,
		}

		rpc.Call = w.call
		rpc.Close = w.close
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
	return rpc, nil
}
