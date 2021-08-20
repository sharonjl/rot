package rot

import (
	"sync"
)

var (
	mu             sync.Mutex
	c              uint64
	maxC           uint64
	maxCPU, maxMem = 0.8, 0.8
)

func init() {
	go updateStats()
}

func inc() {
	mu.Lock()
	c++
	maxC++
	mu.Unlock()
}

func dec() {
	mu.Lock()
	c--
	mu.Unlock()
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
	return maxC
}

func limited(maxCPU, maxMem float64) bool {
	statRW.RLock()
	defer statRW.RUnlock()

	return memUsage > maxMem || cpuUsage > maxCPU || !hasCPU
}

// SetLimits sets cpu and memory limits against which
// the library tests before launching a go routine.
func SetLimits(cpu, mem float64) {
	statRW.RLock()
	defer statRW.RUnlock()

	maxCPU = cpu
	maxMem = mem
}

// Go launches a go routine if the limit conditions are satisfied.
// Returns true if a routine is launched successfully, otherwise false.
func Go(fn func()) bool {
	if !limited(maxCPU, maxMem) {
		inc()
		go func() {
			defer dec()
			fn()
		}()
		return true
	}
	return false
}

// MustGo blocks until a go routine is launched.
func MustGo(fn func()) {
	for !Go(fn) {

	}
}
