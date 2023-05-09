package testutils

import (
	"net/url"
	"os"
)

func Arai2Uri(scheme string) string {
	if uri := os.Getenv("ARIA2_URI"); uri != "" {
		uri, _ := url.Parse(uri)
		return scheme + uri.Host + uri.Path
	} else {
		return "http://localhost:6800/jsonrpc"
	}
}
