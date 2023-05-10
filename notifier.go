package ario

import (
	"context"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type reply struct {
	Method string  `json:"method"`
	Params []Event `json:"params"`
}

type Event struct {
	Gid string `json:"gid"`
}

type notifier struct {
	conn *websocket.Conn
}

type Notify struct {
	r     chan *reply
	Close func() error
}

const (
	ods = "aria2.onDownloadStart"
	odp = "aria2.onDownloadPause"
	odt = "aria2.onDownloadStop"
	odc = "aria2.onDownloadComplete"
	ode = "aria2.onDownloadError"
	obc = "aria2.onBtDownloadComplete"
)

func (n *notifier) setNotifier(host *url.URL) error {
	switch host.Scheme {
	case "https", "wss":
		host.Scheme = "wss"
	case "http", "ws":
		host.Scheme = "ws"
	}

	conn, _, err := websocket.DefaultDialer.Dial(host.String(), nil)
	if err != nil {
		return err
	}

	n.conn = conn

	return nil
}

func (n *notifier) Listener() (*Notify, error) {
	r := make(chan *reply)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer close(r)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			// read notifications from the connection
			resp := &reply{}
			if err := n.conn.ReadJSON(resp); err != nil {
				log.Printf("reading websocket message: %v", err)
				return
			}
			r <- resp
		}
	}()

	return &Notify{
		r,
		func() error {
			cancel()
			return n.conn.Close()
		},
	}, nil
}

func (n *Notify) notifyFunc(method string) chan string {
	gid := make(chan string)

	go func() {
		defer close(gid)
		for v := range n.r {
			if v.Method == method {
				gid <- v.Params[0].Gid
			}
		}
	}()

	return gid
}

func (n *Notify) Start() chan string {
	return n.notifyFunc(ods)
}

func (n *Notify) Pause() chan string {
	return n.notifyFunc(odp)
}

func (n *Notify) Stop() chan string {
	return n.notifyFunc(odt)
}

func (n *Notify) Complete() chan string {
	return n.notifyFunc(odc)
}

func (n *Notify) Error() chan string {
	return n.notifyFunc(ode)
}

func (n *Notify) BtComplete() chan string {
	return n.notifyFunc(obc)
}
