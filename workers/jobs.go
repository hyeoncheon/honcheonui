package workers

import (
	"errors"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/worker"
)

// constants
const (
	DefaultQueue = "default"
)

// HandlerHolder is an interface for workers.
type HandlerHolder interface {
	Reset() error
	Handler(worker.Args) error
}

// Worker is a struct for managing workers.
type Worker struct {
	HandlerHolder
	Name         string
	IsPeriodic   bool
	InitailDelay time.Duration
	RunPeriod    time.Duration
	LastQueuedAt time.Time
	CountQueued  int32
}

// Workers is a search map for the workers with its name.
type Workers map[string]*Worker

var workers = Workers{}
var logger buffalo.Logger
var aw worker.Worker

//*** initiators

func init() {
	RegisterWorkers(&Worker{
		HandlerHolder: &repeater{},
		Name:          workerRepeater,
	})
}

// InitWorkers registers all application workers after system workers.
func InitWorkers(app *buffalo.App) error {
	logger = app.Logger.WithField("category", "worker")
	aw = app.Worker
	logger.Infof("register workers...")

	for name, wkr := range workers {
		logger.Debugf("---> worker: %v %v", name, wkr)
		if err := wkr.Reset(); err != nil {
			logger.Errorf("could not initialize worker %v: %v", name, err)
			continue
		}
		if err := aw.Register(name, wkr.Handler); err != nil {
			logger.Errorf("could not register worker %v: %v", name, err)
			continue
		}
		if wkr.IsPeriodic && wkr.RunPeriod >= (5*time.Second) {
			logger.Infof("%v is periodic worker. initial queueing...", name)
			err := Queue(
				workerRepeater,
				map[string]interface{}{
					"worker": wkr.Name,
					"args":   worker.Args{},
					"repeat": wkr.RunPeriod,
				},
				wkr.InitailDelay,
			)
			if err != nil {
				logger.Errorf("oops! could not add a queue for %v: %v", name, err)
			}
		}
	}
	return nil
}

//*** helper functions

// Run runs the worker immediately
func Run(name string, args worker.Args) error {
	return Queue(name, args, 0)
}

// Queue enqueues the worker after given delay
func Queue(name string, args worker.Args, delay time.Duration) error {
	w := workers[name]
	if w == nil {
		return errors.New("could not find worker")
	}
	w.CountQueued++
	w.LastQueuedAt = time.Now()

	return aw.PerformIn(worker.Job{
		Queue:   DefaultQueue,
		Handler: name,
		Args:    args,
	}, delay)
}

// RegisterWorkers adds given workers into workers map.
func RegisterWorkers(ws ...*Worker) error {
	for _, w := range ws {
		if w.Name != "" {
			workers[w.Name] = w
		}
	}
	return nil
}

//*** system worker: repeater ------
const workerRepeater = "worker.Repeater"

type repeater struct{}

// repeaterHandler is simple system worker handles periodic jobs.
func (r repeater) Handler(wa worker.Args) error {
	logger.Debugf("------ cron handler invoked with %v", wa)
	name, ok := wa["worker"].(string)
	if !ok {
		logger.Errorf("could not cast worker name: %v", wa["worker"])
		return errors.New("could not cast worker name")
	}
	repeat, ok := wa["repeat"].(time.Duration)
	if !ok {
		logger.Errorf("could not cast worker repeat: %v", wa["repeat"])
		return errors.New("could not cast worker repeat")
	}
	args := wa["args"].(worker.Args)
	logger.Debugf("------ run %v and queue new instance...", name)
	Run(name, args)
	return Queue(workerRepeater, wa, repeat)
}

func (r repeater) Reset() error {
	return nil
}
