package worker

type Job struct {
	job    func()
	status chan struct{}
}

func NewJob(job func()) *Job {
	return &Job{
		job:    job,
		status: make(chan struct{}, 1),
	}
}

func (j *Job) Execute() {
	j.job()
	j.status <- struct{}{}
}

func (j *Job) Done() <-chan struct{} {
	return j.status
}
