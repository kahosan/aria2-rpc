package notifier

import (
	"context"
	"io"
	"log"
	"net/url"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
)

type reply struct {
	Method string  `json:"method"`
	Params []Event `json:"params"`
}

type Event struct {
	Gid string `json:"gid"`
}

type Tasks = map[string]func(gid string)

type ne struct {
	Start      string
	Pause      string
	Stop       string
	Complete   string
	Error      string
	BtComplete string
}

type notifier struct {
	host *url.URL
}

type Notify struct {
	r     *sync.Map
	Close func()
}

var NotifyEvents = &ne{
	Start:      "aria2.onDownloadStart",
	Pause:      "aria2.onDownloadPause",
	Stop:       "aria2.onDownloadStop",
	Complete:   "aria2.onDownloadComplete",
	Error:      "aria2.onDownloadError",
	BtComplete: "aria2.onBtDownloadComplete",
}

func NewNotifier(host *url.URL) *notifier {
	return &notifier{
		host,
	}
}

func (n *notifier) Listener(c context.Context) (*Notify, error) {
	switch n.host.Scheme {
	case "https", "wss":
		n.host.Scheme = "wss"
	case "http", "ws":
		n.host.Scheme = "ws"
	}

	conn, _, err := websocket.DefaultDialer.Dial(n.host.String(), nil)
	if err != nil {
		return nil, err
	}

	r := sync.Map{}
	ctx, cancel := context.WithCancel(c)

	// create channels for each method, and store them in the map
	values := reflect.ValueOf(*NotifyEvents)
	for i := 0; i < values.NumField(); i++ {
		r.Store(values.Field(i).String(), make(chan string, 10))
	}

	go func() {
		defer func() {
			r.Range(func(key, value any) bool {
				close(value.(chan string))
				return true
			})
			conn.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			resp := &reply{}
			// read notifications from the connection
			if err := conn.ReadJSON(resp); err != nil {
				if err == io.ErrUnexpectedEOF {
					log.Println("unexpected EOF | if you are using nginx, please adjust the value of `proxy_read_timeout`")
					return
				}
				log.Printf("reading websocket message: %v", err)
				return
			}

			for _, event := range resp.Params {
				// created channels for all methods in advance
				ch, _ := r.Load(resp.Method)
				select {
				case ch.(chan string) <- event.Gid:
				default:
					// if the channel is full, skip the event, maybe the corresponding subscription does not exist
				}
			}
		}
	}()

	return &Notify{
		&r,
		cancel,
	}, nil
}

func (n *Notify) notifyFunc(method string) <-chan string {
	ch, _ := n.r.Load(method)
	return ch.(chan string)
}

func (n *Notify) ListenMultiple(events map[string]func(gid string)) {
	for method, fn := range events {
		go func(m string, f func(gid string)) {
			for gid := range n.notifyFunc(m) {
				f(gid)
			}
		}(method, fn)
	}
}

func (n *Notify) ListenOnce(method string, fn func(gid string, stop func())) {
	ch := n.notifyFunc(method)

	done := false
	stop := func() { done = true }

	for gid := range ch {
		if done {
			return
		}
		fn(gid, stop)
	}
}

func (n *Notify) Start() <-chan string {
	return n.notifyFunc(NotifyEvents.Start)
}

func (n *Notify) Pause() <-chan string {
	return n.notifyFunc(NotifyEvents.Pause)
}

func (n *Notify) Stop() <-chan string {
	return n.notifyFunc(NotifyEvents.Stop)
}

func (n *Notify) Complete() <-chan string {
	return n.notifyFunc(NotifyEvents.Complete)
}

func (n *Notify) Error() <-chan string {
	return n.notifyFunc(NotifyEvents.Error)
}

func (n *Notify) BtComplete() <-chan string {
	return n.notifyFunc(NotifyEvents.BtComplete)
}
