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
	"apiingolang/activity/business/utils/logging"
	"apiingolang/activity/middleware"
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
	poolgeneral  icore.IPool
}

var actServ *activityService
var conce sync.Once

func NewActivityService(httprepo irepo.IHttpRepo, activityrepo irepo.IActivityRepo, poolhttp icore.IPool, poolgeneral icore.IPool) iusecase.IActivityService {

	conce.Do(func() {
		actServ = &activityService{
			httprepo:     httprepo,
			activityrepo: activityrepo,
			poolhttp:     poolhttp,
			poolgeneral:  poolgeneral,
		}
		fectchedActivitiesFromBoredApi = make(chan dto.Activity, 100)
		allFetchedActivities = utility.SyncList{}
		go actServ.storeFetchedActivities()
	})
	return actServ
}

func (as *activityService) FetchActivities(ctx context.Context, cancel context.CancelFunc) (dto.Activities, error) {
	// fetch the 3 activities and other logic here
	activities := dto.Activities{}
	syncMap := utility.SyncMap{}
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		as.poolgeneral.AddJob(worker.NewJob(func() {
			as.fetchActivity(ctx, &syncMap, &wg, cancel)
		}))
	}
	// panic("this is a mock panic in api flow")
	wg.Wait()
	for _, val := range syncMap.GetAllEntry() {
		ac := dto.Activity{}
		err := utility.MapObjectToAnother(val, &ac)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_mapping_object_and_activity", logging.ErrorLevel, logging.Fields{"error": err, "value": val})
		}
		fectchedActivitiesFromBoredApi <- ac
		activities = append(activities, ac)
	}
	// panic("this is a mock panic in api flow")
	logging.Logger.WriteLogs(ctx, "activities_fetch_success", logging.InfoLevel, logging.Fields{"activities": activities})
	return activities, nil
}

func (as *activityService) fetchActivity(ctx context.Context, syncmap *utility.SyncMap, wg *sync.WaitGroup, cancel context.CancelFunc) (*dto.Activity, error) {
	// fetch a new activity here
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			middleware.Recover(ctx, err)
			cancel()
		}
	}()
	result := as.activityUtil(ctx, syncmap, cancel)
	select {
	case <-result:
	case <-time.After(2 * time.Second):
		logging.Logger.WriteLogs(ctx, "time_limit_exceeded_fetching_activities", logging.WarnLevel, logging.Fields{})
	}
	return nil, errors.New("time limit exceeded")
}

func (as *activityService) activityUtil(ctx context.Context, syncmap *utility.SyncMap, cancel context.CancelFunc) <-chan error {
	errCh := make(chan any, 1)
	result := make(chan error, 1)
	defer close(result)
	deadline := time.Now().Add(2 * time.Second)
	if dl, ok := ctx.Deadline(); ok {
		deadline = dl
	}
	for time.Since(deadline) < 0 {
		var ac *dto.Activity
		job := worker.NewJobWithDeadline(func() {
			defer func() {
				if err := recover(); err != nil {
					middleware.Recover(ctx, err)
					cancel()
					errCh <- err
				}
			}()
			a, err := as.httprepo.GetActivityFromBoredApi(ctx)
			if err != nil {
				logging.Logger.WriteLogs(ctx, "error_fetching_activity_from_bored_api", logging.ErrorLevel, logging.Fields{"error": err})
				return
			}
			ac = a
		}, deadline)
		as.poolhttp.AddJob(job)
		<-job.Done()
		if len(errCh) > 0 {
			e := <-errCh
			err := fmt.Errorf("panic while executing boredapi job. %v", e)
			// logging.Logger.WriteLogs(ctx, "Panic", logging.ErrorLevel, logging.Fields{"error": err})
			result <- err
			return result
		}
		if ac == nil || len(ac.Key) == 0 {
			continue
		}
		flag := syncmap.PutIfNotPresent(ac.Key, ac)
		if flag {
			result <- nil
			return result
		}
	}
	result <- nil
	return result
}

var allFetchedActivities utility.SyncList
var fectchedActivitiesFromBoredApi chan dto.Activity
var maxBatchSize = 5

// this will be called from cron worker
func (as *activityService) SaveFetchedActivitiesTillNow(ctx context.Context) error {
	activities := core.Activities{}
	actsTillNow := allFetchedActivities.GetAllEntryList()
	var errs []error
	err := utility.MapObjectToAnother(actsTillNow, &activities)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_json_unmarshal_activities", logging.ErrorLevel, logging.Fields{"error": err, "activities_till_now": actsTillNow})
		return err
	}
	var wg sync.WaitGroup
	totalActivity := len(activities)

	for i := 0; i < totalActivity; i = i + maxBatchSize {
		end := utils.Min(i+maxBatchSize, totalActivity)
		batch := activities[i:end]
		wg.Add(1)
		// todo: should be done using worker pool to avoid goroutine outburst in case of huge number of activities
		as.poolgeneral.AddJob(worker.NewJob(func() {
			defer wg.Done()
			err := as.activityrepo.BatchInsertActivities(ctx, batch)
			if err != nil {
				logging.Logger.WriteLogs(ctx, "error_batch_insert_into_db", logging.ErrorLevel, logging.Fields{"error": err, "batch": batch})
				errs = append(errs, err)
			}
		}))

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
