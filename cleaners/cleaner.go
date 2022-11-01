package cleaners

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/axieinfinity/bridge-cleaner/configs"
	bridgeCoreStores "github.com/axieinfinity/bridge-core/stores"
	"github.com/axieinfinity/bridge-v2/stores"
	"github.com/axieinfinity/bridge-v2/task"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-co-op/gocron"
)

var (
	ErrNotEnoughRecord = errors.New("not enough records")
)

type CleanOptions struct {
	Threshold  int64 `json:"threshold"`  // Limitation number of record to start cleaning
	Expiration int64 `json:"expiration"` // Expiration in second
}

type Cleaner struct {
	scheduler *gocron.Scheduler
	store     stores.BridgeStore
	mainStore bridgeCoreStores.MainStore
}

func (c *Cleaner) ClearSuccessTasks(opts *CleanOptions) error {
	log.Info("Doing clear success tasks")
	count := c.store.GetTaskStore().Count()
	if count <= opts.Threshold {
		log.Info("Skipping clear success tasks")
		return ErrNotEnoughRecord
	}
	if err := c.store.GetTaskStore().DeleteTasks([]string{task.STATUS_DONE}, uint64(time.Now().Unix())-uint64(opts.Expiration)); err != nil {
		log.Error("[ClearSuccessTasks] error while deleting done tasks", "err", err)
		return err
	}

	log.Info("Cleared success tasks")
	return nil
}

func (c *Cleaner) ClearFailedTasks(opts *CleanOptions) error {
	log.Info("Doing clear failed tasks")
	count := c.store.GetTaskStore().Count()
	if count <= opts.Threshold {
		log.Info("Skipping clear failed tasks")
		return ErrNotEnoughRecord
	}

	if err := c.store.GetTaskStore().DeleteTasks([]string{task.STATUS_FAILED}, uint64(time.Now().Unix())-uint64(opts.Expiration)); err != nil {
		log.Error("[ClearFailedTasks] error while deleting failed tasks", "err", err)
	}
	log.Info("Cleared failed tasks")

	return nil
}

func (c *Cleaner) ClearEvents(opts *CleanOptions) error {
	log.Info("Doing clear events")
	count := c.mainStore.GetEventStore().Count()

	if count <= opts.Threshold {
		log.Info("Skipping clear events")
		return ErrNotEnoughRecord
	}
	if err := c.mainStore.GetEventStore().DeleteEvents(uint64(time.Now().Unix()) - uint64(opts.Expiration)); err != nil {
		log.Error("[ClearEvents] error while deleting failed tasks", "err", err)
		return err
	}

	log.Info("Cleared events")
	return nil
}

func (c *Cleaner) ClearSuccessJobs(opts *CleanOptions) error {
	log.Info("Doing clear success jobs")
	count := c.mainStore.GetJobStore().Count()
	if count <= opts.Threshold {
		log.Info("Skipping clear success jobs")
		return ErrNotEnoughRecord
	}
	if err := c.mainStore.GetJobStore().DeleteJobs([]string{task.STATUS_DONE}, uint64(time.Now().Unix())-uint64(opts.Expiration)); err != nil {
		log.Error("[ClearSuccessJobs] error while deleting done tasks", "err", err)
		return nil
	}
	log.Info("Cleared success jobs")

	return nil
}

func (c *Cleaner) ClearFailedJobs(opts *CleanOptions) error {
	log.Info("Doing clear failed jobs")
	count := c.mainStore.GetJobStore().Count()
	if count <= opts.Threshold {
		log.Info("Skipping clear failed jobs")
		return ErrNotEnoughRecord
	}
	if err := c.mainStore.GetJobStore().DeleteJobs([]string{task.STATUS_FAILED}, uint64(time.Now().Unix())-uint64(opts.Expiration)); err != nil {
		log.Error("[ClearFailedJobs] error while deleting ExecCleanFailedJobs tasks", "err", err)
		return err
	}
	log.Info("Cleared failed jobs")

	return nil
}

func (c *Cleaner) Start() {
	v := reflect.ValueOf(c)
	count := v.Type().NumMethod()
	for i := 0; i < count; i++ {
		m := v.Method(i)
		name := v.Type().Method(i).Name
		if cronjob, ok := configs.AppConfig.CronJob[name]; ok {
			log.Info(fmt.Sprintf("Set cronjob for %v by using config: %+v", name, cronjob))

			c.scheduler.Cron(cronjob.Cron).Do(func() {
				opts := CleanOptions{
					Threshold:  cronjob.Threshold,
					Expiration: cronjob.Expiration,
				}

				m.Call([]reflect.Value{reflect.ValueOf(&opts)})
			})
		}
	}

	c.scheduler.StartAsync()
}

func (c *Cleaner) Stop() {
	c.scheduler.Stop()
}

func NewCleaner(store stores.BridgeStore, mainStore bridgeCoreStores.MainStore) *Cleaner {
	return &Cleaner{
		scheduler: gocron.NewScheduler(time.Local),
		store:     store,
		mainStore: mainStore,
	}
}
