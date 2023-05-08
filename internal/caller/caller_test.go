package caller

import (
	"fmt"
	"testing"

	"github.com/kahosan/aria2-rpc/internal/resp"
)

func TestHTTPRPC(t *testing.T) {
	c, err := NewCaller("http://localhost:6800/jsonrpc")
	if err != nil {
		fmt.Println(err)
		t.Errorf("NewCaller should not return error")
	}

	if c.Call == nil {
		t.Errorf("Call function should not be nil")
	}

	t.Run("call getVersion", func(t *testing.T) {
		r := resp.Version{}
		err = c.Call("aria2.getVersion", nil, &r)
		if err != nil {
			t.Errorf("getVersion failed: %v", err)
		}

		t.Log(r)
	})
}

func TestWSRPC(t *testing.T) {
	c, err := NewCaller("ws://localhost:6800/jsonrpc")
	if err != nil {
		t.Fatal(err)
	}

	if c.Call == nil {
		t.Fatal("Call function should not be nil")
	}

	t.Run("call getVersion", func(t *testing.T) {
		r := resp.Version{}
		err = c.Call("aria2.getVersion", nil, &r)
		if err != nil {
			t.Errorf("getVersion failed: %v", err)
		}

		t.Log(r)
	})
}
