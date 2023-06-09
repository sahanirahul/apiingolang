package irepo

import (
	"apiingolang/activity/business/entities/dto"
	"context"
)

type IActivityRepo interface {
	InsertActivities(ctx context.Context) error
}

type IHttpRepo interface {
	GetActivityFromBoredApi(ctx context.Context) (*dto.Activity, error)
}
