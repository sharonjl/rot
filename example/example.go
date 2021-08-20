package main

import (
	"fmt"

	"rot"
)

func main() {
	// Launch go routine
	rot.Go(func() {
		// do something
	})

	// Return false if rot was unable to launch the go routine.
	didRun := rot.GoTry(func() {
		// do something
	})
	fmt.Println("didRun =", didRun)
}
