package workers

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/buffalo/worker"
	"github.com/hyeoncheon/spec"
	"github.com/jinzhu/copier"

	"github.com/hyeoncheon/honcheonui/models"
	"github.com/hyeoncheon/honcheonui/plugins"
)

//*** background worker implementation

// constants belongs to this worker
const (
	WorkerNotificationWatch             = "worker.NotificationWatch"
	workerNotificationWatchInitailDelay = 30 * time.Second
	workerNotificationWatchRunPeriod    = 24 * time.Hour
)

// NotificationWatch is worker to sync resources via plugin in batch mode.
type NotificationWatch struct{}

func init() {
	RegisterWorkers(&Worker{
		HandlerHolder: &NotificationWatch{},
		Name:          WorkerNotificationWatch,
		IsPeriodic:    true,
		InitailDelay:  workerNotificationWatchInitailDelay,
		RunPeriod:     workerNotificationWatchRunPeriod,
		LastQueuedAt:  time.Time{},
		CountQueued:   0,
	})
}

// Handler implements HandlerHolder
func (j NotificationWatch) Handler(args worker.Args) error {
	providerID := args["provider_id"]
	return watchNotification(providerID)
}

// Reset implements HandlerHolder
func (j NotificationWatch) Reset() error {
	return nil
}

//*** local task functions

func watchNotification(id interface{}) error {
	providers := &models.Providers{}
	query := models.DB.Q()
	if id != nil {
		logger.Errorf("single mode")
		query = query.Where("id = ?", id)
	}
	if err := query.All(providers); err != nil {
		logger.Errorf("database error: %v", err)
		return err
	}
	logger.Infof("watch notifications via plugins for %v", *providers)

	// TODO: use brand account without credential loop
	for _, provider := range *providers {
		logger.Debugf("watch notifications for %v...", provider)

		plugin, err := plugins.GetPlugin(provider.Provider, "provider")
		if err != nil {
			logger.Errorf("could not find plugin %v: %v", provider.Provider, err)
			continue
		}
		since := time.Now().AddDate(0, -5, 0)
		notes, err := plugin.GetNotifications(provider.User, provider.Pass, since)
		if err != nil {
			logger.Errorf("could not get notifications via plugin: %v", err)
		}
		logger.Debugf("got %v notifications. create/update...", len(notes))

		//var ids []uuid.UUID
		for _, n := range notes {
			note, ok := n.(spec.HoncheonuiNotification)
			if !ok {
				logger.Errorf("unrecognized data format: %T", n)
			}
			if jb, err := json.Marshal(note); err == nil {
				logger.Debugf("------ found: %v", string(jb))
			}

			inci := &models.Incident{}
			if err := copier.Copy(inci, note); err != nil {
				logger.Errorf("object copying error for %v", note)
				continue
			}
			if err := inci.Save(); err != nil {
				logger.Errorf("could not save incident record: %v", err)
			}

			inci.LinkResourcesByOrigIDs(note.ResourceIDs...)
			inci.LinkUsers(note.UserIDs...)
			if jb, err := json.Marshal(inci); err == nil {
				logger.Debugf("------ note: %v", string(jb))
			}
		}

		logger.Debugf("resources for %v are synced successfully (%v)", provider, len(notes))
	}
	return nil
}
