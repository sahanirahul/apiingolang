package activity

import (
	"apiingolang/activity/business/interfaces/irepo"
	"apiingolang/activity/business/interfaces/iusecase"
	"sync"
)

type activityService struct {
	httprepo     irepo.IHttpRepo
	activityrepo irepo.IActivityRepo
}

var c *activityService
var conce sync.Once

func NewActivityService(httprepo irepo.IHttpRepo, activityrepo irepo.IActivityRepo) iusecase.IActivityService {

	conce.Do(func() {
		c = &activityService{
			httprepo:     httprepo,
			activityrepo: activityrepo,
		}
	})
	return c
}

func (c *activityService) FetchActivities() error {
	// fetch the 3 activities and other logic here
	return nil
}
