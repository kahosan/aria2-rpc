package ario_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	ario "github.com/kahosan/aria2-rpc"
	"github.com/kahosan/aria2-rpc/internal/testutils"
)

func TestClient(t *testing.T) {

	t.Run("unsupported scheme", func(t *testing.T) {
		_, err := ario.NewClient("localhost:6800/jsonrpc", "", false)
		if err.Error() != "unsupported scheme: localhost" {
			t.Fatal("unexpected scheme check error")
		}
	})

	// method test
	client, err := ario.NewClient(testutils.Arai2Uri("https://"), "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	t.Run("add uri", func(t *testing.T) {
		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, nil)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(gid)

		time.Sleep(time.Second)
		err = client.Remove(gid)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("add uri with options", func(t *testing.T) {
		op := ario.Options{}
		op.Dir = "/tmp"

		// the returned value does not include storage path information
		// manual verification is required to confirm if it has been downloaded to the specified path.
		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, &op)
		if err != nil {
			t.Fatal(err)
		}

		// err = client.ChangeOption(gid, &op)
		// if err != nil {
		// 	t.Fatal(err)
		// }

		t.Log(gid)

		time.Sleep(time.Second)
		err = client.Remove(gid)
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("gid status listener", func(t *testing.T) {
		op := ario.Options{}
		op.Dir = "/tmp"
		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, &op)
		if err != nil {
			t.Fatal(err)
		}

		ctx, stopListen := context.WithCancel(context.Background())
		defer stopListen()

		statusChan := client.StatusListenerByPolling(ctx, gid)
		for v := range statusChan {
			switch v.Status {
			case "active":
				t.Log("task active")
				pe := client.Pause(gid)
				if pe != nil {
					t.Fatal(pe)
				}
			case "waiting":
				t.Log("task waiting")
			case "paused":
				t.Log("task paused")

				// if you directly delete a task while it is in pause state, aria2 will complain
				ue := client.Unpause(gid)
				if ue != nil {
					t.Fatal(ue)
				}

				re := client.Remove(gid)
				if re != nil {
					t.Fatal(re)
				}

			case "error":
				t.Log("task error")
				return
			case "complete":
				t.Log("task complete")
				return
			case "removed":
				t.Log("task removed")
				return
			}
		}
	})

	t.Run("get uris", func(t *testing.T) {
		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, nil)
		if err != nil {
			t.Fatal(err)
		}

		uris, err := client.GetURIs(gid)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(uris)

		time.Sleep(time.Second)
		err = client.Remove(gid)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get files", func(t *testing.T) {
		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, nil)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(gid)

		files, err := client.GetFiles(gid)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(files)

		time.Sleep(time.Second)
		err = client.Remove(gid)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get version", func(t *testing.T) {
		version, err := client.GetVersion()
		if err != nil {
			t.Fatal(err)
		}

		t.Log(version)
	})

	t.Run("get methods", func(t *testing.T) {
		methods, err := client.ListMethods()
		if err != nil {
			t.Fatal(err)
		}

		t.Log(methods)
	})
}

func TestMultiCall(t *testing.T) {
	client, err := ario.NewClient(testutils.Arai2Uri("https://"), "", false)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("parameter for empty", func(t *testing.T) {
		_, err := client.MultiCall(nil)
		if err == nil {
			t.Fatal("unexpected parameter check error")
		}
	})

	t.Run("call multi add uri", func(t *testing.T) {
		t.Run("multicall", func(t *testing.T) {
			mc1 := ario.MultiCallMethod{
				Name:   "aria2.addUri",
				Params: []any{[]any{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}},
			}
			mc2 := ario.MultiCallMethod{
				Name:   "aria2.addUri",
				Params: []any{[]any{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}},
			}

			result, err := client.MultiCall(&[]ario.MultiCallMethod{mc1, mc2})
			if err != nil {
				t.Fatal(err)
			}

			for _, v := range result {
				t.Log(v.([]any)[0].(string)) // gid

				time.Sleep(time.Second)

				// different methods will return different data structures. For details, please refer to the aria2 rpc documentation.
				// https://aria2.github.io/manual/en/html/aria2c.html#methods

				if reflect.ValueOf(v).Kind() == reflect.Slice {
					err = client.Remove(v.([]any)[0].(string))
					if err != nil {
						t.Fatal(err)
					}
				}
			}
		})
	})
}
