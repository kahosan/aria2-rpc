package ario

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"time"

	"github.com/kahosan/aria2-rpc/internal/caller"
	"github.com/kahosan/aria2-rpc/internal/resp"
	"github.com/kahosan/aria2-rpc/notifier"
)

type Client struct {
	Call           func(method string, params, reply any) error
	Close          func() error
	token          string
	NotifyListener func(ctx context.Context) (*notifier.Notify, error)
}

func NewClient(host string, token string, notify bool) (*Client, error) {
	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	c, err := caller.NewCaller(uri)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Call:  c.Call,
		Close: c.Close,
		token: token,
		NotifyListener: func(context.Context) (*notifier.Notify, error) {
			return nil, fmt.Errorf("please set the notify parameter to true in the NewClient function")
		},
	}

	if notify {
		not := notifier.NewNotifier(uri)
		client.NotifyListener = not.Listener
	}

	return client, nil
}

// only use when websocket is not supported, or if you want to use it yourself.
// instructions for use -> https://github.com/kahosan/aria2-rpc/blob/master/client_test.go#L68
func (c *Client) StatusListenerByPolling(ctx context.Context, gid string) (status chan *resp.Status) {
	status = make(chan *resp.Status)

	go func() {
		defer close(status)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				s, e := c.TellStatus(gid)
				if e != nil {
					log.Printf("listener error: %v", e)
					return
				} else if s.Gid == "" {
					log.Println("gid not found, maybe it was removed")
					return
				}

				status <- &s
			}

			time.Sleep(time.Second)
		}
	}()

	return
}

func (c *Client) AddURI(uris []string, options *Options) (gid string, err error) {
	err = c.Call(method.AddURI, c.makeParams(uris, options), &gid)
	return
}

func (c *Client) AddTorrent(torrent *[]byte, uris *[]string, options *Options) (gid string, err error) {
	et := base64.StdEncoding.EncodeToString(*torrent)
	err = c.Call(method.AddTorrent, c.makeParams(et, uris, options), &gid)
	return
}

func (c *Client) AddMetalink(metalink *[]byte, options *Options) (gid []string, err error) {
	em := base64.StdEncoding.EncodeToString(*metalink)
	err = c.Call(method.AddMetalink, c.makeParams(em, options), &gid)
	return
}

func (c *Client) Remove(gid string) error {
	return c.Call(method.Remove, c.makeParams(gid), nil)
}

func (c *Client) ForceRemove(gid string) error {
	return c.Call(method.ForceRemove, c.makeParams(gid), nil)
}

func (c *Client) Pause(gid string) error {
	return c.Call(method.Pause, c.makeParams(gid), nil)
}

func (c *Client) PauseAll() error {
	return c.Call(method.PauseAll, c.makeParams(), nil)
}

func (c *Client) ForcePause(gid string) error {
	return c.Call(method.ForcePause, c.makeParams(), nil)
}

func (c *Client) ForcePauseAll() error {
	return c.Call(method.ForcePauseAll, c.makeParams(), nil)
}

func (c *Client) Unpause(gid string) error {
	return c.Call(method.Unpause, c.makeParams(gid), nil)
}

func (c *Client) UnpauseAll() error {
	return c.Call(method.UnpauseAll, c.makeParams(), nil)
}

func (c *Client) TellStatus(gid string, keys ...string) (status resp.Status, err error) {
	err = c.Call(method.TellStatus, c.makeParams(gid, keys), &status)
	return
}

func (c *Client) GetURIs(gid string) (uris []resp.URIs, err error) {
	err = c.Call(method.GetURIs, c.makeParams(gid), &uris)
	return
}

func (c *Client) GetFiles(gid string) (files []resp.Files, err error) {
	err = c.Call(method.GetFiles, c.makeParams(gid), &files)
	return
}

func (c *Client) GetPeers(gid string) (peers []resp.Peers, err error) {
	err = c.Call(method.GetPeers, c.makeParams(gid), &peers)
	return
}

func (c *Client) GetServers(gid string) (servers []resp.Servers, err error) {
	err = c.Call(method.GetServers, c.makeParams(gid), &servers)
	return
}

