package activity

import (
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/entities/utility"
	"apiingolang/activity/business/entities/worker"
	"apiingolang/activity/business/interfaces/icore"
	"apiingolang/activity/business/interfaces/irepo"
	"apiingolang/activity/business/interfaces/iusecase"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type activityService struct {
	httprepo     irepo.IHttpRepo
	activityrepo irepo.IActivityRepo
	poolhttp     icore.IPool
}

var c *activityService
var conce sync.Once

func NewActivityService(httprepo irepo.IHttpRepo, activityrepo irepo.IActivityRepo, poolhttp icore.IPool) iusecase.IActivityService {

	conce.Do(func() {
		c = &activityService{
			httprepo:     httprepo,
			activityrepo: activityrepo,
			poolhttp:     poolhttp,
		}
	})
	return c
}

func (as *activityService) FetchActivities(ctx context.Context) (dto.Activities, error) {
	// fetch the 3 activities and other logic here
	activities := dto.Activities{}
	syncMap := utility.SyncMap{}
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go c.fetchActivity(ctx, &syncMap, &wg)
	}
	wg.Wait()
	for _, val := range syncMap.GetAllEntry() {
		ac := dto.Activity{}
		err := utility.MapObjectToAnother(val, &ac)
		if err != nil {
			fmt.Println(err)
			continue
		}
		activities = append(activities, ac)
	}
	return activities, nil
}

func (as *activityService) fetchActivity(ctx context.Context, syncmap *utility.SyncMap, wg *sync.WaitGroup) (*dto.Activity, error) {
	// fetch a new activity here
	start := time.Now()
	defer wg.Done()
	for time.Since(start) <= time.Second*2 {
		var ac *dto.Activity
		job := worker.NewJob(func() {
			a, err := c.httprepo.GetActivityFromBoredApi(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			ac = a
		})
		as.poolhttp.AddJob(job)
		<-job.Done()
		if ac == nil || len(ac.Key) == 0 {
			continue
		}
		flag := syncmap.PutIfNotPresent(ac.Key, ac)
		if flag {
			return ac, nil
		}
	}
	return nil, errors.New("time limit exceeded")
}
