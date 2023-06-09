package routes

import (
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/interfaces/iusecase"
	"apiingolang/activity/business/repository/db"
	"apiingolang/activity/business/repository/http"
	"apiingolang/activity/business/usecase/activity"
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
	activityService := activity.NewActivityService(httprepo, dbrepo)
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
	activitiesList := &dto.Activities{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go maxtimewait(&wg, &timelimitexceeded)
	go ar.getActivities(c, activitiesList, &wg)
	wg.Wait()

	if timelimitexceeded {
		c.JSON(corehttp.StatusUnprocessableEntity, gin.H{
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

func maxtimewait(wg *sync.WaitGroup, timelimitexceeded *bool) {
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	*timelimitexceeded = true
	wg.Done()
	wg.Done()
}

func (ar *activityRouter) getActivities(ctx context.Context, activitiesList *dto.Activities, wg *sync.WaitGroup) {
	activities, err := ar.activityManager.FetchActivities(ctx)
	// time.Sleep(3 * time.Second)
	if err != nil {
		fmt.Println(err)
	}
	activitiesList = &activities
	wg.Done()
	wg.Done()
}
