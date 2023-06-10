package routes

import (
	"apiingolang/activity/business/cron"
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/interfaces/iusecase"
	"apiingolang/activity/business/repository/db"
	"apiingolang/activity/business/repository/http"
	"apiingolang/activity/business/usecase/activity"
	"apiingolang/activity/business/utils/logging"
	"apiingolang/activity/business/worker"
	"apiingolang/activity/middleware"
	"context"
	"database/sql"
	corehttp "net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func provideActivityRouter(dbconn *sql.DB) *activityRouter {
	httprepo := http.NewActivityHttpRepo()
	dbrepo := db.NewActivityRepo(dbconn)
	httpWorkerPool := worker.NewWorkerPool(3, 3) // this pool wil make sure that the boredapi is called on 3 at a time
	activityService := activity.NewActivityService(httprepo, dbrepo, httpWorkerPool)
	cron.StartNewCron(activityService)
	return newActivityRouter(activityService)

}

func ActivityRoutes(apigroup *gin.RouterGroup, db *sql.DB) {
	r := provideActivityRouter(db)
	apigroup.GET("/v1/activities", r.processActivitiesRequest)
}

type activityRouter struct {
	activityManager iusecase.IActivityService
}

var routeOnce sync.Once
var arouter *activityRouter

func newActivityRouter(as iusecase.IActivityService) *activityRouter {
	routeOnce.Do(func() {
		arouter = &activityRouter{
			activityManager: as,
		}
	})
	return arouter
}

func (ar *activityRouter) processActivitiesRequest(c *gin.Context) {
	ctx, cancel := context.WithCancel(c)
	timelimitexceeded := false
	activitiesList := dto.Activities{}
	wg := make(chan struct{}, 1)

	go maxtimewait(wg, &timelimitexceeded)
	go ar.getActivities(ctx, &activitiesList, wg, cancel)
	select {
	case <-wg:
		if timelimitexceeded {
			logging.Logger.WriteLogs(ctx, "time_limit_exceeded", logging.WarnLevel, logging.Fields{})
			c.JSON(corehttp.StatusRequestTimeout, gin.H{
				"status":  false,
				"message": "failure",
				"error":   "(Activity-API not available)",
			})
		} else {
			c.JSON(corehttp.StatusOK, gin.H{
				"status":     true,
				"message":    "success",
				"activities": activitiesList,
			})
		}
	case <-ctx.Done():
		logging.Logger.WriteLogs(ctx, "cancel button hit", logging.InfoLevel, logging.Fields{})
		c.JSON(corehttp.StatusInternalServerError, gin.H{
			"message": "something went wrong",
			"status":  false,
		})
		return
	}

}

func maxtimewait(wg chan struct{}, timelimitexceeded *bool) {
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	*timelimitexceeded = true
	wg <- struct{}{}
}

func (ar *activityRouter) getActivities(ctx context.Context, activitiesList *dto.Activities, wg chan struct{}, cancel context.CancelFunc) {
	defer func() {
		if err := recover(); err != nil {
			middleware.Recover(ctx, err)
			cancel()
		}
	}()
	activities, err := ar.activityManager.FetchActivities(ctx, cancel)
	// time.Sleep(3 * time.Second)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_processing_activity_request", logging.ErrorLevel, logging.Fields{"error": err})
	}
	*activitiesList = activities
	wg <- struct{}{}
}
