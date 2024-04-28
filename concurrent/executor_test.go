package concurrent

import (
	"fmt"
	"sync"
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
		RunPeriodically(func() {
			if atomic.AddInt64(&i, -1) == 0 {
				c <- true
			}
		}, c, testInterval)
	})
}

func TestExecutor(t *testing.T) {

	t.Run("zero worker count returns error", func(t *testing.T) {
		if ex, err := NewPoolExecutor(0); ex != nil || err == nil {
			t.Errorf("Expected executor to be nil and error non-nil, but got %v and %v", ex, err)
		}
	})

	t.Run("single executor start stop", func(t *testing.T) {
		ex, err := NewPoolExecutor(1)
		if err != nil {
			t.Fatalf("Unexpected error when creating an executor: %s", err)
		}
		stop := sync.Mutex{}
		stop.Lock()
		ex.Execute(Run(func() {
			stop.Unlock()
		}))
		stop.Lock()
	})

	t.Run("single executor start end shutdown", func(t *testing.T) {
		execs := 100
		ex, err := NewPoolExecutor(execs)
		if err != nil {
			t.Fatalf("Unexpected error when creating an executor: %s", err)
		}

		wg := sync.WaitGroup{}
		wg.Add(execs)
		for i := range execs {
			ex.Execute(Run(func() {
				fmt.Printf("Thread %d\n", i)
				wg.Done()
			}))
		}
		wg.Wait()
	})
}
