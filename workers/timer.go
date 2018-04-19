// This is just a sample code for background job worker

package workers

import (
	"time"

	"github.com/gobuffalo/buffalo/worker"
)

// constants belongs to this worker
const (
	WorkerTimer             = "worker.Timer"
	workerTimerInitailDelay = 60 * time.Second
	workerTimerRunPeriod    = 30 * time.Second
)

// Timer is
type Timer struct{}

func init() {
	RegisterWorkers(&Worker{
		HandlerHolder: &Timer{},
		Name:          WorkerTimer,
		IsPeriodic:    true,
		InitailDelay:  workerTimerInitailDelay,
		RunPeriod:     workerTimerRunPeriod,
	})
}

// Handler implements HandlerHolder
func (j Timer) Handler(args worker.Args) error {
	logger.Debugf("-----> %v invoked with %v", WorkerTimer, args)
	return nil
}

// Reset implements HandlerHolder
func (j Timer) Reset() error {
	logger.Debugf("-----> reset %v...", WorkerTimer)
	return nil
}
