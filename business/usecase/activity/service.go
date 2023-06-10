package activity

import (
	"apiingolang/activity/business/entities/core"
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/entities/utility"
	"apiingolang/activity/business/entities/worker"
	"apiingolang/activity/business/interfaces/icore"
	"apiingolang/activity/business/interfaces/irepo"
	"apiingolang/activity/business/interfaces/iusecase"
	"apiingolang/activity/business/utils"
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
		fectchedActivitiesFromBoredApi = make(chan dto.Activity, 100)
		allFetchedActivities = utility.SyncList{}
		go c.storeFetchedActivities()
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
		fectchedActivitiesFromBoredApi <- ac
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

var allFetchedActivities utility.SyncList
var fectchedActivitiesFromBoredApi chan dto.Activity
var maxBatchSize = 50

// this will be called from cron worker
func (as *activityService) SaveFetchedActivitiesTillNow(ctx context.Context) error {
	activities := core.Activities{}
	actsTillNow := allFetchedActivities.GetAllEntryList()
	var errs []error
	err := utility.MapObjectToAnother(actsTillNow, &activities)
	if err != nil {
		fmt.Println("error while mapping activities : ", err.Error())
		return err
	}
	var wg sync.WaitGroup
	totalActivity := len(activities)
	batchNo := 0
	for i := 0; i < totalActivity; i = i + maxBatchSize {
		batchNo++
		end := utils.Min(i+maxBatchSize, totalActivity)
		batch := activities[i:end]
		fmt.Println("inserting batch number ", batchNo)
		wg.Add(1)
		// todo: should be done using worker pool to avoid goroutine outburst in case of huge number of activities
		go func(acts core.Activities) {
			defer wg.Done()
			err := as.activityrepo.BatchInsertActivities(ctx, batch)
			if err != nil {
				fmt.Println("error while inserting bacth : ", err.Error())
				errs = append(errs, err)
			}
		}(batch)

	}
	wg.Wait()
	if len(errs) > 0 {
		return utils.WrapErrors(errs)
	}
	return nil
}

func (as *activityService) storeFetchedActivities() {
	for {
		ac := <-fectchedActivitiesFromBoredApi
		allFetchedActivities.Append(ac)
	}
}
