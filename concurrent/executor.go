package concurrent

import (
	"time"
)

// RunPeriodically invokes function "f" every "interval" until something
// is sent to the "stop" channel. If the "stop" channel  is null, then the
// invocation never stops.
func RunPeriodically(f func(), stop chan bool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			f()

			if stop == nil {
				<-ticker.C
				continue
			}

			select {
			case <-ticker.C:
				continue
			case <-stop:
				break
			}
		}
		ticker.Stop()
	}()
}
