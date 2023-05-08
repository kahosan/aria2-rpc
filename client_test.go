package ario

import "testing"

func TestClient(t *testing.T) {

	t.Run("unsupported scheme", func(t *testing.T) {
		_, err := NewClient("localhost:6800/jsonrpc", "")
		if err.Error() != "unsupported scheme: localhost" {
			t.Fatal("unexpected scheme check error")
		}
	})

	t.Run("http or https", func(t *testing.T) {
		_, err := NewClient("http://localhost:6800/jsonrpc", "")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ws or wss", func(t *testing.T) {
		_, err := NewClient("ws://localhost:6800/jsonrpc", "")
		if err != nil {
			t.Fatal(err)
		}
	})

	url := "http://localhost:6800/jsonrpc"

	// method test
	client, err := NewClient(url, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("add uri", func(t *testing.T) {
		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, nil)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(gid)
	})

	t.Run("add uri with options", func(t *testing.T) {
		op := Options{}
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
