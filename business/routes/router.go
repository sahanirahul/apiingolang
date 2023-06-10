package routes

import (
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/interfaces/iusecase"
	"apiingolang/activity/business/repository/db"
	"apiingolang/activity/business/repository/http"
	"apiingolang/activity/business/usecase/activity"
	"apiingolang/activity/business/worker"
	"context"
	"database/sql"
	"fmt"
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
	//todo: CleanArch complete
	//call activity service
	timelimitexceeded := false
	activitiesList := dto.Activities{}
	wg := make(chan struct{}, 1)

	go maxtimewait(wg, &timelimitexceeded)
	go ar.getActivities(c, &activitiesList, wg)
	<-wg

	if timelimitexceeded {
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
}

func maxtimewait(wg chan struct{}, timelimitexceeded *bool) {
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	*timelimitexceeded = true
	wg <- struct{}{}
}

func (ar *activityRouter) getActivities(ctx context.Context, activitiesList *dto.Activities, wg chan struct{}) {
	activities, err := ar.activityManager.FetchActivities(ctx)
	// time.Sleep(3 * time.Second)
	if err != nil {
		fmt.Println(err)
	}
	*activitiesList = activities
	wg <- struct{}{}
}
