package iusecase

import (
	"apiingolang/activity/business/entities/dto"
	"context"
)

type IActivityService interface {
	FetchActivities(ctx context.Context) (dto.Activities, error)
}
