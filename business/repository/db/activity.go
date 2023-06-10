package db

import (
	"apiingolang/activity/business/entities/core"
	"apiingolang/activity/business/interfaces/icore"
	"apiingolang/activity/business/interfaces/irepo"
	"apiingolang/activity/business/utils/logging"
	"context"
	"fmt"
	"strings"
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

func (ar *activityrepo) InsertActivity(ctx context.Context, act core.Activity) error {
	query := "INSERT INTO activity (activity_key,activity_content) VALUES ($1,$2)"
	conn, err := ar.db.Conn(ctx)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetch_db_conn_insert", logging.ErrorLevel, logging.Fields{"error": err})
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, query, act.Key, act.Activity)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_batch_insert_execute", logging.ErrorLevel, logging.Fields{"error": err, "query": query, "activity": act})
		return err
	}
	return nil
}

func (ar *activityrepo) BatchInsertActivities(ctx context.Context, activities core.Activities) error {
	//insert activities in db
	if len(activities) == 0 {
		return nil
	}
	if len(activities) > 100 {
		return fmt.Errorf("batch size=%d is greater than 100", len(activities))
	}
	values := []interface{}{}
	for _, ac := range activities {
		values = append(values, ac.Key, ac.Activity)
	}
	query := "INSERT INTO activity (activity_key,activity_content) VALUES %s"
	valueStr := ""
	for i := 1; i <= len(values); i = i + 2 {
		valueStr = valueStr + fmt.Sprintf("($%d,$%d),", i, i+1)
	}
	query = fmt.Sprintf(query, valueStr)
	query = strings.TrimSuffix(query, ",")
	logging.Logger.WriteLogs(ctx, "batch_insert_activities_query", logging.InfoLevel, logging.Fields{"query": query, "qargs": values})

	conn, err := ar.db.Conn(ctx)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetch_db_conn_insert_batch", logging.ErrorLevel, logging.Fields{"error": err})
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, query, values...)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_batch_insert_execute", logging.ErrorLevel, logging.Fields{"error": err, "query": query, "values": values})
		return err
	}
	return nil
}
