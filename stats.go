package rot

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var statRW = sync.RWMutex{}
var (
	hasCPU     = false
	lastCPU    cpu.TimesStat
	cpuUsage   = 0.0
	memUsage   = 0.0
	updateRate = time.Millisecond
)

func SetUpdateRate(d time.Duration) {
	updateRate = d
}

func updateStats() {
	for range time.NewTicker(updateRate).C {
		p, _ := cpu.Times(false)
		v, _ := mem.VirtualMemory()

		statRW.Lock()
		if hasCPU {
			b := calculateBusy(lastCPU, p[0])
			memUsage = v.UsedPercent / 100
			if b > 0 {
				cpuUsage = b
				lastCPU = p[0]
			}
		} else {
			lastCPU = p[0]
			hasCPU = true
		}
		statRW.Unlock()
	}
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
