# Aria2 JSON-RPC Client

This is a Go client for the Aria2 JSON-RPC interface, providing a way to interact with the Aria2 download manager programmatically.

## Installation

To install the package, run:

```bash
go get github.com/kahosan/aria2-rpc
```

## Usage

To use the client, first, import the package:

```go
import ario "github.com/kahosan/aria2-rpc"
```

Then create a new client with the `NewClient` function, passing in the host and token for the Aria2 instance:

```go
client, err := ario.NewClient("http://localhost:6800/jsonrpc", "mytoken")
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

Note that the methods take different parameters depending on the specific method being called. Refer to the [Aria2 documentation](https://aria2.github.io/manual/en/html/aria2c.html#methods) for details on each method.

## License

This library is licensed under the MIT License. See the LICENSE file for details.
