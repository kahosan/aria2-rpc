# Aria2 JSON-RPC Client

This is a Go client for the Aria2 JSON-RPC interface, providing a way to interact with the Aria2 download manager programmatically.

> **Note:** This library is still in development and is not yet ready for production use.

## Installation

To install the package, run:

```bash
go get -u github.com/kahosan/aria2-rpc
```

## Usage

To use the client, first, import the package:

```go
import ario "github.com/kahosan/aria2-rpc"
```

Then create a new client with the `NewClient` function, passing in the host and token for the Aria2 instance:

```go
client, err := ario.NewClient("http://localhost:6800/jsonrpc", "token", false)
// client, err := ario.NewClient("ws://localhost:6800/jsonrpc", "token", false)
if err != nil {
    // handle error
}
defer client.Close()
```

Once you have a client, you can use it to call any of the Aria2 methods:

```go
// download a file
gid, err := client.AddURI([]string{"http://example.com/file.txt"}, nil)
if err != nil {
    // handle error
}

// if you want to add an option 
opts := ario.Options{}
opts.Dir = "/path/to/dir"
gid, err := client.AddURI([]string{"http://example.com/file.txt"}, &opts)

// or
gid, err := client.AddURI([]string{"http://example.com/file.txt"}, nil)
if err != nil {
    // handle error
}

opts := ario.Options{}
opts.Dir = "/path/to/dir"
err := client.changeOption(gid, &opts)

// get status
status, err := client.TellStatus(gid, "key1", "key2")
if err != nil {
    // handle error
}
fmt.Println(status)
```

### Listener

You can also use the client to listen for events from Aria2. To do so, use the `NotifyListener` method to get an instance that has some events and return the channel with a value of gid.

Please refer to this [document](https://aria2.github.io/manual/en/html/aria2c.html#notifications) for the supported notification events.

```go
// both HTTP and WebSocket protocols can be used, but WebSocket protocol connection is required
// please refer to the events related to support in `notifier.go` file
client, err := ario.NewClient("http://localhost:6800/jsonrpc", "token", true)
if err != nil {
    // handle error
}

ctx, stopListener := context.WithCancel(context.Background())
notify, err := client.NotifyListener(ctx)
if err != nil {
    fmt.Println(err)
    return
}
defer stopListener()

gid, _ := client.AddURI([]string{"http://example.com/file.txt"}, nil)

// return a channel whose value is GID 
<-notify.Complete()

// or
for g := range notify.Complete() {
    if g == gid {
        // do something
    }
}
```

Note that the methods take different parameters depending on the specific method being called. Refer to the [Aria2 documentation](https://aria2.github.io/manual/en/html/aria2c.html#methods) for details on each method.

## License

This library is licensed under the MIT License. See the LICENSE file for details.
