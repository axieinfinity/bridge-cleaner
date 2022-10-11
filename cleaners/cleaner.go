package cleaners

import (
	"time"

	"github.com/axieinfinity/bridge-cleaner/configs"
	bridgeCoreStores "github.com/axieinfinity/bridge-core/stores"
	"github.com/axieinfinity/bridge-v2/stores"
	"github.com/axieinfinity/bridge-v2/task"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-co-op/gocron"
)

type Cleaner struct {
	scheduler *gocron.Scheduler
	store     stores.BridgeStore
	mainStore bridgeCoreStores.MainStore
}

func (c *Cleaner) clearSuccessTasks() error {
	_, err := c.scheduler.Cron(configs.AppConfig.Cleaner.ClearSuccessTaskScheduler).Do(func() {
		log.Info("Doing clear success tasks")
		if c.store.GetTaskStore().Count() <= configs.AppConfig.Cleaner.ClearSuccessTaskThreshold {
			log.Info("Skipping clear success tasks")
			return
		}
		if err := c.store.GetTaskStore().DeleteTasks([]string{task.STATUS_DONE}, uint64(time.Now().Unix())-uint64(configs.AppConfig.Cleaner.ClearSuccessTaskExpiration)); err != nil {
			log.Error("[ExecClearSuccessTask] error while deleting done tasks", "err", err)
		}
	})
	return err
}

func (c *Cleaner) clearFailedTasks() error {
	_, err := c.scheduler.Cron(configs.AppConfig.Cleaner.ClearFailedTaskScheduler).Do(func() {
		log.Info("Doing clear failed tasks")
		if c.store.GetTaskStore().Count() <= configs.AppConfig.Cleaner.ClearSuccessTaskThreshold {
			log.Info("Skipping clear failed tasks")
			return
		}
		if err := c.store.GetTaskStore().DeleteTasks([]string{task.STATUS_FAILED}, uint64(time.Now().Unix())-uint64(configs.AppConfig.Cleaner.ClearFailedTaskExpiration)); err != nil {
			log.Error("[ExecClearFailedTask] error while deleting failed tasks", "err", err)
		}
	})
	return err
}

func (c *Cleaner) clearEvents() error {
	_, err := c.scheduler.Cron(configs.AppConfig.Cleaner.ClearEventScheduler).Do(func() {
		log.Info("Doing clear events")
		if c.mainStore.GetEventStore().Count() <= configs.AppConfig.Cleaner.ClearEventThreshold {
			log.Info("Skipping clear events")
			return
		}
		if err := c.mainStore.GetEventStore().DeleteEvents(uint64(time.Now().Unix()) - uint64(configs.AppConfig.Cleaner.ClearEventExpiration)); err != nil {
			log.Error("[ExecClearFailedTask] error while deleting failed tasks", "err", err)
		}
	})
	return err
}

func (c *Cleaner) clearSuccessJobs() error {
	_, err := c.scheduler.Cron(configs.AppConfig.Cleaner.CLearSuccessJobScheduler).Do(func() {
		log.Info("Doing clear success jobs")
		if c.mainStore.GetJobStore().Count() <= configs.AppConfig.Cleaner.CLearSuccessJobThreshold {
			log.Info("Skipping clear success jobs")
			return
		}
		if err := c.mainStore.GetJobStore().DeleteJobs([]string{task.STATUS_DONE}, uint64(time.Now().Unix())-uint64(configs.AppConfig.Cleaner.ClearSuccessJobExpiration)); err != nil {
			log.Error("[ExecClearSuccessTask] error while deleting done tasks", "err", err)
		}
	})
	return err
}

func (c *Cleaner) clearFailedJobs() error {
	_, err := c.scheduler.Cron(configs.AppConfig.Cleaner.CLearFailedJobScheduler).Do(func() {
		log.Info("Doing clear failed jobs")
		if c.mainStore.GetJobStore().Count() <= configs.AppConfig.Cleaner.CLearFailedJobThreshold {
			log.Info("Skipping clear failed jobs")
			return
		}
		if err := c.mainStore.GetJobStore().DeleteJobs([]string{task.STATUS_FAILED}, uint64(time.Now().Unix())-uint64(configs.AppConfig.Cleaner.ClearFailedJobExpiration)); err != nil {
			log.Error("[ExecClearFailedJobs] error while deleting ExecCleanFailedJobs tasks", "err", err)
		}
	})
	return err
}

func (c *Cleaner) Start() {
	c.clearEvents()
	c.clearFailedJobs()
	c.clearFailedTasks()
	c.clearSuccessJobs()
	c.clearSuccessTasks()

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
