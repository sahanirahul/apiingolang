package activity

import (
	"apiingolang/activity/business/entities/core"
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/entities/utility"
	"apiingolang/activity/business/entities/worker"
	"apiingolang/activity/business/utils"
	"apiingolang/activity/business/utils/logging"
	"context"
	"sync"
)

var allFetchedActivities utility.SyncList
var fectchedActivitiesFromBoredApi chan dto.Activity

func (as *activityService) storeFetchedActivities() {
	for {
		ac := <-fectchedActivitiesFromBoredApi
		allFetchedActivities.Append(ac)
	}
}

var maxBatchSize = 50

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
