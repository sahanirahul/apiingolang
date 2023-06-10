package cron

import (
	"apiingolang/activity/business/interfaces/iusecase"
	"context"
	"fmt"
	"sync"
	"time"

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
		err := cro.activityService.SaveFetchedActivitiesTillNow(ctx)
		if err != nil {
			fmt.Println(fmt.Sprintf("%s error while executing cron:\\n%s", time.Now(), err.Error()))
		}
		fmt.Println(fmt.Sprintf("%s | cron execution success", time.Now()))
	})
	c.Start()
}
