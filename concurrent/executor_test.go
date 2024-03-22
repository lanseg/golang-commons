package concurrent

import (
    "sync/atomic"
	"testing"
	"time"
)

const (
	testPeriodicInterval = 100 * time.Millisecond
)

func TestRunPeriodically(t *testing.T) {

	t.Run("run and stop", func(t *testing.T) {
        c := make(chan bool)
        i := int64(10)
        RunPeriodically(func () {
           if atomic.AddInt64(&i, -1) == 0 {
               c <- true
           }
        }, c, testInterval)
    })
}
