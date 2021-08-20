package rot

import (
	"runtime"
	"sync"
	"time"
)

var (
	mu                 sync.Mutex
	c                  uint64
	C                  uint64
	polling            bool
	rate               time.Duration = time.Millisecond * 100
	cpuLimit, memLimit               = 0.8, 0.8
)

func inc() {
	mu.Lock()
	c++
	C++
	mu.Unlock()
}

func dec() {
	mu.Lock()
	c--
	mu.Unlock()
}

func pollIfNotPolling() {
	if polling {
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	polling = true
	go updateStats(rate)
}

// SetPollRate set how often we want to poll for cpu and memory usage.
func SetPollRate(t time.Duration) {
	mu.Lock()
	rate = t
	mu.Unlock()
}

// SetLimits sets cpu and memory limits against which
// the library tests before launching a go routine.
func SetLimits(cpu, mem float64) {
	rw.Lock()
	cpuLimit = cpu
	memLimit = mem
	rw.Unlock()
}

// Count returns the active number of go routines.
func Count() uint64 {
	mu.Lock()
	defer mu.Unlock()
	return c
}

// Max returns the total number of go routines launched.
func Max() uint64 {
	mu.Lock()
	defer mu.Unlock()
	return C
}

// GoTry launches a go routine if the limit conditions are satisfied.
// Returns true if a routine is launched successfully, otherwise false.
func GoTry(fn func()) bool {
	pollIfNotPolling()
	if !limited(cpuLimit, memLimit) {
		inc()
		go func() {
			defer dec()
			fn()
		}()
		return true
	}
	return false
}

// Go blocks until a go routine is launched.
func Go(fn func()) {
	for !GoTry(fn) {

	}
}
