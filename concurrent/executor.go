package concurrent

import (
	"fmt"
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
	}()
}

type Runnable interface {
	Run()
}

type runnableImpl struct {
	Runnable

	f func()
}

func (ri *runnableImpl) Run() {
	ri.f()
}

func Run(f func()) Runnable {
	return &runnableImpl{
		f: f,
	}
}

type Executor interface {
	Execute(task Runnable)
	Shutdown()
}

type poolExecutor struct {
	Executor

	tasksToRun  chan Runnable
	stopSignal  chan bool
	workerCount int
}

func (pe *poolExecutor) Execute(task Runnable) {
	pe.tasksToRun <- task
}

func (pe *poolExecutor) Shutdown() {
	for range pe.workerCount {
		pe.stopSignal <- true
	}
	close(pe.tasksToRun)
	close(pe.stopSignal)
}

func (pe *poolExecutor) initWorkers() {
	for range pe.workerCount {
		go func() {
		loop:
			for {
				select {
				case task := <-pe.tasksToRun:
					task.Run()
				case <-pe.stopSignal:
					break loop
				}
			}
		}()
	}
}

func NewPoolExecutor(workerCount int) (Executor, error) {
	if workerCount <= 0 {
		return nil, fmt.Errorf("worker count should be greater than zero, but got %d", workerCount)
	}
	ex := &poolExecutor{
		workerCount: workerCount,
		tasksToRun:  make(chan Runnable, workerCount),
		stopSignal:  make(chan bool, workerCount),
	}
	ex.initWorkers()
	return ex, nil
}
