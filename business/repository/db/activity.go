package db

import (
	"apiingolang/activity/business/interfaces/icore"
	"apiingolang/activity/business/interfaces/irepo"
	"context"
	"sync"
)

type activityrepo struct {
	db icore.IDB
}

var once sync.Once
var repo *activityrepo

func NewActivityRepo(db icore.IDB) irepo.IActivityRepo {
	once.Do(func() {
		repo = &activityrepo{
			db: db,
		}
	})
	return repo
}

func (cr *activityrepo) InsertActivities(ctx context.Context) error {
	//insert activities in db
	return nil
}
