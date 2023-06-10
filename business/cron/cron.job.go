package cron

import (
	"apiingolang/activity/business/interfaces/iusecase"
	"apiingolang/activity/business/utils/logging"
	"context"
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
	ctx := context.Background()
	c := cron.New()

	c.AddFunc("@every 15s", func() {
		logging.Logger.WriteLogs(ctx, "cron_started", logging.InfoLevel, logging.Fields{})
		err := cro.activityService.SaveFetchedActivitiesTillNow(ctx)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_while_executing_cron", logging.ErrorLevel, logging.Fields{"error": err})
		}
		logging.Logger.WriteLogs(ctx, "cron_finisthed", logging.InfoLevel, logging.Fields{})
	})
	c.Start()
}
