package irepo

import (
	"apiingolang/activity/business/entities/core"
	"apiingolang/activity/business/entities/dto"
	"context"
)

type IActivityRepo interface {
	InsertActivity(ctx context.Context, act core.Activity) error
	BatchInsertActivities(ctx context.Context, activities core.Activities) error
}

type IHttpRepo interface {
	GetActivityFromBoredApi(ctx context.Context) (*dto.Activity, error)
}
