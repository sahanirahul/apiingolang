package http

import (
	"apiingolang/activity/business/entities/dto"
	"apiingolang/activity/business/interfaces/irepo"
	"apiingolang/activity/business/utils/logging"
	"context"
	"net/http"
	"sync"
	"time"
)

type httprepo struct {
}

var once sync.Once
var repo *httprepo

func NewActivityHttpRepo() irepo.IHttpRepo {
	once.Do(func() {
		repo = &httprepo{}
	})
	return repo
}

func (cr *httprepo) GetActivityFromBoredApi(ctx context.Context) (*dto.Activity, error) {
	// call https://www.boredapi.com/api/activity here
	var response dto.Activity
	url := "https://www.boredapi.com/api/activity"
	httpreq := HttpRequest{URL: url, Body: nil, Timeout: 2 * time.Second, Method: http.MethodGet}
	status, err := httpreq.InitiateHttpCall(ctx, &response)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetching_activity_http_request", logging.ErrorLevel, logging.Fields{"error": err})
		return nil, err
	}
	if status != http.StatusOK {
		logging.Logger.WriteLogs(ctx, "fstatus_code_not_ok", logging.ErrorLevel, logging.Fields{"statusCode": status})
	}
	// panic("mock panic in boredapi")
	return &response, nil
}
