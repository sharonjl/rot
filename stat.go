package rot

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var (
	last               cpu.TimesStat
	rw                 = sync.RWMutex{}
	hasCPU             = false
	cpuUsage, memUsage = 0.0, 0.0
)

func updateStats(r time.Duration) {
	for range time.NewTicker(r).C {
		p, _ := cpu.Times(false)
		v, _ := mem.VirtualMemory()

		rw.Lock()
		if hasCPU {
			b := calculateBusy(last, p[0])
			memUsage = v.UsedPercent / 100
			if b > 0 {
				cpuUsage = b
				last = p[0]
			}
		} else {
			last = p[0]
			hasCPU = true
		}
		rw.Unlock()
	}
}

func limited(maxCPU, maxMem float64) bool {
	rw.RLock()
	defer rw.RUnlock()

	return memUsage > maxMem || cpuUsage > maxCPU || !hasCPU
}

func getAllBusy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal
	return busy + t.Idle, busy
}

func calculateBusy(t1, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return (t2Busy - t1Busy) / (t2All - t1All)
}
