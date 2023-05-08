package ario

import (
	"encoding/base64"
	"reflect"

	"github.com/kahosan/aria2-rpc/internal/caller"
	"github.com/kahosan/aria2-rpc/internal/resp"
)

type Client struct {
	call  func(method string, params, reply any) error
	close func() error
	token string
}

func NewClient(host string, token string) (*Client, error) {
	c, err := caller.NewCaller(host)
	if err != nil {
		return nil, err
	}

	return &Client{
		call:  c.Call,
		close: c.Close,
		token: token,
	}, nil
}

func (c *Client) Close() error {
	return c.close()
}

func (c *Client) AddURI(uris []string, options *Options) (gid string, err error) {
	c.call(method.AddURI, c.makeParams(uris, options), &gid)
	return
}

func (c *Client) AddTorrent(torrent *[]byte, uris *[]string, options *Options) (gid string, err error) {
	et := base64.StdEncoding.EncodeToString(*torrent)
	err = c.call(method.AddTorrent, c.makeParams(et, uris, options), &gid)
	return
}

func (c *Client) AddMetalink(metalink *[]byte, options *Options) (gid []string, err error) {
	em := base64.StdEncoding.EncodeToString(*metalink)
	err = c.call(method.AddMetalink, c.makeParams(em, options), &gid)
	return
}

func (c *Client) Remove(gid string) error {
	return c.call(method.Remove, c.makeParams(gid), nil)
}

func (c *Client) ForceRemove(gid string) error {
	return c.call(method.ForceRemove, c.makeParams(gid), nil)
}

func (c *Client) Pause(gid string) error {
	return c.call(method.Pause, c.makeParams(gid), nil)
}

func (c *Client) PauseAll() error {
	return c.call(method.PauseAll, c.makeParams(), nil)
}

func (c *Client) ForcePause(gid string) error {
	return c.call(method.ForcePause, c.makeParams(), nil)
}

func (c *Client) ForcePauseAll() error {
	return c.call(method.ForcePauseAll, c.makeParams(), nil)
}

func (c *Client) Unpause(gid string) error {
	return c.call(method.Unpause, c.makeParams(gid), nil)
}

func (c *Client) UnpauseAll() error {
	return c.call(method.UnpauseAll, c.makeParams(), nil)
}

func (c *Client) TellStatus(gid string, keys ...string) (status resp.Status, err error) {
	err = c.call(method.TellStatus, c.makeParams(gid, keys), &status)
	return
}

func (c *Client) GetURIs(gid string) (uris []resp.URIs, err error) {
	err = c.call(method.GetURIs, c.makeParams(gid), &uris)
	return
}

func (c *Client) GetFiles(gid string) (files []resp.Files, err error) {
	err = c.call(method.GetFiles, c.makeParams(gid), &files)
	return
}

func (c *Client) GetPeers(gid string) (peers []resp.Peers, err error) {
	err = c.call(method.GetPeers, c.makeParams(gid), &peers)
	return
}

func (c *Client) GetServers(gid string) (servers []resp.Servers, err error) {
	err = c.call(method.GetServers, c.makeParams(gid), &servers)
	return
}

func (c *Client) TellActive(keys ...string) (active []resp.Status, err error) {
	err = c.call(method.TellActive, c.makeParams(keys), &active)
	return
}

func (c *Client) TellWaiting(offset, num int, keys ...string) (waiting []resp.Status, err error) {
	err = c.call(method.TellWaiting, c.makeParams(offset, num, keys), &waiting)
	return
}

func (c *Client) TellStopped(offset, num int, keys ...string) (stopped []resp.Status, err error) {
	err = c.call(method.TellStopped, c.makeParams(offset, num, keys), &stopped)
	return
}

func (c *Client) ChangePosition(gid string, pos int, how string) (err error) {
	err = c.call(method.ChangePosition, c.makeParams(gid, pos, how), nil)
	return
}

func (c *Client) ChangeURI(gid string, fileIndex int, delURIs, addURIs *[]string, position int) (err error) {
	err = c.call(method.ChangeURI, c.makeParams(gid, fileIndex, delURIs, addURIs, position), nil)
	return
}

func (c *Client) GetOption(gid string) (options Options, err error) {
	err = c.call(method.GetOption, c.makeParams(gid), &options)
	return
}

func (c *Client) ChangeOption(gid string, options *Options) (err error) {
	err = c.call(method.ChangeOption, c.makeParams(gid, options), nil)
	return
}

func (c *Client) GetGlobalOption() (options Options, err error) {
	err = c.call(method.GetGlobalOption, c.makeParams(), &options)
	return
}

func (c *Client) ChangeGlobalOption(options *Options) (err error) {
	err = c.call(method.ChangeGlobalOption, c.makeParams(options), nil)
	return
}

func (c *Client) GetGlobalStat() (stat resp.GlobalStat, err error) {
	err = c.call(method.GetGlobalStat, c.makeParams(), &stat)
	return
}

func (c *Client) PurgeDownloadResult() error {
	return c.call(method.PurgeDownloadResult, c.makeParams(), nil)
}

func (c *Client) RemoveDownloadResult(gid string) error {
	return c.call(method.RemoveDownloadResult, c.makeParams(gid), nil)
}

func (c *Client) GetVersion() (version resp.Version, err error) {
	err = c.call(method.GetVersion, c.makeParams(), &version)
	return
}

func (c *Client) GetSessionInfo() (session resp.SessionInfo, err error) {
	err = c.call(method.GetSessionInfo, c.makeParams(), &session)
	return
}

func (c *Client) Shutdown() error {
	return c.call(method.Shutdown, c.makeParams(), nil)
}

func (c *Client) ForceShutdown() error {
	return c.call(method.ForceShutdown, c.makeParams(), nil)
}

func (c *Client) SaveSession() error {
	return c.call(method.SaveSession, c.makeParams(), nil)
}

// Method is an element of parameters used in system.multicall
type multiCallMethod struct {
	Name   string `json:"methodName"` // Method name to call
	Params []any  `json:"params"`     // Array containing parameters to the method call
}

func (c *Client) MultiCall(methods *[]multiCallMethod) (result []any, err error) {
	c.call(method.Multicall, c.makeParams(methods), &result)
	return
}

func (c *Client) ListMethods() (methods []string, err error) {
	err = c.call(method.ListMethods, c.makeParams(), &methods)
	return
}

func (c *Client) makeParams(p ...any) []any {
	params := make([]any, 0, 3)
	if c.token != "" {
		params = append(params, "token:"+c.token)
	}

	for _, u := range p {
		switch v := u.(type) {
		case []string, []byte:
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
