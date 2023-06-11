package worker

import "time"

type Job struct {
	job      func()
	status   chan struct{}
	deadline time.Time
}

func NewJobWithDeadline(job func(), dl time.Time) *Job {
	return &Job{
		job:      job,
		deadline: dl,
		status:   make(chan struct{}, 1),
	}
}

func NewJob(job func()) *Job {
	return &Job{
		job:    job,
		status: make(chan struct{}, 1),
	}
}

func (j *Job) Execute() {
	if j.deadline.IsZero() || time.Since(j.deadline) < 0 {
		// execute the job only if deadline is not reached or deadline is not set
		j.job()
	}
	j.status <- struct{}{}
}

func (j *Job) Done() <-chan struct{} {
	return j.status
}
