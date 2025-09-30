package helpFlag

import (
	"fmt"
)

func Help() {
	fmt.Println(`
Usage: go run . [flags]

Flags:
  -product string
        your product name (default "gopher")
  -short
        shorten image URLs in output
  -webSocketDebuggerUrl string
        your websocket url
  -help
        show help

Examples:
  go run . -product="iphone15"
  go run . -product="samsung" -short
  go run . -webSocketDebuggerUrl="ws://localhost:9222/devtools/browser/..."
	`)
}