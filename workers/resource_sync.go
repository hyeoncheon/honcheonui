package workers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/uuid"
	"github.com/jinzhu/copier"

	"github.com/hyeoncheon/honcheonui/models"
	"github.com/hyeoncheon/honcheonui/plugins"
)

//*** background worker implementation

// constants belongs to this worker
const (
	WorkerResourceSync             = "worker.ResourceSync"
	workerResourceSyncInitailDelay = 10 * time.Second
	workerResourceSyncRunPeriod    = 24 * time.Hour
)

// ResourceSync is worker to sync resources via plugin in batch mode.
type ResourceSync struct{}

func init() {
	RegisterWorkers(&Worker{
		HandlerHolder: &ResourceSync{},
		Name:          WorkerResourceSync,
		IsPeriodic:    true,
		InitailDelay:  workerResourceSyncInitailDelay,
		RunPeriod:     workerResourceSyncRunPeriod,
	})
}

// Handler implements HandlerHolder
func (j ResourceSync) Handler(args worker.Args) error {
	providerID := args["provider_id"]
	return syncResources(providerID)
}

// Reset implements HandlerHolder
func (j ResourceSync) Reset() error {
	return nil
}

//*** local task functions

func syncResources(id interface{}) error {
	providers := &models.Providers{}
	query := models.DB.Q()
	if id != nil {
		query = query.Where("id = ?", id)
	}
	if err := query.All(providers); err != nil {
		logger.Errorf("database error: %v", err)
		return err
	}
	logger.Infof("sync resources via plugins for %v", *providers)

	for _, provider := range *providers {
		logger.Debugf("sync resources for %v...", provider)

		plugin, err := plugins.GetPlugin(provider.Provider, "provider")
		if err != nil {
			logger.Errorf("could not find plugin %v: %v", provider.Provider, err)
			continue
		}
		resources, err := plugin.GetResources(provider.User, provider.Pass)
		if err != nil {
			logger.Errorf("could not get resources via plugin: %v", err)
		}
		logger.Debugf("got %v resources. create/update...", len(resources))

		var ids []uuid.UUID
		for _, r := range resources {
			res := &models.Resource{}
			if jr, err := json.Marshal(r); err == nil {
				logger.Debugf("------ found: %v", string(jr))
				re := &plugins.HoncheonuiResource{}
				if err := json.Unmarshal(jr, re); err != nil {
					logger.Errorf("error: %v", err)
					return errors.New("cloud not recognize data format")
				}
				copier.Copy(res, re)
				// universally unique identifier, uuid is not perfectly uniq but almost.
				// but we can assume it is uniq anyway.
				// buffalo/pop uses uuid version 4 based on random number generator and
				// softlayer seems to use real random string as uuid. :-(
				if res.UUID != uuid.Nil {
					res.ID = res.UUID
				}
				ids = append(ids, res.ID)
				if err := res.Save(); err != nil {
					logger.Errorf("saving error: %v", err)
					logger.Errorf("---- resource: %v", res.JSON())
				}
				for k, v := range re.Attributes {
					res.AddAttribute(k, v)
				}
				if err := res.LinkTags(re.Tags); err != nil {
					logger.Debugf("problem on mapping tags")
				}
				if err := res.LinkUsers(re.UserIDs); err != nil {
					logger.Debugf("problem on mapping users")
				}
				// TODO: support IntegerAttributes
				logger.Debugf("------ re.IntegerAttributes: %v", re.IntegerAttributes)
			}
		}
		if err := provider.LinkResources(ids); err != nil {
			logger.Debugf("problem on mapping provider")
		}

		logger.Debugf("resources for %v are synced successfully", provider)
	}
	return nil
}
