package cron

import (
	"apiingolang/activity/business/interfaces/iusecase"
	"apiingolang/activity/business/utils/logging"
	"apiingolang/activity/middleware"
	"apiingolang/activity/middleware/corel"
	"sync"

	"github.com/robfig/cron"
)

type cronn struct {
	activityService iusecase.IActivityService
}

var once sync.Once
var cronObj *cronn

func StartNewCron(activityService iusecase.IActivityService) {
	once.Do(func() {
		cronObj := cronn{
			activityService: activityService,
		}
		cronObj.storeBoredApiActivities()
	})
}

func (cro *cronn) storeBoredApiActivities() {
	c := cron.New()
	c.AddFunc("@every 15s", func() {
		ctx := corel.CreateNewContext()
		// adding recovery for worker go routines
		defer func() {
			if err := recover(); err != nil {
				middleware.Recover(ctx, err)
			}
		}()
		logging.Logger.WriteLogs(ctx, "cron_started", logging.InfoLevel, logging.Fields{})
		err := cro.activityService.SaveFetchedActivitiesTillNow(ctx)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_while_executing_cron", logging.ErrorLevel, logging.Fields{"error": err})
		}
		logging.Logger.WriteLogs(ctx, "cron_finisthed", logging.InfoLevel, logging.Fields{})
	})
	c.Start()
}
