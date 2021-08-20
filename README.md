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

	// Return false if rot was unable to launch the go routine,
	// because the cpu and mem usage has exceeded.
	ok := rot.GoTry(func() {
		// do something
	})
	fmt.Println("did launch =", ok)
}

```