package rot

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	done := make(chan int)
	hc := 0
	lc := 0
	<-time.NewTimer(time.Second * 2).C
	go func() {
		for {
			if rand.Float32() > 0.95 {
				Go(highUsage(done))
				hc++
			} else {
				Go(lowUsage)
				lc++
			}
		}
	}()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	ticker := time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("cpu: %0.2f mem: %0.2f routines active: %d max: %d HighUsage: %d LowUsage: %d\n ", cpuUsage, memUsage, Count(), Max(), hc, lc)
		case <-ctx.Done():
			close(done)
			return
		}
	}
}

func lowUsage() {
	c := rand.Uint64()
	for i := 0; i < rand.Intn(10000); i++ {
		c = (c * uint64(i) / c) + c
	}
}

func highUsage(done chan int) func() {
	return func() {
		deadline := time.NewTimer(time.Millisecond * time.Duration(rand.Intn(3000)))
		for {
			select {
			case <-deadline.C:
				return
			case <-done:
				return
			default:
			}
		}
	}
}
