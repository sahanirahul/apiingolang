package http

import (
	"apiingolang/activity/business/interfaces/irepo"
	"sync"
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

func (cr *httprepo) CallExternal() error {
	// call https://www.boredapi.com/api/activity here
	return nil
}
