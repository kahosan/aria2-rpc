<h1 align="center">Ario</h1>

This is a Go client for the Aria2 JSON-RPC interface

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

Note that the methods take different parameters depending on the specific method being called. Refer to the [Aria2 documentation](https://aria2.github.io/manual/en/html/aria2c.html#methods) for details on each method

### Listener

You can also use the client to listen for events from Aria2. To do so, use the `NotifyListener` method to get an instance that has some events and return the channel with a value of gid

Please refer to this [document](https://aria2.github.io/manual/en/html/aria2c.html#notifications) for the supported notification events

```go
// both HTTP and WebSocket protocols can be used, but WebSocket protocol connection is required
// please refer to the events related to support in `notifier.go` file
client, err := ario.NewClient("http://localhost:6800/jsonrpc", "token", true)
if err != nil {
    // handle error
}

notify, err := client.NotifyListener(context.Background())
if err != nil {
    //
}
defer notify.Close()

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

Using the callback method:

```go
notify, err := client.NotifyListener(context.Background())
if err != nil {
    //
}
defer notify.Close()

// blocking
notify.ListenOnce(notifier.NotifyEvents.Complete, func(g string, stop func()) {
    fmt.Println("Stop: ", g)
    stop()
})

// subscribe to multiple at the same time
wg := sync.WaitGroup{}
wg.Add(1) // Any task is completed, the listener is closed

go notify.ListenOnce(notifier.NotifyEvents.Start, func(g string, stop func()) {
    fmt.Println("Start: ", g)
    wg.Done()
})

go notify.ListenOnce(notifier.NotifyEvents.Pause, func(g string, stop func()) {
    fmt.Println("Pause: ", g)
    wg.Done()
})

go notify.ListenOnce(notifier.NotifyEvents.Stop, func(g string, stop func()) {
    fmt.Println("Stop: ", g)
    wg.Done()
})

wg.Wait()

// if there is still blocked code, it is recommended to execute the `notify.Close` function first
```

If you want to listen to multiple events at the same time:

```go
tasks := notifier.Tasks{
    notifier.NotifyEvents.Start: func(gid string) {
        fmt.Println("aria2.onDownloadStart", gid)
    },
    notifier.NotifyEvents.Pause: func(gid string) {
        fmt.Println("aria2.onDownloadPause", gid)
    },
    notifier.NotifyEvents.Stop: func(gid string) {
        fmt.Println("aria2.onDownloadStop", gid)
    },
}

// non-blocking
notify.ListenMultiple(tasks)
```

## License

This library is licensed under the MIT License. See the LICENSE file for details.
