package caller

import (
	"fmt"
	"testing"

	"github.com/kahosan/aria2-rpc/internal/resp"
	"github.com/kahosan/aria2-rpc/internal/testutils"
)

func TestHTTPRPC(t *testing.T) {
	c, err := NewCaller(testutils.Arai2Uri("https://"))
	if err != nil {
		fmt.Println(err)
		t.Fatal("NewCaller should not return error")
	}

	t.Run("connect should not be error", func(t *testing.T) {
		r := resp.Version{}
		err = c.Call("aria2.getVersion", nil, &r)
		if err != nil {
			t.Fatal("get version failed: ", err)
		}

		t.Log(r)
	})

	t.Run("when reply is nil, error should not be returned.", func(t *testing.T) {
		err = c.Call("aria2.getVersion", nil, nil)
		if err != nil {
			t.Fatal("get version failed: ", err)
		}
	})
}

func TestWSRPC(t *testing.T) {
	c, err := NewCaller(testutils.Arai2Uri("ws://"))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("connect should not be error", func(t *testing.T) {
		r := resp.Version{}
		err = c.Call("aria2.getVersion", nil, &r)
		if err != nil {
			t.Fatal("get version failed: ", err)
		}

		t.Log(r)
	})

	t.Run("when reply is nil, error should not be returned.", func(t *testing.T) {
		err = c.Call("aria2.getVersion", nil, nil)
		if err != nil {
			t.Fatal("get version failed: ", err)
		}
	})
}