func (c *Client) TellActive(keys ...string) (active []resp.Status, err error) {
	err = c.Call(method.TellActive, c.makeParams(keys), &active)
	return
}

func (c *Client) TellWaiting(offset, num int, keys ...string) (waiting []resp.Status, err error) {
	err = c.Call(method.TellWaiting, c.makeParams(offset, num, keys), &waiting)
	return
}

func (c *Client) TellStopped(offset, num int, keys ...string) (stopped []resp.Status, err error) {
	err = c.Call(method.TellStopped, c.makeParams(offset, num, keys), &stopped)
	return
}

func (c *Client) ChangePosition(gid string, pos int, how string) (err error) {
	err = c.Call(method.ChangePosition, c.makeParams(gid, pos, how), nil)
	return
}

func (c *Client) ChangeURI(gid string, fileIndex int, delURIs, addURIs *[]string, position ...int) (err error) {
	err = c.Call(method.ChangeURI, c.makeParams(gid, fileIndex, delURIs, addURIs, position), nil)
	return
}

func (c *Client) GetOption(gid string) (options Options, err error) {
	err = c.Call(method.GetOption, c.makeParams(gid), &options)
	return
}

func (c *Client) ChangeOption(gid string, options *Options) (err error) {
	err = c.Call(method.ChangeOption, c.makeParams(gid, options), nil)
	return
}

func (c *Client) GetGlobalOption() (options Options, err error) {
	err = c.Call(method.GetGlobalOption, c.makeParams(), &options)
	return
}

func (c *Client) ChangeGlobalOption(options *Options) (err error) {
	err = c.Call(method.ChangeGlobalOption, c.makeParams(options), nil)
	return
}

func (c *Client) GetGlobalStat() (stat resp.GlobalStat, err error) {
	err = c.Call(method.GetGlobalStat, c.makeParams(), &stat)
	return
}

func (c *Client) PurgeDownloadResult() error {
	return c.Call(method.PurgeDownloadResult, c.makeParams(), nil)
}

func (c *Client) RemoveDownloadResult(gid string) error {
	return c.Call(method.RemoveDownloadResult, c.makeParams(gid), nil)
}

func (c *Client) GetVersion() (version resp.Version, err error) {
	err = c.Call(method.GetVersion, c.makeParams(), &version)
	return
}

func (c *Client) GetSessionInfo() (session resp.SessionInfo, err error) {
	err = c.Call(method.GetSessionInfo, c.makeParams(), &session)
	return
}

func (c *Client) Shutdown() error {
	return c.Call(method.Shutdown, c.makeParams(), nil)
}

func (c *Client) ForceShutdown() error {
	return c.Call(method.ForceShutdown, c.makeParams(), nil)
}

func (c *Client) SaveSession() error {
	return c.Call(method.SaveSession, c.makeParams(), nil)
}

// Method is an element of parameters used in system.multicall
type MultiCallMethod struct {
	Name   string `json:"methodName"` // Method name to call
	Params []any  `json:"params"`     // Array containing parameters to the method call
}

// if MultiCallMethod for empty method name is given, it will be ignored
func (c *Client) MultiCall(methods *[]MultiCallMethod) (result []any, err error) {
	if methods == nil {
		return nil, fmt.Errorf("invalid parameter")
	}

	err = c.Call(method.Multicall, c.makeParams(methods), &result)
	return
}

func (c *Client) ListMethods() (methods []string, err error) {
	err = c.Call(method.ListMethods, c.makeParams(), &methods)
	return
}

func (c *Client) makeParams(p ...any) []any {
	params := make([]any, 0, 3)
	if c.token != "" {
		params = append(params, "token:"+c.token)
	}

	for _, u := range p {
		switch v := u.(type) {
		case []string, []byte, []int:
			if reflect.ValueOf(v).Len() != 0 {
				params = append(params, v)
			}
		case *Options:
			if v != nil {
				params = append(params, v)
			}
		default:
			params = append(params, v)
		}
	}

	return params
}
