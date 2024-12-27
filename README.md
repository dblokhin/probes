# Probes Package

The `probes` package provides a simple solution for managing Kubernetes HTTP probes (startup, readiness, liveness).

Probes are served on the following endpoints:
- `/startup` - Startup probe.
- `/ready` - Readiness probe (by default is `false`).
- `/live` - Liveness probe (by default is `true`).

## Starting the Server

Use `RunServer` to start the probe server.

```go
package main

import (
	"log"
	"time"
	"github.com/dblokhin/probes"
)

func main() {
	go func() {
		if err := probes.RunServer("", 12087); err != nil {
			log.Fatalf("Failed to start probe server: %v", err)
		}
	}()

	// you app is ready
	probes.Ready()

	// your app is not ready
	probes.Unready()
}
```