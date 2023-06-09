package activity

import (
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/interfaces/irepo"
	"apiingolang/activity/business/interfaces/iusecase"
	"context"
	"fmt"
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

func (c *activityService) FetchActivities(ctx context.Context) (dto.Activities, error) {
	// fetch the 3 activities and other logic here
	ac, err := c.httprepo.GetActivityFromBoredApi(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return dto.Activities{dto.CreateActivity(ac)}, nil
}
