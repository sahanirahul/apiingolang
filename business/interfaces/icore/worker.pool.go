package icore

import (
	"apiingolang/activity/business/entities/worker"
)

type IPool interface {
	AddJob(job *worker.Job)
}
