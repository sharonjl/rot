# rot

rot is a go routine management library that limits excessive spawning 
of go routines by monitoring cpu and memory usage.

# example
```go
package main

import (
	"fmt"
	
	"github.com/sharonjl/rot"
)

func main() {
	// Go routines are launched only when cpu
	// and mem usages are under this limit.
	rot.SetLimits(0.8, 0.8)

	// Launch go routine, blocks until the routine is launched.
	rot.Go(func() {
		// do something
	})

	// Try launch go routine, returns false if it is not launched.
	didRun := rot.GoTry(func() {
		// do something
	})
	fmt.Println("didRun =", didRun)
}

```