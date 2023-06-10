package worker

import (
	entitiesworker "apiingolang/activity/business/entities/worker"
	"apiingolang/activity/business/interfaces/icore"
	"fmt"
	"sync/atomic"
)

type worker struct {
	numberOfWorker int
	jobs           chan *entitiesworker.Job
}

func (w *worker) AddJob(job *entitiesworker.Job) {
	w.jobs <- job
}

func (w *worker) start() {
	fmt.Println("starting workers")
	for i := 1; i <= w.numberOfWorker; i++ {
		go w.run(i)
	}
}

var numberOfJobs atomic.Int64

func (w *worker) run(workerId int) {
	fmt.Println("starting worker ", workerId)
	for {
		job := <-w.jobs
		val := numberOfJobs.Add(1)
		fmt.Println(fmt.Sprintf("worker %d executing job %d", workerId, val))
		job.Execute()
	}
}

func NewWorkerPool(maxworkers int, jobQSize int) icore.IPool {
	w := new(worker)
	w.jobs = make(chan *entitiesworker.Job, jobQSize)
	w.numberOfWorker = maxworkers
	w.start()
	return w
}
