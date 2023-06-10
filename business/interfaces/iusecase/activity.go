package iusecase

import (
	"apiingolang/activity/business/entities/dto"
	"context"
)

type IActivityService interface {
	FetchActivities(ctx context.Context, cancel context.CancelFunc) (dto.Activities, error)
	SaveFetchedActivitiesTillNow(ctx context.Context) error
}
