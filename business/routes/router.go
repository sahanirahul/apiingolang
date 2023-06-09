package routes

import (
	"apiingolang/activity/business/interfaces/iusecase"
	"apiingolang/activity/business/repository/db"
	"apiingolang/activity/business/repository/http"
	"apiingolang/activity/business/usecase/activity"
	"database/sql"
	corehttp "net/http"
	"sync"

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
	apigroup.GET("/v1/activities", r.getActivities)
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

func (ar *activityRouter) getActivities(c *gin.Context) {
	//todo: CleanArch complete
	//call activity service
	activities, err := ar.activityManager.FetchActivities(c)
	if err != nil {
		c.JSON(corehttp.StatusUnprocessableEntity, gin.H{
			"status":  true,
			"message": "failure",
			"error":   err.Error(),
		})
	}
	c.JSON(corehttp.StatusOK, gin.H{
		"status":     true,
		"message":    "success",
		"activities": activities,
	})
}
